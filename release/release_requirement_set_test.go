package release_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/kiln/internal/cargo"
	. "github.com/pivotal-cf/kiln/release"
)

var _ = Describe("ReleaseRequirementSet", func() {
	const (
		release1Name    = "release-1"
		release1Version = "1.2.3"
		release2Name    = "release-2"
		release2Version = "2.3.4"
		stemcellName    = "some-os"
		stemcellVersion = "9.8.7"
	)

	var (
		rrs                    ReleaseRequirementSet
		release1ID, release2ID ReleaseID
	)

	BeforeEach(func() {
		kilnfileLock := cargo.KilnfileLock{
			Releases: []cargo.ReleaseLock{
				{Name: release1Name, Version: release1Version},
				{Name: release2Name, Version: release2Version},
			},
			Stemcell: cargo.Stemcell{OS: stemcellName, Version: stemcellVersion},
		}
		rrs = NewReleaseRequirementSet(kilnfileLock)
		release1ID = ReleaseID{Name: release1Name, Version: release1Version}
		release2ID = ReleaseID{Name: release2Name, Version: release2Version}
	})

	Describe("NewReleaseRequirementSet", func() {
		It("constructs a requirement set based on the Kilnfile.lock", func() {
			Expect(rrs).To(HaveLen(2))
			Expect(rrs).To(HaveKeyWithValue(release1ID,
				ReleaseRequirement{Name: release1Name, Version: release1Version, StemcellOS: stemcellName, StemcellVersion: stemcellVersion},
			))
			Expect(rrs).To(HaveKeyWithValue(release2ID,
				ReleaseRequirement{Name: release2Name, Version: release2Version, StemcellOS: stemcellName, StemcellVersion: stemcellVersion},
			))
		})
	})

	Describe("Partition", func() {
		var (
			releaseSet                                           []SatisfyingLocalRelease
			extraReleaseID                                       ReleaseID
			satisfyingRelease, unsatisfyingRelease, extraRelease SatisfyingLocalRelease
		)

		BeforeEach(func() {
			satisfyingRelease = NewLocalBuiltRelease(release1ID, "satisfying-path")
			unsatisfyingRelease = NewLocalBuiltRelease(ReleaseID{Name: release2Name, Version: "4.0.4"}, "unsatisfying-path")

			extraReleaseID = ReleaseID{Name: "extra", Version: "2.3.5"}
			extraRelease = NewLocalBuiltRelease(extraReleaseID, "so-extra")

			releaseSet = []SatisfyingLocalRelease{satisfyingRelease, unsatisfyingRelease, extraRelease}
		})

		It("returns the intersecting, missing, and extra releases", func() {
			intersection, missing, extra := rrs.Partition(releaseSet)

			Expect(intersection).To(HaveLen(1))
			Expect(intersection).To(ConsistOf(
				LocalRelease{ReleaseID: release1ID, LocalPath: satisfyingRelease.LocalPath},
			))

			Expect(missing).To(HaveLen(1))
			Expect(missing).To(HaveKeyWithValue(release2ID, rrs[release2ID]))

			Expect(extra).To(HaveLen(2))
			Expect(extra).To(ConsistOf(
				LocalRelease{
					ReleaseID: ReleaseID{Name: release2Name, Version: "4.0.4"},
					LocalPath: unsatisfyingRelease.LocalPath,
				},
				LocalRelease{ReleaseID: extraReleaseID, LocalPath: extraRelease.LocalPath},
			))
		})

		It("does not modify itself", func() {
			rrs.Partition(releaseSet)
			Expect(rrs).To(HaveLen(2))
			Expect(rrs).To(HaveKey(release1ID))
			Expect(rrs).To(HaveKey(release2ID))
		})

		It("does not modify the given release set", func() {
			rrs.Partition(releaseSet)
			Expect(releaseSet).To(ConsistOf(satisfyingRelease, unsatisfyingRelease, extraRelease))
		})
	})
})