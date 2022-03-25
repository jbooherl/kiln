package component_test

import (
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/kiln/internal/component"
	"github.com/pivotal-cf/kiln/pkg/cargo"
)

var _ = Describe("ReleaseSourceList", func() {
	var logger *log.Logger

	BeforeEach(func() {
		logger = log.New(GinkgoWriter, "", log.LstdFlags)
	})

	Describe("NewReleaseSourceRepo", func() {
		var kilnfile cargo.Kilnfile

		Context("happy path", func() {
			BeforeEach(func() {
				kilnfile = cargo.Kilnfile{
					ReleaseSources: cargo.ReleaseSourceList{
						cargo.S3ReleaseSource{Bucket: "compiled-releases", Region: "us-west-1", PathTemplate: "template", Publishable: true},
						cargo.S3ReleaseSource{Bucket: "built-releases", Region: "us-west-1", PathTemplate: "template", Publishable: false},
						cargo.BOSHIOReleaseSource{Publishable: false},
						cargo.GitHubReleaseSource{Org: "cloudfoundry", GithubToken: "banana"},
					},
				}
			})

			It("constructs all the release sources", func() {
				releaseSources := component.NewReleaseSourceRepo(kilnfile, logger)
				Expect(len(releaseSources)).To(Equal(4)) // not using HaveLen because S3 struct is so huge
			})

			It("constructs the compiled release source properly", func() {
				releaseSources := component.NewReleaseSourceRepo(kilnfile, logger)
				Expect(releaseSources[0]).To(BeAssignableToTypeOf(component.S3ReleaseSource{}))
			})

			It("sets the release source id to bucket id for s3", func() {
				releaseSources := component.NewReleaseSourceRepo(kilnfile, logger)

				Expect(releaseSources[0].ID()).To(Equal(kilnfile.ReleaseSources[0].(cargo.S3ReleaseSource).Bucket))
				Expect(releaseSources[1].ID()).To(Equal(kilnfile.ReleaseSources[1].(cargo.S3ReleaseSource).Bucket))
			})

			It("constructs the built release source properly", func() {
				releaseSources := component.NewReleaseSourceRepo(kilnfile, logger)

				Expect(releaseSources[1]).To(BeAssignableToTypeOf(component.S3ReleaseSource{}))
				Expect(releaseSources[1].ID()).To(Equal(kilnfile.ReleaseSources[1].(cargo.S3ReleaseSource).Bucket))
			})

			It("constructs the github release source properly", func() {
				releaseSources := component.NewReleaseSourceRepo(kilnfile, logger)

				Expect(releaseSources[3]).To(BeAssignableToTypeOf(&component.GithubReleaseSource{}))
				Expect(releaseSources[3].ID()).To(Equal(kilnfile.ReleaseSources[3].(cargo.GitHubReleaseSource).Org))
			})
		})

		Context("when bosh.io is publishable", func() {
			BeforeEach(func() {
				kilnfile = cargo.Kilnfile{
					ReleaseSources: cargo.ReleaseSourceList{
						cargo.BOSHIOReleaseSource{Publishable: true},
					},
				}
			})

			It("marks it correctly", func() {
				releaseSources := component.NewReleaseSourceRepo(kilnfile, logger)

				Expect(releaseSources).To(HaveLen(1))
				var (
					boshIOReleaseSource *component.BOSHIOReleaseSource
				)

				Expect(releaseSources[0]).To(BeAssignableToTypeOf(boshIOReleaseSource))
				Expect(releaseSources[0].IsPublishable()).To(BeTrue())
			})
		})

		Context("when the Kilnfile gives explicit IDs", func() {
			BeforeEach(func() {
				kilnfile = cargo.Kilnfile{
					ReleaseSources: cargo.ReleaseSourceList{
						cargo.S3ReleaseSource{Identifier: "comp", Publishable: true, Bucket: "compiled-releases", Region: "us-west-1", PathTemplate: "template"},
						cargo.S3ReleaseSource{Identifier: "buil", Publishable: false, Bucket: "built-releases", Region: "us-west-1", PathTemplate: "template"},
						cargo.BOSHIOReleaseSource{Identifier: "bosh", Publishable: false},
					},
				}
			})

			It("gives the correct IDs to the release sources", func() {
				releaseSources := component.NewReleaseSourceRepo(kilnfile, logger)

				Expect(releaseSources).To(HaveLen(3))
				Expect(releaseSources[0].ID()).To(Equal("comp"))
				Expect(releaseSources[1].ID()).To(Equal("buil"))
				Expect(releaseSources[2].ID()).To(Equal("bosh"))
			})
		})
	})

	Describe("FindReleaseUploader", func() {
		var (
			repo     component.ReleaseSourceList
			kilnfile cargo.Kilnfile
		)

		JustBeforeEach(func() {
			repo = component.NewReleaseSourceRepo(kilnfile, logger)
		})

		BeforeEach(func() {
			kilnfile = cargo.Kilnfile{
				ReleaseSources: cargo.ReleaseSourceList{
					cargo.S3ReleaseSource{Bucket: "bucket-1", Region: "us-west-1", AccessKeyId: "ak1", SecretAccessKey: "shhhh!",
						PathTemplate: `2.8/{{trimSuffix .Name "-release"}}/{{.Name}}-{{.Version}}-{{.StemcellOS}}-{{.StemcellVersion}}.tgz`},
					cargo.S3ReleaseSource{Bucket: "bucket-2", Region: "us-west-2", AccessKeyId: "aki", SecretAccessKey: "shhhh!",
						PathTemplate: `2.8/{{trimSuffix .Name "-release"}}/{{.Name}}-{{.Version}}.tgz`},
					cargo.BOSHIOReleaseSource{},
					cargo.S3ReleaseSource{Bucket: "bucket-3", Region: "us-west-2", AccessKeyId: "aki", SecretAccessKey: "shhhh!",
						PathTemplate: `{{.Name}}-{{.Version}}.tgz`},
				},
			}
		})

		Context("when the named source exists and accepts uploads", func() {
			It("returns a valid release uploader", func() {
				uploader, err := repo.FindReleaseUploader("bucket-2")
				Expect(err).NotTo(HaveOccurred())

				var s3ReleaseSource component.S3ReleaseSource
				Expect(uploader).To(BeAssignableToTypeOf(s3ReleaseSource))
			})
		})

		Context("when no sources accept uploads", func() {
			BeforeEach(func() {
				kilnfile = cargo.Kilnfile{
					ReleaseSources: cargo.ReleaseSourceList{cargo.BOSHIOReleaseSource{}},
				}
			})

			It("errors", func() {
				_, err := repo.FindReleaseUploader("bosh.io")
				Expect(err).To(MatchError(ContainSubstring("no upload-capable release sources were found")))
			})
		})

		Context("when the named source doesn't accept uploads", func() {
			It("errors with a list of valid sources", func() {
				_, err := repo.FindReleaseUploader("bosh.io")
				Expect(err).To(MatchError(ContainSubstring("could not find a valid matching release source")))
				Expect(err).To(MatchError(ContainSubstring("bucket-1")))
				Expect(err).To(MatchError(ContainSubstring("bucket-2")))
				Expect(err).To(MatchError(ContainSubstring("bucket-3")))
			})
		})

		Context("when the named source doesn't exist", func() {
			It("errors with a list of valid sources", func() {
				_, err := repo.FindReleaseUploader("bucket-42")
				Expect(err).To(MatchError(ContainSubstring("could not find a valid matching release source")))
				Expect(err).To(MatchError(ContainSubstring("bucket-1")))
				Expect(err).To(MatchError(ContainSubstring("bucket-2")))
				Expect(err).To(MatchError(ContainSubstring("bucket-3")))
			})
		})
	})

	Describe("RemotePather", func() {
		var (
			list     component.ReleaseSourceList
			kilnfile cargo.Kilnfile
		)

		JustBeforeEach(func() {
			list = component.NewReleaseSourceRepo(kilnfile, logger)
		})

		BeforeEach(func() {
			kilnfile = cargo.Kilnfile{
				ReleaseSources: cargo.ReleaseSourceList{
					cargo.S3ReleaseSource{Bucket: "bucket-1", Region: "us-west-1", AccessKeyId: "ak1", SecretAccessKey: "shhhh!",
						PathTemplate: `2.8/{{trimSuffix .Name "-release"}}/{{.Name}}-{{.Version}}-{{.StemcellOS}}-{{.StemcellVersion}}.tgz`},
					cargo.S3ReleaseSource{Bucket: "bucket-2", Region: "us-west-2", AccessKeyId: "aki", SecretAccessKey: "shhhh!",
						PathTemplate: `2.8/{{trimSuffix .Name "-release"}}/{{.Name}}-{{.Version}}.tgz`},
					cargo.BOSHIOReleaseSource{},
					cargo.S3ReleaseSource{Bucket: "bucket-3", Region: "us-west-2", AccessKeyId: "aki", SecretAccessKey: "shhhh!",
						PathTemplate: `{{.Name}}-{{.Version}}.tgz`},
				},
			}
		})

		Context("when the named source exists and implements RemotePath", func() {
			It("returns a valid release uploader", func() {
				uploader, err := list.FindRemotePather("bucket-2")
				Expect(err).NotTo(HaveOccurred())

				var s3ReleaseSource component.S3ReleaseSource
				Expect(uploader).To(BeAssignableToTypeOf(s3ReleaseSource))
			})
		})

		Context("when no sources implement RemotePath", func() {
			BeforeEach(func() {
				kilnfile = cargo.Kilnfile{
					ReleaseSources: cargo.ReleaseSourceList{cargo.BOSHIOReleaseSource{}},
				}
			})

			It("errors", func() {
				_, err := list.FindRemotePather("bosh.io")
				Expect(err).To(MatchError(ContainSubstring("no path-generating release sources were found")))
			})
		})

		Context("when the named source doesn't implement RemotePath", func() {
			It("errors with a list of valid sources", func() {
				_, err := list.FindRemotePather("bosh.io")
				Expect(err).To(MatchError(ContainSubstring("could not find a valid matching release source")))
				Expect(err).To(MatchError(ContainSubstring("bucket-1")))
				Expect(err).To(MatchError(ContainSubstring("bucket-2")))
				Expect(err).To(MatchError(ContainSubstring("bucket-3")))
			})
		})

		Context("when the named source doesn't exist", func() {
			It("errors with a list of valid sources", func() {
				_, err := list.FindRemotePather("bucket-42")
				Expect(err).To(MatchError(ContainSubstring("could not find a valid matching release source")))
				Expect(err).To(MatchError(ContainSubstring("bucket-1")))
				Expect(err).To(MatchError(ContainSubstring("bucket-2")))
				Expect(err).To(MatchError(ContainSubstring("bucket-3")))
			})
		})
	})
})
