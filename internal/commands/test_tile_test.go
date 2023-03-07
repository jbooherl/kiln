package commands_test

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/pivotal-cf/kiln/internal/commands"
	commandsFakes "github.com/pivotal-cf/kiln/internal/commands/fakes"
)

var _ = FDescribe("kiln test docker", func() {
	var (
		fakeTasDirectory string
	)
	BeforeEach(func() {
		fakeTasDirectory = setupFakeTAS(GinkgoT())
	})

	Context("locally missing docker image is built", func() {
		var (
			ctx    context.Context
			writer strings.Builder
			logger *log.Logger
		)

		BeforeEach(func() {
			ctx = context.Background()
			logger = log.New(&writer, "", 0)
		})

		It("runs the tests for ist", func() {
			oldStdin := os.Stdin
			format.MaxLength = 100000
			passwd := "password\n"
			content := []byte(passwd)
			tmpfile, err := os.CreateTemp(GinkgoT().TempDir(), GinkgoT().Name())
			Expect(err).To(BeNil())
			_, err = tmpfile.Write(content)
			Expect(err).To(BeNil())
			_, err = tmpfile.Seek(0, 0)
			Expect(err).To(BeNil())
			os.Stdin = tmpfile

			fakeMobyClient := &commandsFakes.MobyClient{}
			fakeMobyClient.PingReturns(types.Ping{}, nil)
			r, w := io.Pipe()
			_ = w.Close()
			fakeConn := &fakeConn{r: r, w: stdcopy.NewStdWriter(io.Discard, stdcopy.Stdout)}
			fakeMobyClient.DialHijackReturns(fakeConn, nil)
			rc := io.NopCloser(strings.NewReader(`{"error": "", "message": "tagged dont_push_me_vmware_confidential:123"}`))
			imageBuildResponse := types.ImageBuildResponse{
				Body: rc,
			}
			fakeMobyClient.ImageBuildReturns(imageBuildResponse, nil)
			createResp := container.CreateResponse{
				ID: "some id",
			}
			fakeMobyClient.ContainerCreateReturns(createResp, nil)
			responses := make(chan container.WaitResponse)
			go func() {
				responses <- container.WaitResponse{
					Error:      nil,
					StatusCode: 0,
				}
			}()
			fakeMobyClient.ContainerWaitReturns(responses, nil)
			rcLog := io.NopCloser(strings.NewReader(`{"error": "", "message": "manifest tests completed successfully"}"`))
			fakeMobyClient.ContainerLogsReturns(rcLog, nil)
			fakeMobyClient.ContainerStartReturns(nil)
			fakeSshThinger := commandsFakes.SshProvider{}
			fakeSshThinger.NeedsKeysReturns(false, nil)
			subjectUnderTest := commands.NewManifestTest(logger, ctx, fakeMobyClient, &fakeSshThinger)

			err = subjectUnderTest.Execute([]string{"--tile-path", filepath.Join(fakeTasDirectory, "ist"), "--ginkgo-manifest-flags", "-r -slowSpecThreshold 1"})
			Expect(err).To(BeNil())

			_, config, _, _, _, _ := fakeMobyClient.ContainerCreateArgsForCall(0)
			Expect(len(config.Env)).To(Equal(2))
			Expect(config.Env[0]).To(Equal("TAS_METADATA_PATH=/tas/ist/test/manifest/fixtures/tas_metadata.yml"))
			Expect(config.Env[1]).To(Equal("TAS_CONFIG_FILE=/tas/ist/test/manifest/fixtures/tas_config.yml"))

			dockerCmd := fmt.Sprintf("cd /tas/%s/test/manifest && PRODUCT=ist RENDERER=ops-manifest ginkgo %s", "ist", "-r -slowSpecThreshold 1")
			Expect(config.Cmd).To(Equal(strslice.StrSlice{"/bin/bash", "-c", dockerCmd}))
			GinkgoT().Log(writer.String())
			Expect((&writer).String()).To(ContainSubstring("tagged dont_push_me_vmware_confidential:123"))
			Expect((&writer).String()).To(ContainSubstring("Building / restoring cached docker image"))
			os.Stdin = oldStdin
		})

		It("runs the tests for tas", func() {
			oldStdin := os.Stdin
			format.MaxLength = 100000
			passwd := "password\n"
			content := []byte(passwd)
			tmpfile, err := os.CreateTemp(GinkgoT().TempDir(), GinkgoT().Name())
			Expect(err).To(BeNil())
			_, err = tmpfile.Write(content)
			Expect(err).To(BeNil())
			_, err = tmpfile.Seek(0, 0)
			Expect(err).To(BeNil())
			os.Stdin = tmpfile

			fakeMobyClient := &commandsFakes.MobyClient{}
			fakeMobyClient.PingReturns(types.Ping{}, nil)
			r, w := io.Pipe()
			_ = w.Close()
			fakeConn := &fakeConn{r: r, w: stdcopy.NewStdWriter(io.Discard, stdcopy.Stdout)}
			fakeMobyClient.DialHijackReturns(fakeConn, nil)

			rc := io.NopCloser(strings.NewReader(`{"error": "", "message": "tagged dont_push_me_vmware_confidential:123"}`))
			imageBuildResponse := types.ImageBuildResponse{
				Body: rc,
			}
			fakeMobyClient.ImageBuildReturns(imageBuildResponse, nil)
			createResp := container.CreateResponse{
				ID: "some id",
			}
			fakeMobyClient.ContainerCreateReturns(createResp, nil)
			responses := make(chan container.WaitResponse)
			go func() {
				responses <- container.WaitResponse{
					Error:      nil,
					StatusCode: 0,
				}
			}()
			fakeMobyClient.ContainerWaitReturns(responses, nil)
			rcLog := io.NopCloser(strings.NewReader(`{"error": "", "message": "manifest tests completed successfully"}"`))
			fakeMobyClient.ContainerLogsReturns(rcLog, nil)
			fakeMobyClient.ContainerStartReturns(nil)
			fakeSshThinger := commandsFakes.SshProvider{}
			fakeSshThinger.NeedsKeysReturns(false, nil)
			subjectUnderTest := commands.NewManifestTest(logger, ctx, fakeMobyClient, &fakeSshThinger)
			err = subjectUnderTest.Execute([]string{"--tile-path", filepath.Join(fakeTasDirectory, "tas")})

			Expect(err).To(BeNil())
			_, config, _, _, _, _ := fakeMobyClient.ContainerCreateArgsForCall(0)
			Expect(len(config.Env)).To(Equal(0))

			GinkgoT().Log(writer.String())
			Expect((&writer).String()).To(ContainSubstring("tagged dont_push_me_vmware_confidential:123"))
			Expect((&writer).String()).To(ContainSubstring("Building / restoring cached docker image"))
			os.Stdin = oldStdin
		})

		It("exits with an error if docker isn't running", func() {
			fakeMobyClient := &commandsFakes.MobyClient{}
			fakeMobyClient.PingReturns(types.Ping{}, errors.New("docker not running"))
			fakeSshThinger := commandsFakes.SshProvider{}
			fakeSshThinger.NeedsKeysReturns(false, nil)
			subjectUnderTest := commands.NewManifestTest(logger, ctx, fakeMobyClient, &fakeSshThinger)
			err := subjectUnderTest.Execute([]string{"tas_fake"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Docker daemon is not running"))
		})
	})

	Context("logs output from container", func() {
		It("logs when the manifest tests complete successfully", func() {
			fakeMobyClient := &commandsFakes.MobyClient{}

			ctx := context.Background()
			fakeMobyClient.PingReturns(types.Ping{}, nil)
			r, w := io.Pipe()
			_ = w.Close()
			fakeConn := &fakeConn{r: r, w: stdcopy.NewStdWriter(io.Discard, stdcopy.Stdout)}
			fakeMobyClient.DialHijackReturns(fakeConn, nil)
			output := "Successfully built 1234"
			rc := io.NopCloser(strings.NewReader(fmt.Sprintf(`{"error": "", "message": "%s"}`, output)))
			imageBuildResponse := types.ImageBuildResponse{
				Body: rc,
			}
			fakeMobyClient.ImageBuildReturns(imageBuildResponse, nil)
			createResp := container.CreateResponse{
				ID: "some id",
			}
			fakeMobyClient.ContainerCreateReturns(createResp, nil)
			var logLogOutput bytes.Buffer
			logOut := log.New(&logLogOutput, "", 0)

			rcLog := io.NopCloser(strings.NewReader(`"manifest tests completed successfully"`))
			fakeMobyClient.ContainerLogsReturns(rcLog, nil)
			fakeMobyClient.ContainerStartReturns(nil)
			fakeSshThinger := commandsFakes.SshProvider{}
			fakeSshThinger.NeedsKeysReturns(false, nil)
			subjectUnderTest := commands.NewManifestTest(logOut, ctx, fakeMobyClient, &fakeSshThinger)
			fakeMobyClient.ContainerCreateReturns(createResp, nil)
			responses := make(chan container.WaitResponse)
			go func() {
				responses <- container.WaitResponse{
					Error:      nil,
					StatusCode: 0,
				}
			}()
			fakeMobyClient.ContainerWaitReturns(responses, nil)

			err := subjectUnderTest.Execute([]string{"--tile-path", filepath.Join(fakeTasDirectory, "tas")})
			Expect(err).To(BeNil())

			Expect(logLogOutput.String()).To(ContainSubstring("manifest tests completed successfully"))
		})
	})
})

type fakeConn struct {
	r io.Reader
	w io.Writer
}

func (c *fakeConn) LocalAddr() net.Addr {
	return nil
}

func (c *fakeConn) RemoteAddr() net.Addr {
	return &net.UnixAddr{Name: "test", Net: "test"}
}

func (c *fakeConn) SetDeadline(time.Time) error {
	return nil
}

func (c *fakeConn) SetReadDeadline(time.Time) error {
	return nil
}

func (c *fakeConn) SetWriteDeadline(time.Time) error {
	return nil
}

func (c *fakeConn) Read(p []byte) (int, error) {
	return c.r.Read(p)
}

func (c *fakeConn) Write(p []byte) (int, error) {
	return c.w.Write(p)
}

func (c *fakeConn) Close() error {
	return nil
}

func setupFakeTAS(t GinkgoTInterface) string {
	fakeTasDirectory := os.TempDir() // t.TempDir is making a dir in the current directory
	t.Cleanup(func() {
		_ = os.RemoveAll(fakeTasDirectory)
	})
	tasFakeZipFile, err := os.Open(filepath.Join("testdata", "tas_fake.zip"))
	if err != nil {
		t.Fatal(err)
	}
	zipStat, err := tasFakeZipFile.Stat()
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(tasFakeZipFile, zipStat.Size())
	if err != nil {
		t.Fatal(err)
	}
	walkErr := fs.WalkDir(zr, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		filePathDirSegments := strings.Split(path.Dir(filePath), "/")
		err = os.MkdirAll(filepath.Join(append([]string{fakeTasDirectory}, filePathDirSegments...)...), 0755)
		if err != nil {
			return err
		}
		filePathSegments := strings.Split(filePath, "/")
		f, err := os.Create(filepath.Join(append([]string{fakeTasDirectory}, filePathSegments...)...))
		if err != nil {
			return err
		}
		defer closeAndIgnoreError(f)
		zf, err := zr.Open(filePath)
		if err != nil {
			return err
		}
		defer closeAndIgnoreError(zf)
		_, err = io.Copy(f, zf)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(walkErr)
	}
	return fakeTasDirectory
}
