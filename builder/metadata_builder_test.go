package builder_test

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf/kiln/builder"
	"github.com/pivotal-cf/kiln/builder/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MetadataBuilder", func() {
	var (
		releaseManifestReader  *fakes.ReleaseManifestReader
		stemcellManifestReader *fakes.StemcellManifestReader
		handcraftReader        *fakes.HandcraftReader
		logger                 *fakes.Logger
		tileBuilder            builder.MetadataBuilder
	)

	BeforeEach(func() {
		releaseManifestReader = &fakes.ReleaseManifestReader{}
		stemcellManifestReader = &fakes.StemcellManifestReader{}
		handcraftReader = &fakes.HandcraftReader{}
		logger = &fakes.Logger{}

		tileBuilder = builder.NewMetadataBuilder(releaseManifestReader, stemcellManifestReader, handcraftReader, logger)
	})

	Describe("Build", func() {
		It("creates a metadata with the correct information", func() {
			releaseManifestReader.ReadCall.Stub = func(path string) (builder.ReleaseManifest, error) {
				switch path {
				case "/tmp/release-1.tgz":
					return builder.ReleaseManifest{
						Name:    "release-1",
						Version: "version-1",
					}, nil
				case "/tmp/release-2.tgz":
					return builder.ReleaseManifest{
						Name:    "release-2",
						Version: "version-2",
					}, nil
				default:
					return builder.ReleaseManifest{}, fmt.Errorf("could not read release %q", path)
				}
			}
			stemcellManifestReader.ReadCall.Returns.StemcellManifest = builder.StemcellManifest{
				Version:         "2332",
				OperatingSystem: "ubuntu-trusty",
			}

			handcraftReader.ReadCall.Returns.Handcraft = builder.Handcraft{
				"metadata_version":          "some-metadata-version",
				"provides_product_versions": "some-provides-product-versions",
			}

			metadata, err := tileBuilder.Build([]string{
				"/tmp/release-1.tgz",
				"/tmp/release-2.tgz",
			}, "/tmp/test-stemcell.tgz", "/some/path/handcraft.yml", "cool-product", "1.2.3")
			Expect(err).NotTo(HaveOccurred())
			Expect(stemcellManifestReader.ReadCall.Receives.Path).To(Equal("/tmp/test-stemcell.tgz"))
			Expect(handcraftReader.ReadCall.Receives.Path).To(Equal("/some/path/handcraft.yml"))

			Expect(metadata.Name).To(Equal("cool-product"))
			Expect(metadata.Releases).To(Equal([]builder.MetadataRelease{
				{
					Name:    "release-1",
					Version: "version-1",
					File:    "release-1.tgz",
				},
				{
					Name:    "release-2",
					Version: "version-2",
					File:    "release-2.tgz",
				},
			}))
			Expect(metadata.StemcellCriteria).To(Equal(builder.MetadataStemcellCriteria{
				Version:     "2332",
				OS:          "ubuntu-trusty",
				RequiresCPI: false,
			}))
			Expect(metadata.Handcraft).To(Equal(builder.Handcraft{
				"metadata_version":          "some-metadata-version",
				"provides_product_versions": "some-provides-product-versions",
			}))

			Expect(logger.PrintfCall.Receives.LogLines).To(Equal([]string{
				"Creating metadata for .pivotal...",
				"Read manifest for release release-1",
				"Read manifest for release release-2",
				"Read manifest for stemcell version 2332",
				"Read handcraft",
			}))
		})

		Context("failure cases", func() {
			Context("when the release tarball cannot be read", func() {
				It("returns an error", func() {
					releaseManifestReader.ReadCall.Returns.Error = errors.New("failed to read release tarball")

					_, err := tileBuilder.Build([]string{"release-1.tgz"}, "", "", "", "")
					Expect(err).To(MatchError("failed to read release tarball"))
				})
			})

			Context("when the stemcell tarball cannot be read", func() {
				It("returns an error", func() {
					stemcellManifestReader.ReadCall.Returns.Error = errors.New("failed to read stemcell tarball")

					_, err := tileBuilder.Build([]string{}, "stemcell.tgz", "", "", "")
					Expect(err).To(MatchError("failed to read stemcell tarball"))
				})
			})

			Context("when the handcraft cannot be read", func() {
				It("returns an error", func() {
					handcraftReader.ReadCall.Returns.Error = errors.New("failed to read handcraft")

					_, err := tileBuilder.Build([]string{}, "", "handcraft.yml", "", "")
					Expect(err).To(MatchError("failed to read handcraft"))
				})
			})
		})
	})
})