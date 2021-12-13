package cargo

import (
	"fmt"

	"github.com/Masterminds/semver"
)

func Validate(spec Kilnfile, lock KilnfileLock) []error {
	var result []error

	for index, componentSpec := range spec.Releases {
		if componentSpec.Name == "" {
			result = append(result, fmt.Errorf("spec at index %d missing name", index))
			continue
		}

		componentLock, err := lock.FindReleaseWithName(componentSpec.Name)
		if err != nil {
			result = append(result,
				fmt.Errorf("component spec for release %q not found in lock", componentSpec.Name))
			continue
		}

		if err := checkComponentVersionsAndConstraint(componentSpec, componentLock, index); err != nil {
			result = append(result, err)
		}
	}

	if len(result) > 0 {
		return result
	}

	return nil
}

func checkComponentVersionsAndConstraint(spec ComponentSpec, lock ComponentLock, index int) error {
	v, err := semver.NewVersion(lock.Version)
	if err != nil {
		return fmt.Errorf("spec %s (index %d in Kilnfile.lock) has invalid lock version %q: %w",
			spec.Name, index, lock.Version, err)
	}

	if spec.Version != "" {
		c, err := semver.NewConstraint(spec.Version)
		if err != nil {
			return fmt.Errorf("spec %s (index %d in Kilnfile) has invalid version constraint: %w",
				spec.Name, index, err)
		}

		matches, errs := c.Validate(v)
		if !matches {
			return fmt.Errorf("spec %s version in lock %q does not match constraint %q: %v",
				spec.Name, lock.Version, spec.Version, errs)
		}
	}

	return nil
}

func checkStemcell(spec Kilnfile, lock KilnfileLock) []error {
	v, err := semver.NewVersion(lock.Stemcell.Version)
	if err != nil {
		return []error{fmt.Errorf("invalid lock version %q in Kilnfile.lock: %w",
			lock.Stemcell.Version, err)}
	}

	if spec.Stemcell.Version != "" {
		c, err := semver.NewConstraint(spec.Stemcell.Version)
		if err != nil {
			return []error{fmt.Errorf("invalid version constraint %q in Kilnfile: %w",
				lock.Stemcell.Version, err)}
		}

		matches, errs := c.Validate(v)
		if !matches {
			return []error{fmt.Errorf("stemcell version %s in Kilnfile.lock does not match constraint %q: %v",
				lock.Stemcell.Version, spec.Stemcell.Version, errs)}
		}
	}

	var result []error
	for index, componentLock := range lock.Releases {
		if componentLock.StemcellOS == "" {
			continue
		}
		if componentLock.StemcellOS != lock.Stemcell.OS {
			result = append(result, fmt.Errorf("spec %s (index %d in Kilnfile) has stemcell os that does not match the stemcell lock os",
				componentLock.Name, index))
		}
		if componentLock.StemcellVersion != lock.Stemcell.Version {
			result = append(result, fmt.Errorf("spec %s (index %d in Kilnfile) has stemcell version that does not match the stemcell lock (expected %s but got %s)",
				componentLock.Name, index, lock.Stemcell.Version, componentLock.StemcellVersion))
		}
	}

	return result
}
