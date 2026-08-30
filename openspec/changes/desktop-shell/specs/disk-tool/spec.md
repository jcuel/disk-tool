# Delta Spec: Desktop shell

**Change:** `desktop-shell`
**Domain:** `disk-tool`

## ADDED Requirements

### Requirement: Desktop application distribution

The system SHALL provide a native desktop application that launches the disk-tool UI in an OS webview without requiring the user to run `disk-tool serve` manually or open a separate browser tab. The desktop app SHALL bundle the Go server as a sidecar process, navigate the webview to the sidecar localhost URL, and terminate the sidecar when the application exits.

#### Scenario: Desktop launch

- **WHEN** the user starts the desktop application
- **THEN** a native window opens showing the disk-tool UI
- **AND** the Go sidecar serves the embedded UI and API on localhost

#### Scenario: Desktop quit

- **WHEN** the user closes the desktop application
- **THEN** the Go sidecar process is stopped
- **AND** no orphan `disk-tool serve` process remains

#### Scenario: Sidecar dynamic port

- **WHEN** the sidecar starts with `--port 0`
- **THEN** the server binds an available localhost port
- **AND** emits a readiness line on stdout for the shell to parse

### Requirement: CLI distribution preserved

The system SHALL continue to ship standalone CLI binaries (`disk-tool serve`, `disk-tool scan`, `disk-tool version`) independent of the desktop installer.

#### Scenario: CLI serve unchanged

- **WHEN** a user runs `disk-tool serve` from the CLI release binary
- **THEN** the localhost web UI and API behave as before desktop packaging
