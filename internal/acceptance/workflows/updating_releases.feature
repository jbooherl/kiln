Feature: Updating a BOSH Release
  Scenario: Find a version on GitHub
    Given I have a "hello-tile" repository checked out at v0.1.1
    And GitHub repository "crhntr/hello-release" has release with tag "v0.1.4"
    When I invoke kiln find-release-version for "hello-release"
    Then stdout contains substring: "0.1.5"

  Scenario: Find a version on bosh.io
    Given I set the version constraint to "1.1.18" for release "bpm"
    When I invoke kiln find-release-version for "bpm"
    Then stdout contains substring: "1.1.18"

  Scenario: Update a component to a new release
    Given I have a "hello-tile" repository checked out at v0.1.1
    And the Kilnfile.lock specifies version "v0.1.3" for release "hello-release"
    And GitHub repository "crhntr/hello-release" has release with tag "v0.1.4"
    When I invoke kiln update-release for releas "hello-release" with version "0.1.4"
    Then the Kilnfile.lock specifies version "0.1.4" for release "hello-release"
    And kiln validate succeeds