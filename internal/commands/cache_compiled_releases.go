package commands

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/pivotal-cf/jhanda"
	"github.com/pivotal-cf/om/api"
	"gopkg.in/yaml.v2"

	"github.com/pivotal-cf/kiln/internal/commands/flags"
	"github.com/pivotal-cf/kiln/internal/fetcher"
	"github.com/pivotal-cf/kiln/internal/om"
	"github.com/pivotal-cf/kiln/pkg/cargo"
	"github.com/pivotal-cf/kiln/pkg/release"
)

//go:generate counterfeiter -o ./fakes/ops_manager_release_cache_source.go --fake-name OpsManagerReleaseCacheSource . OpsManagerReleaseCacheSource
//go:generate counterfeiter -o ./fakes/release_cache_bucket.go --fake-name ReleaseCacheBucket . ReleaseCacheBucket

type (
	OpsManagerReleaseCacheSource interface {
		om.GetBoshEnvironmentAndSecurityRootCACertificateProvider
		GetStagedProductManifest(guid string) (string, error)
		GetStagedProductByName(productName string) (api.StagedProductsFindOutput, error)
	}

	ReleaseCacheBucket interface {
		UploadRelease(spec release.Requirement, file io.Reader) (release.Remote, error)
	}
)

type CacheCompiledReleases struct {
	Options struct {
		flags.Standard
		om.ClientConfiguration

		UploadTargetID string `           long:"upload-target-id"   required:"true"    description:"the ID of the release source where the built release will be uploaded"`
		ReleasesDir    string `short:"rd" long:"releases-directory" default:"releases" description:"path to a directory to download releases into"`
		Name           string `short:"n"  long:"name"               default:"cf"       description:"name of the tile"` // TODO: parse from base.yml
	}

	Logger *log.Logger
	FS     billy.Filesystem

	ReleaseCache   func(kilnfile cargo.Kilnfile) fetcher.MultiReleaseSource
	Bucket         func(kilnfile cargo.Kilnfile) (ReleaseCacheBucket, error)
	OpsManager     func(om.ClientConfiguration) (OpsManagerReleaseCacheSource, error)
	Director       func(om.ClientConfiguration, om.GetBoshEnvironmentAndSecurityRootCACertificateProvider) (boshdir.Director, error)
}

func NewCacheCompiledReleases() *CacheCompiledReleases {
	cmd := &CacheCompiledReleases{
		FS:     osfs.New(""),
		Logger: log.Default(),
	}
	cmd.ReleaseCache = func(kilnfile cargo.Kilnfile) fetcher.MultiReleaseSource {
		releaseSourceProvider := fetcher.NewReleaseSourceRepo(kilnfile, cmd.Logger)
		return releaseSourceProvider.MultiReleaseSource(false)
	}
	cmd.Bucket = func(kilnfile cargo.Kilnfile) (ReleaseCacheBucket, error) {
		return cmd.s3Bucket(kilnfile)
	}
	cmd.OpsManager = func(conf om.ClientConfiguration) (OpsManagerReleaseCacheSource, error) {
		return conf.API()
	}
	cmd.Director = om.BoshDirector
	return cmd
}

func (cmd *CacheCompiledReleases) WithLogger(logger *log.Logger) *CacheCompiledReleases {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	cmd.Logger = logger
	return cmd
}

