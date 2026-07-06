## ADDED Requirements

### Requirement: OS safety zones

The system SHALL classify paths into safety zones and block deletion of protected zones.

#### Scenario: Protected path delete rejected

- **WHEN** the user attempts to delete a path in critical_os or diagnostic zone
- **THEN** the API returns an error and no file is removed
