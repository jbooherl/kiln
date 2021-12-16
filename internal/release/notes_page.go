package release

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

const DefaultReleasesSentinel = "\n## <a id='releases'></a> Releases\n\n"

var releaseNoteExp = regexp.MustCompile(`(?m)(?P<notes>### <a id='(?P<version>\d+\.\d+\.\d+(-.+)?)'></a> (\d+\.\d+\.\d+(-.+)?)\w*(?P<header_suffix>.*)\n*((\*.*\n)|\n|(\</?.*)|(  .*)|(####+.*)|(\t.*)|(\w.*))*)`)

type NotesPage struct {
	Exp *regexp.Regexp

	Prefix, Suffix string
	Releases       []VersionNote
}

func ParseNotesPage(input string) (NotesPage, error) {
	return ParseNotesPageWithExpressionAndReleasesSentinel(input, releaseNoteExp.String(), DefaultReleasesSentinel)
}

func ParseNotesPageWithExpressionAndReleasesSentinel(input, releaseRegularExpression, releasesSentinel string) (NotesPage, error) {
	const (
		versionCaptureGroup = "version"
		notesCaptureGroup   = "notes"
	)

	if !strings.Contains(input, releasesSentinel) {
		return NotesPage{}, fmt.Errorf("releases sentinal not found in input: expected input to contain %q", releasesSentinel)
	}

	exp, err := regexp.Compile(releaseRegularExpression)
	if err != nil {
		return NotesPage{}, fmt.Errorf(`release regular expression parse failure: %w`, err)
	}
	if !stringsSliceContains(exp.SubexpNames(), versionCaptureGroup) {
		return NotesPage{}, fmt.Errorf(`release regular expression must contain named capture group %q`, versionCaptureGroup)
	}
	if !stringsSliceContains(exp.SubexpNames(), notesCaptureGroup) {
		return NotesPage{}, fmt.Errorf(`release regular expression must contain named capture group %q`, notesCaptureGroup)
	}

	page := NotesPage{
		Exp: exp,
	}

	matchIndices := page.Exp.FindAllStringSubmatchIndex(input, -1)

	switch len(matchIndices) {
	case 0:
		index := strings.Index(input, releasesSentinel)
		page.Prefix = input[:index+len(releasesSentinel)]
		page.Suffix = input[index+len(releasesSentinel):]
	default:
		versionSubExpIndex := page.Exp.SubexpIndex(versionCaptureGroup)
		notesSubExpIndex := page.Exp.SubexpIndex(notesCaptureGroup)
		matchStrings := page.Exp.FindAllStringSubmatch(input, -1)

		for _, match := range matchStrings {
			page.Releases = append(page.Releases, VersionNote{
				Version: match[versionSubExpIndex],
				Notes:   match[notesSubExpIndex],
			})
		}

		page.Prefix = input[:matchIndices[0][notesSubExpIndex+1]]
		page.Suffix = input[matchIndices[len(matchIndices)-1][notesSubExpIndex+2]:]
	}

	return page, nil
}

func (page NotesPage) validateRelease(tile VersionNote) error {
	_, err := tile.version()
	if err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}
	if !page.Exp.MatchString(tile.Notes) {
		return fmt.Errorf("notes do not match expression")
	}
	return nil
}

func (page *NotesPage) Add(versionNote VersionNote) error {
	err := page.validateRelease(versionNote)
	if err != nil {
		return err
	}

	if len(page.Releases) == 0 {
		page.Releases = []VersionNote{versionNote}
		return nil
	}

	nv, _ := versionNote.version()
	for i, t := range page.Releases {
		tv, err := t.version()
		if err != nil {
			continue
		}
		if nv.Equal(tv) {
			page.Releases[i] = versionNote
			return nil
		}
		if !nv.GreaterThan(tv) {
			continue
		}
		page.Releases = append(page.Releases[:i], append([]VersionNote{versionNote}, page.Releases[i:]...)...)
		return nil
	}

	page.Releases = append(page.Releases, versionNote)

	return nil
}

func (page *NotesPage) WriteTo(w io.Writer) (int64, error) {
	buf := new(bytes.Buffer)
	n, err := buf.WriteString(page.Prefix)
	if err != nil {
		return int64(n), err
	}
	for _, r := range page.Releases {
		n, err = buf.WriteString(r.Notes)
		if err != nil {
			return int64(n), err
		}
	}
	n, err = buf.WriteString(page.Suffix)
	if err != nil {
		return int64(n), err
	}
	return buf.WriteTo(w)
}

type VersionNote struct {
	Version string
	Notes   string
}

func (notes VersionNote) version() (*semver.Version, error) {
	return semver.NewVersion(notes.Version)
}

func stringsSliceContains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}
