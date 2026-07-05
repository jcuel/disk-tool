# Proposal: Add project README

**Change:** `add-project-readme`
**Status:** archived
**Domain:** `disk-tool`

## Summary

Add a root README.md documenting the project purpose and SPECBOOT development workflow.

## Motivation

New contributors need a single entry point before reading OpenSpec artifacts.

## Scope

### In scope

- Root README.md with purpose and workflow commands

### Out of scope

- Feature implementation, CI, containers

## Risks

| Risk | Mitigation |
|------|------------|
| README drifts from workflow | Reference `.cursor/commands/` as source |