func (cmd CacheCompiledReleases) Execute(args []string) error {
	err := flags.LoadFlagsWithDefaults(&cmd.Options, args, cmd.FS.Stat)
	if err != nil {
		return err
	}

	kilnfile, lock, err := cmd.Options.LoadKilnfiles(cmd.FS, nil)
	if err != nil {
		return fmt.Errorf("failed to load kilnfiles: %w", err)
	}

	omAPI, deploymentName, stagedStemcellOS, stagedStemcellVersion, err := cmd.fetchProductDeploymentData()
	if err != nil {
		return err
	}

	if stagedStemcellOS != lock.Stemcell.OS || stagedStemcellVersion != lock.Stemcell.Version {
		return fmt.Errorf(
			"staged stemcell (%s %s) and lock stemcell (%s %s) do not match",
			stagedStemcellOS, stagedStemcellVersion,
			lock.Stemcell.OS, lock.Stemcell.Version,
		)
	}

	var nonCompiledReleases []cargo.ReleaseLock

	cache := cmd.ReleaseCache(kilnfile)
	for _, rel := range lock.Releases {
		remote, found, err := cache.GetMatchedRelease(release.Requirement{
			Name:            rel.Name,
			Version:         rel.Version,
			StemcellOS:      lock.Stemcell.OS,
			StemcellVersion: lock.Stemcell.Version,
		})
		if err != nil {
			return fmt.Errorf("failed check for matched release: %w", err)
		}
		if !found {
			nonCompiledReleases = append(nonCompiledReleases, rel)
			continue
		}
		err = updateLock(lock, remote)
		if err != nil {
			return fmt.Errorf("failed to update lock file: %w", err)
		}
	}

	switch len(nonCompiledReleases) {
	case 0:
		cmd.Logger.Print("cache already contains releases matching constraint\n")
		return nil
	case 1:
		cmd.Logger.Printf("1 release is not publishable\n")
	default:
		cmd.Logger.Printf("%d releases are not publishable\n", len(nonCompiledReleases))
	}

	for _, rel := range nonCompiledReleases {
		cmd.Logger.Printf("\t%s %s compiled with %s %s not found in cache\n", rel.Name, rel.Version, lock.Stemcell.OS, lock.Stemcell.Version)
	}

	bucket, err := cmd.Bucket(kilnfile)
	if err != nil {
		return fmt.Errorf("failed to configure release cache: %w", err)
	}

	bosh, err := cmd.Director(cmd.Options.ClientConfiguration, omAPI)
	if err != nil {
		return err
	}

	osVersionSlug := boshdir.NewOSVersionSlug(stagedStemcellOS, stagedStemcellVersion)

	deployment, err := bosh.FindDeployment(deploymentName)
	if err != nil {
		return err
	}

	cmd.Logger.Printf("exporting from bosh deployment %s\n", deploymentName)

	err = cmd.FS.MkdirAll(cmd.Options.ReleasesDir, 0777)
	if err != nil {
		return fmt.Errorf("failed to create release directory: %w", err)
	}

	for _, rel := range nonCompiledReleases {
		requirement := release.Requirement{
			Name:            rel.Name,
			Version:         rel.Version,
			StemcellOS:      stagedStemcellOS,
			StemcellVersion: stagedStemcellVersion,
		}

		newRemote, err := cmd.cacheRelease(bosh, bucket, deployment, requirement, osVersionSlug)
		if err != nil {
			return fmt.Errorf("failed to cache release %s: %w", rel.Name, err)
		}

		err = updateLock(lock, newRemote)
		if err != nil {
			return fmt.Errorf("failed to lock release %s: %w", rel.Name, err)
		}
	}

	err = cmd.Options.Standard.SaveKilnfileLock(cmd.FS, lock)
	if err != nil {
		return err
	}

	cmd.Logger.Printf("DON'T FORGET TO MAKE A COMMIT AND PR\n")

	return nil
}

func (cmd CacheCompiledReleases) fetchProductDeploymentData() (_ OpsManagerReleaseCacheSource, deploymentName, stemcellOS, stemcellVersion string, _ error) {
	omAPI, err := cmd.OpsManager(cmd.Options.ClientConfiguration)
	if err != nil {
		return nil, "", "", "", err
	}

	stagedProduct, err := omAPI.GetStagedProductByName(cmd.Options.Name)
	if err != nil {
		return nil, "", "", "", err
	}

	stagedManifest, err := omAPI.GetStagedProductManifest(stagedProduct.Product.GUID)
	if err != nil {
		return nil, "", "", "", err
	}

	var manifest struct {
		Name      string `yaml:"name"`
		Stemcells []struct {
			OS      string `yaml:"os"`
			Version string `yaml:"version"`
		} `yaml:"stemcells"`
	}

	if err := yaml.Unmarshal([]byte(stagedManifest), &manifest); err != nil {
		return nil, "", "", "", err
	}

	if len(manifest.Stemcells) == 0 {
		return nil, "", "", "", errors.New("manifest stemcell not set")
	}
	stagedStemcell := manifest.Stemcells[0]

	return omAPI, manifest.Name, stagedStemcell.OS, stagedStemcell.Version, nil
}

