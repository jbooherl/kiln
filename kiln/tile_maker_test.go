package kiln_test

import (
	"errors"

	"github.com/pivotal-cf/kiln/builder"
	"github.com/pivotal-cf/kiln/kiln"
	"github.com/pivotal-cf/kiln/kiln/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TileMaker", func() {
	var (
		fakeMetadataBuilder *fakes.MetadataBuilder
		fakeTileWriter      *fakes.TileWriter
		fakeLogger          *fakes.Logger

		config    kiln.ApplicationConfig
		tileMaker kiln.TileMaker
	)

	BeforeEach(func() {
		fakeMetadataBuilder = &fakes.MetadataBuilder{}
		fakeTileWriter = &fakes.TileWriter{}
		fakeLogger = &fakes.Logger{}

		config = kiln.ApplicationConfig{
			Name:                 "cool-product",
			Version:              "1.2.3-build.4",
			FinalVersion:         "1.2.3",
			StemcellTarball:      "some-stemcell-tarball",
			ReleaseTarballs:      []string{"some-release-tarball", "some-other-release-tarball"},
			Handcraft:            "some-handcraft",
			Migrations:           []string{"some-migration", "some-other-migration"},
			BaseContentMigration: "some-base-content-migration",
			ContentMigrations:    []string{"some-content-migration", "some-other-content-migration"},
			OutputDir:            "some-output-dir",
			StubReleases:         true,
		}
		tileMaker = kiln.NewTileMaker(fakeMetadataBuilder, fakeTileWriter, fakeLogger)
	})

	It("builds the metadata", func() {
		err := tileMaker.Make(config)
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeMetadataBuilder.BuildCallCount()).To(Equal(1))

		releaseTarballs, stemcellTarball, handcraft, name, version := fakeMetadataBuilder.BuildArgsForCall(0)
		Expect(releaseTarballs).To(Equal([]string{"some-release-tarball", "some-other-release-tarball"}))
		Expect(stemcellTarball).To(Equal("some-stemcell-tarball"))
		Expect(handcraft).To(Equal("some-handcraft"))
		Expect(name).To(Equal("cool-product"))
		Expect(version).To(Equal("1.2.3"))
	})

	It("makes the tile", func() {
		fakeMetadataBuilder.BuildReturns(builder.Metadata{
			Name: "cool-product",
			Releases: []builder.MetadataRelease{{
				Name:    "some-release",
				File:    "some-release-tarball",
				Version: "1.2.3-build.4",
			}},
			StemcellCriteria: builder.MetadataStemcellCriteria{
				Version:     "2.3.4",
				OS:          "an-operating-system",
				RequiresCPI: false,
			},
		}, nil)

		err := tileMaker.Make(config)
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeTileWriter.WriteCallCount()).To(Equal(1))

		metadataContents, writeCfg := fakeTileWriter.WriteArgsForCall(0)
		Expect(metadataContents).To(MatchYAML(`
name: cool-product
releases:
- name: some-release
  file: some-release-tarball
  version: 1.2.3-build.4
stemcell_criteria:
  version: 2.3.4
  os: an-operating-system
  requires_cpi: false`))
		Expect(writeCfg).To(Equal(builder.WriteConfig{
			Name:                 "cool-product",
			Version:              "1.2.3-build.4",
			FinalVersion:         "1.2.3",
			StemcellTarball:      "some-stemcell-tarball",
			ReleaseTarballs:      []string{"some-release-tarball", "some-other-release-tarball"},
			Handcraft:            "some-handcraft",
			Migrations:           []string{"some-migration", "some-other-migration"},
			BaseContentMigration: "some-base-content-migration",
			ContentMigrations:    []string{"some-content-migration", "some-other-content-migration"},
			OutputDir:            "some-output-dir",
			StubReleases:         true,
		}))
	})

	It("logs its step", func() {
		err := tileMaker.Make(config)
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeLogger.PrintlnCall.Receives.LogLines).To(Equal([]string{"Marshaling metadata file..."}))
	})

	Context("failure cases", func() {
		Context("when metadata builder fails", func() {
			It("returns an error", func() {
				fakeMetadataBuilder.BuildReturns(builder.Metadata{}, errors.New("some-error"))

				err := tileMaker.Make(config)
				Expect(err).To(MatchError("some-error"))
			})
		})

		Context("when the tile writer fails", func() {
			It("returns an error", func() {
				fakeTileWriter.WriteReturns(errors.New("tile writer has failed"))

				err := tileMaker.Make(config)
				Expect(err).To(MatchError("tile writer has failed"))
			})
		})
	})
})