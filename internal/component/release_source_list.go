package component

import (
	"errors"
	"fmt"
	"log"

	"github.com/pivotal-cf/kiln/pkg/cargo"
)

type ReleaseSourceList []ReleaseSource

func NewReleaseSourceRepo(kilnfile cargo.Kilnfile, logger *log.Logger) ReleaseSourceList {
	var list ReleaseSourceList

	for _, releaseConfig := range kilnfile.ReleaseSources {
		list = append(list, ReleaseSourceFactory(releaseConfig, logger))
	}

	panicIfDuplicateIDs(list)

	return list
}

func (list ReleaseSourceList) Filter(allowOnlyPublishable bool) ReleaseSourceList {
	var sources ReleaseSourceList
	for _, source := range list {
		if allowOnlyPublishable && !source.Configuration().Publishable {
			continue
		}
		sources = append(sources, source)
	}
	return sources
}

func (list ReleaseSourceList) FindReleaseUploader(sourceID string) (ReleaseUploader, error) {
	var (
		uploader     ReleaseUploader
		availableIDs []string
	)
	for _, src := range list {
		u, ok := src.(ReleaseUploader)
		if !ok {
			continue
		}
		availableIDs = append(availableIDs, src.Configuration().ID)
		if src.Configuration().ID == sourceID {
			uploader = u
			break
		}
	}

	if len(availableIDs) == 0 {
		return nil, errors.New("no upload-capable release sources were found in the Kilnfile")
	}

	if uploader == nil {
		return nil, fmt.Errorf(
			"could not find a valid matching release source in the Kilnfile, available upload-compatible sources are: %q",
			availableIDs,
		)
	}

	return uploader, nil
}

func (list ReleaseSourceList) FindRemotePather(sourceID string) (RemotePather, error) {
	var (
		pather       RemotePather
		availableIDs []string
	)

	for _, src := range list {
		u, ok := src.(RemotePather)
		if !ok {
			continue
		}
		id := src.Configuration().ID
		availableIDs = append(availableIDs, id)
		if id == sourceID {
			pather = u
			break
		}
	}

	if len(availableIDs) == 0 {
		return nil, errors.New("no path-generating release sources were found in the Kilnfile")
	}

	if pather == nil {
		return nil, fmt.Errorf(
			"could not find a valid matching release source in the Kilnfile, available path-generating sources are: %q",
			availableIDs,
		)
	}

	return pather, nil
}

func panicIfDuplicateIDs(releaseSources []ReleaseSource) {
	indexOfID := make(map[string]int)
	for index, rs := range releaseSources {
		id := rs.Configuration().ID
		previousIndex, seen := indexOfID[id]
		if seen {
			panic(fmt.Sprintf(`release_sources must have unique IDs; items at index %d and %d both have ID %q`, previousIndex, index, id))
		}
		indexOfID[id] = index
	}
}

func NewMultiReleaseSource(sources ...ReleaseSource) ReleaseSourceList {
	return sources
}

func (list ReleaseSourceList) SetDownloadThreads(n int) {
	for i, rs := range list {
		s3rs, ok := rs.(S3ReleaseSource)
		if ok {
			s3rs.DownloadThreads = n
			list[i] = s3rs
		}
	}
}

func (list ReleaseSourceList) FindByID(id string) (ReleaseSource, error) {
	if id == "" {
		return nil, newReleaseSourceNotFoundError(id, list)
	}
	for _, src := range list {
		if src.Configuration().ID == id {
			return src, nil
		}
	}
	return nil, newReleaseSourceNotFoundError(id, list)
}

func (list ReleaseSourceList) GetReleaseCache() (ReleaseSource, error) {
	for _, src := range list {
		if src.Configuration().Publishable {
			return src, nil
		}
	}
	return nil, fmt.Errorf("publishable release source not found")
}

func newReleaseSourceNotFoundError(id string, list ReleaseSourceList) error {
	ids := make([]string, 0, len(list))
	for _, src := range list {
		ids = append(ids, src.Configuration().ID)
	}
	return fmt.Errorf("release source with name %q not found in Kilnfile release_sources (the named release sources are: %q)", id, ids)
}

func scopedError(sourceID string, err error) error {
	return fmt.Errorf("error from release source %q: %w", sourceID, err)
}