func (cmd CacheCompiledReleases) cacheRelease(bosh boshdir.Director, bucket ReleaseCacheBucket, deployment boshdir.Deployment, req release.Requirement, osVersionSlug boshdir.OSVersionSlug) (release.Remote, error) {
	cmd.Logger.Printf("\texporting %s %s\n", req.Name, req.Version)
	result, err := deployment.ExportRelease(boshdir.NewReleaseSlug(req.Name, req.Version), osVersionSlug, nil)
	if err != nil {
		return release.Remote{}, err
	}

	cmd.Logger.Printf("\tdownloading %s %s\n", req.Name, req.Version)
	releaseFilePath, err := cmd.saveReleaseLocally(bosh, cmd.Options.ReleasesDir, req, result)
	if err != nil {
		return release.Remote{}, err
	}

	cmd.Logger.Printf("\tuploading %s %s\n", req.Name, req.Version)
	remoteRelease, err := cmd.uploadLocalRelease(req, releaseFilePath, bucket)
	if err != nil {
		return release.Remote{}, err
	}

	return remoteRelease, nil
}

func (cmd *CacheCompiledReleases) s3Bucket(kilnfile cargo.Kilnfile) (fetcher.S3ReleaseSource, error) {
	for _, source := range kilnfile.ReleaseSources {
		if source.ID != cmd.Options.UploadTargetID {
			continue
		}
		return fetcher.S3ReleaseSourceFromConfig(source, cmd.Logger), nil
	}
	return fetcher.S3ReleaseSource{}, errors.New("release source not found")
}

func updateLock(lock cargo.KilnfileLock, release release.Remote) error {
	for index, releaseLock := range lock.Releases {
		if release.Name != releaseLock.Name {
			continue
		}
		lock.Releases[index] = cargo.ReleaseLock{
			Name:         release.Name,
			Version:      release.Version,
			RemoteSource: release.SourceID,
			RemotePath:   release.RemotePath,
			SHA1:         release.SHA,
		}
		return nil
	}
	return fmt.Errorf("existing release not found in Kilnfile.lock")
}

func (cmd *CacheCompiledReleases) uploadLocalRelease(spec release.Requirement, fp string, uploader ReleaseCacheBucket) (release.Remote, error) {
	f, err := cmd.FS.Open(fp)
	if err != nil {
		return release.Remote{}, err
	}
	defer func() {
		_ = f.Close()
	}()
	return uploader.UploadRelease(spec, f)
}

func (cmd *CacheCompiledReleases) saveReleaseLocally(director boshdir.Director, relDir string, req release.Requirement, res boshdir.ExportReleaseResult) (string, error) {
	fileName := fmt.Sprintf("%s-%s-%s-%s.tgz", req.Name, req.Version, req.StemcellOS, req.StemcellVersion)
	filePath := filepath.Join(relDir, fileName)

	f, err := cmd.FS.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	checkSum := sha256.New()

	w := io.MultiWriter(f, checkSum)

	err = director.DownloadResourceUnchecked(res.BlobstoreID, w)
	if err != nil {
		_ = os.Remove(filePath)
		return "", err
	}

	if sum := fmt.Sprintf("sha256:%x", checkSum.Sum(nil)); sum != res.SHA1 {
		return "", fmt.Errorf("checksums do not match got %q but expected %q", sum, res.SHA1)
	}

	return filePath, nil
}

func (cmd CacheCompiledReleases) Usage() jhanda.Usage {
	return jhanda.Usage{
		Description:      "Downloads compiled bosh releases from an Tanzu Ops Manager bosh director and then uploads them to a bucket",
		ShortDescription: "Cache compiled releases",
		Flags:            cmd.Options,
	}
}