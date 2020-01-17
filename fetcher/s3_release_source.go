package fetcher

import (
	"bytes"
	"fmt"
	"github.com/pivotal-cf/kiln/release"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pivotal-cf/kiln/internal/cargo"
)

//go:generate counterfeiter -o ./fakes/s3_downloader.go --fake-name S3Downloader . S3Downloader
type S3Downloader interface {
	Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error)
}

//go:generate counterfeiter -o ./fakes/s3_uploader.go --fake-name S3Uploader . S3Uploader
type S3Uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

//go:generate counterfeiter -o ./fakes/s3_head_objecter.go --fake-name S3HeadObjecter . S3HeadObjecter
type S3HeadObjecter interface {
	HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
}

type S3ReleaseSource struct {
	Logger       *log.Logger
	S3Client     S3HeadObjecter
	S3Downloader S3Downloader
	S3Uploader   S3Uploader
	Bucket       string
	PathTemplate string
}

func (src *S3ReleaseSource) Configure(config cargo.ReleaseSourceConfig) {
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyId,
			config.SecretAccessKey,
			"",
		),
	}))
	client := s3.New(sess)

	src.S3Client = client
	src.S3Downloader = s3manager.NewDownloaderWithClient(client)
	src.S3Uploader = s3manager.NewUploaderWithClient(client)

	src.Bucket = config.Bucket
	src.PathTemplate = config.PathTemplate
}

func (src S3ReleaseSource) ID() string {
	return src.Bucket
}

//go:generate counterfeiter -o ./fakes/s3_request_failure.go --fake-name S3RequestFailure github.com/aws/aws-sdk-go/service/s3.RequestFailure
func (src S3ReleaseSource) GetMatchedRelease(requirement release.ReleaseRequirement) (release.RemoteRelease, bool, error) {
	t := src.pathTemplate()

	pathBuf := new(bytes.Buffer)
	err := t.Execute(pathBuf, requirement)
	if err != nil {
		return release.RemoteRelease{}, false, fmt.Errorf("unable to evaluate path_template: %w", err)
	}

	remotePath := pathBuf.String()

	headRequest := new(s3.HeadObjectInput)
	headRequest.SetBucket(src.Bucket)
	headRequest.SetKey(remotePath)

	_, err = src.S3Client.HeadObject(headRequest)
	if err != nil {
		requestFailure, ok := err.(s3.RequestFailure)
		if ok && requestFailure.StatusCode() == 404 {
			return release.RemoteRelease{}, false, nil
		}
		return release.RemoteRelease{}, false, err
	}

	return release.RemoteRelease{
		ReleaseID:  release.ReleaseID{Name: requirement.Name, Version: requirement.Version},
		RemotePath: remotePath,
	}, true, nil
}

func (src S3ReleaseSource) DownloadRelease(releaseDir string, remoteRelease release.RemoteRelease, downloadThreads int) (release.LocalRelease, error) {
	setConcurrency := func(dl *s3manager.Downloader) {
		if downloadThreads > 0 {
			dl.Concurrency = downloadThreads
		} else {
			dl.Concurrency = s3manager.DefaultDownloadConcurrency
		}
	}

	src.Logger.Printf("downloading %s %s from %s", remoteRelease.Name, remoteRelease.Version, src.Bucket)

	outputFile := filepath.Join(releaseDir, fmt.Sprintf("%s-%s.tgz", remoteRelease.Name, remoteRelease.Version))

	file, err := os.Create(outputFile)
	if err != nil {
		return release.LocalRelease{}, fmt.Errorf("failed to create file %q: %w", outputFile, err)
	}

	_, err = src.S3Downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(src.Bucket),
		Key:    aws.String(remoteRelease.RemotePath),
	}, setConcurrency)
	if err != nil {
		return release.LocalRelease{}, fmt.Errorf("failed to download file: %w\n", err)
	}

	return release.LocalRelease{ReleaseID: remoteRelease.ReleaseID, LocalPath: outputFile}, nil
}

func (src S3ReleaseSource) UploadRelease(name, version string, file io.Reader) error {
	pathBuf := new(bytes.Buffer)
	err := src.pathTemplate().Execute(pathBuf, release.ReleaseID{Name: name, Version: version})
	if err != nil {
		return fmt.Errorf("unable to evaluate path_template: %w", err)
	}

	remotePath := pathBuf.String()

	src.Logger.Printf("Uploading release to %s at %q...\n", src.ID(), remotePath)

	_, err = src.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(src.Bucket),
		Key:    aws.String(remotePath),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func (src S3ReleaseSource) pathTemplate() *template.Template {
	return template.Must(
		template.New("remote-path").
			Funcs(template.FuncMap{"trimSuffix": strings.TrimSuffix}).
			Parse(src.PathTemplate))
}