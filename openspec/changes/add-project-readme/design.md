# Design: Add project README

**Change:** `add-project-readme`

## Approach

Single Markdown file at repo root. No build tooling required.

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| README only | Minimal, standard | None for this scope | Selected |
| docs/ site | Richer docs | Overkill for greenfield | Rejected |

## Components

- `README.md` (new)

## Data & API Touchpoints

None.
