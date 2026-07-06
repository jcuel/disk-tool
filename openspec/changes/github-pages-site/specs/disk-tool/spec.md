# Delta spec: GitHub Pages site

## ADDED Requirements

### Requirement: GitHub Pages product site

The project SHALL publish a static GitHub Pages site at `/disk-tool/` with a product landing page and links to install or try the app.

#### Scenario: Landing page loads

- **WHEN** a visitor opens `https://<owner>.github.io/disk-tool/`
- **THEN** they see product overview, features, safety zones, and install instructions

### Requirement: Interactive browser demo

The project SHALL provide a static demo at `/disk-tool/demo/` that reuses the disk-tool web UI with sample scan data and no Go backend.

#### Scenario: Demo auto-scan

- **WHEN** a visitor opens the demo URL
- **THEN** a sample overview scan completes and shows folder tree, charts, and insights

#### Scenario: Demo blocks destructive actions

- **WHEN** the user attempts delete, open-in-explorer, or cleanup in demo mode
- **THEN** the UI explains that local install is required

### Requirement: Pages deployment pipeline

The project SHALL build and deploy the site on push to `master` via GitHub Actions, and SHALL verify site + demo builds on pull requests.
