package component

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
)

type BOSHIOReleaseSource struct {
	Identifier  string `yaml:"id,omitempty"`
	Publishable bool   `yaml:"publishable,omitempty"`

	CustomURI string `yaml:"customURI,omitempty"`
}

func (src *BOSHIOReleaseSource) ConfigurationErrors() []error {
	return nil
}

func (src *BOSHIOReleaseSource) ID() string {
	if src.Identifier != "" {
		return src.Identifier
	}
	return ReleaseSourceTypeBOSHIO
}
func (src *BOSHIOReleaseSource) IsPublishable() bool { return src.Publishable }
func (src *BOSHIOReleaseSource) Type() string        { return ReleaseSourceTypeBOSHIO }

func (src *BOSHIOReleaseSource) GetMatchedRelease(ctx context.Context, _ *log.Logger, requirement Spec) (Lock, error) {
	requirement = requirement.UnsetStemcell()

	for _, repo := range organizations {
		for _, suf := range suffixes {
			fullName := repo + "/" + requirement.Name + suf
			exists, err := src.releaseExistOnBoshIO(ctx, fullName, requirement.Version)
			if err != nil {
				return Lock{}, err
			}

			if exists {
				builtRelease := src.createReleaseRemote(requirement, fullName)
				return builtRelease, nil
			}
		}
	}
	return Lock{}, ErrNotFound
}

func (src *BOSHIOReleaseSource) FindReleaseVersion(ctx context.Context, _ *log.Logger, spec Spec) (Lock, error) {
	spec = spec.UnsetStemcell()

	constraint, err := spec.VersionConstraints()
	if err != nil {
		return Lock{}, err
	}

	var validReleases []releaseResponse

	for _, repo := range organizations {
		for _, suf := range suffixes {
			fullName := repo + "/" + spec.Name + suf
			releaseResponses, err := src.getReleases(ctx, fullName)
			if err != nil {
				return Lock{}, err
			}

			for _, release := range releaseResponses {
				version, _ := semver.NewVersion(release.Version)
				if constraint.Check(version) {
					validReleases = append(validReleases, release)
				}
			}
			if len(validReleases) == 0 {
				continue
			}
			spec.Version = validReleases[0].Version
			lock := src.createReleaseRemote(spec, fullName)
			lock.SHA1 = validReleases[0].SHA
			return lock, nil
		}
	}
	return Lock{}, ErrNotFound
}

func (src *BOSHIOReleaseSource) DownloadRelease(ctx context.Context, logger *log.Logger, releaseDir string, remoteRelease Lock) (Local, error) {
	logger.Printf(logLineDownload, remoteRelease.Name, ReleaseSourceTypeBOSHIO, src.ID())

	downloadURL := remoteRelease.RemotePath

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return Local{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Local{}, err
	}

	filePath := filepath.Join(releaseDir, fmt.Sprintf("%s-%s.tgz", remoteRelease.Name, remoteRelease.Version))

	out, err := os.Create(filePath)
	if err != nil {
		return Local{}, err
	}
	defer closeAndIgnoreError(out)

	_, err = io.Copy(out, resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return Local{}, err
	}

	_, err = out.Seek(0, 0)
	if err != nil {
		return Local{}, fmt.Errorf("error reseting file cursor: %w", err) // untested
	}

	hash := sha1.New()
	_, err = io.Copy(hash, out)
	if err != nil {
		return Local{}, fmt.Errorf("error hashing file contents: %w", err) // untested
	}

	remoteRelease.SHA1 = hex.EncodeToString(hash.Sum(nil))

	return Local{Lock: remoteRelease, LocalPath: filePath}, nil
}

func (src *BOSHIOReleaseSource) createReleaseRemote(spec Spec, fullName string) Lock {
	downloadURL := fmt.Sprintf("%s/d/github.com/%s?v=%s", src.serverURI(), fullName, spec.Version)
	releaseRemote := spec.Lock()
	releaseRemote.RemoteSource = src.ID()
	releaseRemote.RemotePath = downloadURL
	return releaseRemote
}

func (src *BOSHIOReleaseSource) getReleases(ctx context.Context, name string) ([]releaseResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/releases/github.com/%s", src.serverURI(), name), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bosh.io API is down with error: %w", err)
	}
	if resp.StatusCode >= 500 {
		return nil, (*ResponseStatusCodeError)(resp)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode >= 300 {
		// we don't handle redirects yet
		// also this will catch other client request errors (>= 400)
		return nil, (*ResponseStatusCodeError)(resp)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if string(body) == "null" {
		return nil, nil
	}
	var releases []releaseResponse
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}

type releaseResponse struct {
	Version string `json:"version"`
	SHA     string `json:"sha1"`
}

func (src *BOSHIOReleaseSource) releaseExistOnBoshIO(ctx context.Context, name, version string) (bool, error) {
	releaseResponses, err := src.getReleases(ctx, name)
	if err != nil {
		return false, err
	}
	for _, rel := range releaseResponses {
		if rel.Version == version {
			return true, nil
		}
	}
	return false, nil
}

func (src *BOSHIOReleaseSource) serverURI() string {
	if src.CustomURI != "" {
		return src.CustomURI
	}
	return "https://bosh.io"
}

var organizations = []string{
	"cloudfoundry",
	"pivotal-cf",
	"cloudfoundry-incubator",
	"pivotal-cf-experimental",
	"bosh-packages",
	"cppforlife",
	"vito",
	"flavorjones",
	"xoebus",
	"dpb587",
	"jamlo",
	"concourse",
	"cf-platform-eng",
	"starkandwayne",
	"cloudfoundry-community",
	"vmware",
	"DataDog",
	"Dynatrace",
	"SAP",
	"hybris",
	"minio",
	"rakutentech",
	"frodenas",
}

var suffixes = []string{
	"-release",
	"-boshrelease",
	"-bosh-release",
	"",
}
