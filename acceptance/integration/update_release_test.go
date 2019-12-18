package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pivotal-cf/kiln/internal/cargo"
	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Context("Updating a release to a specific version", func() {
	var kilnfileContents, previousKilnfileLock, kilnfileLockPath, kilnfilePath, releasesPath string

	Context("for public releases", func() {
		BeforeEach(func() {
			kilnfileContents = `---
release_sources:
- type: bosh.io

`
			previousKilnfileLock = `---
releases:
- name: "loggregator-agent"
  version: "5.1.0"
  sha1: "a86e10219b0ed9b7b82f0610b7cdc03c13765722"
- name: capi
  sha1: "03ac801323cd23205dde357cc7d2dc9e92bc0c93"
  version: "1.87.0"
stemcell_criteria:
  os: some-os
  version: "4.5.6"
`
			tmpDir, err := ioutil.TempDir("", "kiln-main-test")
			Expect(err).NotTo(HaveOccurred())

			kilnfileLockPath = filepath.Join(tmpDir, "Kilnfile.lock")
			kilnfilePath = filepath.Join(tmpDir, "Kilnfile")
			releasesPath = filepath.Join(tmpDir, "releases")
			ioutil.WriteFile(kilnfilePath, []byte(kilnfileContents), 0600)
			ioutil.WriteFile(kilnfileLockPath, []byte(previousKilnfileLock), 0600)
			os.Mkdir(releasesPath, 0700)
		})

		It("updates the Kilnfile.lock", func() {
			command := exec.Command(pathToMain, "update-release",
				"--name", "capi",
				"--version", "1.88.0",
				"--kilnfile", kilnfilePath,
				"--releases-directory", releasesPath)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 60*time.Second).Should(gexec.Exit(0))
			Expect(session.Out).To(gbytes.Say("Updated capi to 1.88.0"))

			var kilnfileLock cargo.KilnfileLock

			file, err := os.Open(kilnfileLockPath)
			Expect(err).NotTo(HaveOccurred())

			err = yaml.NewDecoder(file).Decode(&kilnfileLock)
			Expect(err).NotTo(HaveOccurred())

			Expect(kilnfileLock).To(Equal(
				cargo.KilnfileLock{
					Releases: []cargo.Release{
						{Name: "loggregator-agent", Version: "5.1.0", SHA1: "a86e10219b0ed9b7b82f0610b7cdc03c13765722"},
						{Name: "capi", Version: "1.88.0", SHA1: "7a7ef183de3252724b6f8e6ca39ad7cf4995fe27"},
					},
					Stemcell: cargo.Stemcell{
						OS:      "some-os",
						Version: "4.5.6",
					},
				}))
		})
	})

	Context("for private releases (on S3)", func() {
		BeforeEach(func() {
			kilnfileContents = `---
release_sources:
- type: s3
  bucket: compiled-releases
  region: us-west-1
  access_key_id: $(variable "aws_access_key_id")
  secret_access_key: $(variable "aws_secret_access_key")
  regex: ^2\.8/.+/(?P<release_name>[a-z-_0-9]+)-(?P<release_version>v?[0-9\.]+)-(?P<stemcell_os>[a-z-_]+)-(?P<stemcell_version>\d+\.\d+)(\.0)?\.tgz$
  compiled: true
  publishable: true
`
			previousKilnfileLock = `---
releases:
- name: "loggregator-agent"
  version: "5.1.0"
  sha1: "a86e10219b0ed9b7b82f0610b7cdc03c13765722"
- name: capi
  sha1: "03ac801323cd23205dde357cc7d2dc9e92bc0c93"
  version: "1.87.0"
stemcell_criteria:
  os: ubuntu-xenial
  version: '456.30'
`
			tmpDir, err := ioutil.TempDir("", "kiln-main-test")
			Expect(err).NotTo(HaveOccurred())

			kilnfileLockPath = filepath.Join(tmpDir, "Kilnfile.lock")
			kilnfilePath = filepath.Join(tmpDir, "Kilnfile")
			releasesPath = filepath.Join(tmpDir, "releases")

			_ = ioutil.WriteFile(kilnfilePath, []byte(kilnfileContents), 0600)
			_ = ioutil.WriteFile(kilnfileLockPath, []byte(previousKilnfileLock), 0600)
			_ = os.Mkdir(releasesPath, 0700)
		})

		It("updates the Kilnfile.lock", func() {
			varsFile := os.Getenv("KILN_ACCEPTANCE_VARS_FILE_CONTENTS")
			if varsFile == "" {
				Fail("please provide the KILN_ACCEPTANCE_VARS_FILE_CONTENTS environment variable")
			}

			tmpfile, err := ioutil.TempFile("", "varsfile")
			Expect(err).NotTo(HaveOccurred())

			tmpfile.Write([]byte(varsFile))
			tmpfile.Close()

			command := exec.Command(pathToMain, "update-release",
				"--name", "capi",
				"--version", "1.86.0",
				"--kilnfile", kilnfilePath,
				"--releases-directory", releasesPath,
				"--variables-file", tmpfile.Name(),
			)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 60*time.Second).Should(gexec.Exit(0))
			Expect(session.Out).To(gbytes.Say("Updated capi to 1.86.0"))

			var kilnfileLock cargo.KilnfileLock

			file, err := os.Open(kilnfileLockPath)
			Expect(err).NotTo(HaveOccurred())

			err = yaml.NewDecoder(file).Decode(&kilnfileLock)
			Expect(err).NotTo(HaveOccurred())

			Expect(kilnfileLock).To(Equal(
				cargo.KilnfileLock{
					Releases: []cargo.Release{
						{Name: "loggregator-agent", Version: "5.1.0", SHA1: "a86e10219b0ed9b7b82f0610b7cdc03c13765722"},
						{Name: "capi", Version: "1.86.0", SHA1: "32f40c3006e3b0b401b855da99cbd701c3c5be33"},
					},
					Stemcell: cargo.Stemcell{
						OS:      "ubuntu-xenial",
						Version: "456.30",
					},
				}))
		})
	})
})