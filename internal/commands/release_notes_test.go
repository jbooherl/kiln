package commands_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	Ω "github.com/onsi/gomega"
	"github.com/pivotal-cf/jhanda"

	"github.com/pivotal-cf/kiln/internal/commands"
	"github.com/pivotal-cf/kiln/internal/commands/fakes"
	"github.com/pivotal-cf/kiln/internal/component"
	"github.com/pivotal-cf/kiln/pkg/cargo"
)

var _ jhanda.Command = commands.ReleaseNotes{}

func TestReleaseNotes_Usage(t *testing.T) {
	please := Ω.NewWithT(t)

	rn := commands.ReleaseNotes{}

	please.Expect(rn.Usage().Description).NotTo(Ω.BeEmpty())
	please.Expect(rn.Usage().ShortDescription).NotTo(Ω.BeEmpty())
	please.Expect(rn.Usage().Flags).NotTo(Ω.BeNil())
}

func TestReleaseNotes_Execute(t *testing.T) {
	t.Run("bump one release and use the default template", func(t *testing.T) {
		please := Ω.NewWithT(t)

		repo, _ := git.Init(memory.NewStorage(), memfs.New())

		revisionResolver := &fakes.RevisionResolver{}
		var (
			initialHash, finalHash plumbing.Hash
		)
		fill(initialHash[:], '1')
		fill(finalHash[:], '9')
		revisionResolver.ResolveRevisionReturnsOnCall(0, &initialHash, nil)
		revisionResolver.ResolveRevisionReturnsOnCall(1, &finalHash, nil)

		historicKilnfileLockFunc := &fakes.HistoricKilnfileLockFunc{}
		historicKilnfileLockFunc.ReturnsOnCall(0, cargo.KilnfileLock{
			Releases: []cargo.ComponentLock{
				{ComponentSpec: component.Spec{Name: "banana", Version: "0.1.0"}},
				{ComponentSpec: component.Spec{Name: "lemon", Version: "1.1.0"}},
			},
		}, nil)
		historicKilnfileLockFunc.ReturnsOnCall(1, cargo.KilnfileLock{
			Releases: []cargo.ComponentLock{
				{ComponentSpec: component.Spec{Name: "banana", Version: "0.1.0"}},
				{ComponentSpec: component.Spec{Name: "lemon", Version: "1.2.0"}},
			},
		}, nil)

		readFileCount := 0
		readFileFunc := func(fp string) ([]byte, error) {
			readFileCount++
			return nil, nil
		}

		var output bytes.Buffer
		rn := commands.ReleaseNotes{
			Repository:           repo,
			RevisionResolver:     revisionResolver,
			KilnfileLockAtCommit: historicKilnfileLockFunc.Spy,
			Writer:               &output,
			ReadFile:             readFileFunc,
		}

		err := rn.Execute([]string{
			"--version=0.1.0",
			"tile/1.1.0",
			"tile/1.2.0",
		})
		please.Expect(err).NotTo(Ω.HaveOccurred())

		please.Expect(revisionResolver.ResolveRevisionCallCount()).To(Ω.Equal(2))
		please.Expect(revisionResolver.ResolveRevisionArgsForCall(0)).To(Ω.Equal(plumbing.Revision("tile/1.1.0")))
		please.Expect(revisionResolver.ResolveRevisionArgsForCall(1)).To(Ω.Equal(plumbing.Revision("tile/1.2.0")))

		please.Expect(readFileCount).To(Ω.Equal(0))
		expected, err := ioutil.ReadFile("testdata/release_notes_output.md")
		please.Expect(err).NotTo(Ω.HaveOccurred())
		please.Expect(output.String()).To(Ω.Equal(string(expected)))
	})
}

func fill(buf []byte, value byte) {
	for i := range buf {
		buf[i] = value
	}
}