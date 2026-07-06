#!/usr/bin/env python3
"""Orchestrate outreach posts after a release tag."""

from __future__ import annotations

import argparse
import json
import os
import re
import subprocess
import sys
from pathlib import Path
from typing import Any

# Allow running as script without package install
sys.path.insert(0, str(Path(__file__).resolve().parent))

from hn_submit import submit_link  # noqa: E402
from reddit_post import submit_posts  # noqa: E402
from render import build_context, load_config, render_hn, render_reddit_posts  # noqa: E402

REPO = os.environ.get("GITHUB_REPOSITORY", "jcuel/disk-tool")


def _tag(version: str) -> str:
    return version if version.startswith("v") else f"v{version}"


def _marker_prefix(config: dict[str, Any]) -> str:
    return config.get("idempotency_marker_prefix", "outreach-posted:")


def _parse_posted_channels(body: str, prefix: str) -> set[str]:
    m = re.search(rf"{re.escape(prefix)}\s*([^\n]+)", body)
    if not m:
        return set()
    return {c.strip() for c in m.group(1).split(",") if c.strip()}


def _gh_json(args: list[str]) -> dict[str, Any] | list[Any] | None:
    cmd = ["gh", *args, "--repo", REPO]
    try:
        out = subprocess.check_output(cmd, stderr=subprocess.PIPE, text=True)
        return json.loads(out) if out.strip() else None
    except (subprocess.CalledProcessError, FileNotFoundError, json.JSONDecodeError):
        return None


def get_release_body(version: str) -> str:
    tag = _tag(version)
    data = _gh_json(["release", "view", tag, "--json", "body"])
    if isinstance(data, dict):
        return data.get("body") or ""
    return ""


def update_release_marker(version: str, channels: set[str], summary: str, config: dict[str, Any]) -> None:
    tag = _tag(version)
    prefix = _marker_prefix(config)
    body = get_release_body(version)
    marker_line = f"{prefix} {','.join(sorted(channels))}"
    if prefix in body:
        body = re.sub(rf"{re.escape(prefix)}[^\n]*", marker_line, body)
    else:
        body = f"{body.rstrip()}\n\n<!-- {marker_line} -->\n\n## Outreach\n\n{summary}\n"
    subprocess.run(
        ["gh", "release", "edit", tag, "--repo", REPO, "--notes", body],
        check=True,
    )


def run_outreach(
    version: str,
    *,
    channel: str = "all",
    dry_run: bool = False,
    force: bool = False,
) -> dict[str, Any]:
    config = load_config()
    context = build_context(version, config)
    prefix = _marker_prefix(config)
    posted = _parse_posted_channels(get_release_body(version), prefix) if not force else set()

    results: dict[str, Any] = {"version": _tag(version), "dry_run": dry_run, "channels": {}}

    if channel in ("all", "hn"):
        if "hn" in posted and not dry_run and not force:
            results["channels"]["hn"] = {"status": "skipped", "reason": "already posted"}
        else:
            hn = render_hn(config, context)
            results["channels"]["hn"] = submit_link(
                hn["title"], hn["url"], dry_run=dry_run
            )

    if channel in ("all", "reddit"):
        if "reddit" in posted and not dry_run and not force:
            results["channels"]["reddit"] = {"status": "skipped", "reason": "already posted"}
        else:
            posts = render_reddit_posts(config, context)
            stagger = int(config.get("stagger_seconds", 600))
            results["channels"]["reddit"] = submit_posts(
                posts, stagger_seconds=stagger, dry_run=dry_run
            )

    if not dry_run and not force:
        new_channels = set(posted)
        for key, val in results["channels"].items():
            if isinstance(val, dict) and val.get("status") in ("posted", "dry_run"):
                if val.get("status") == "posted":
                    new_channels.add(key)
            elif isinstance(val, list) and any(r.get("status") == "posted" for r in val):
                new_channels.add(key)
        if new_channels - posted:
            summary = json.dumps(results["channels"], indent=2)
            try:
                update_release_marker(version, new_channels, summary, config)
            except subprocess.CalledProcessError as e:
                results["marker_error"] = str(e)

    return results


def main() -> int:
    p = argparse.ArgumentParser(description="Run release outreach (HN + Reddit)")
    p.add_argument("--version", required=True, help="Semver e.g. 1.3.0 or v1.3.0")
    p.add_argument("--channel", choices=("all", "hn", "reddit"), default="all")
    p.add_argument("--dry-run", action="store_true", help="Render only, no network posts")
    p.add_argument("--force", action="store_true", help="Ignore idempotency marker")
    args = p.parse_args()

    try:
        out = run_outreach(
            args.version,
            channel=args.channel,
            dry_run=args.dry_run,
            force=args.force,
        )
        print(json.dumps(out, indent=2))
        return 0
    except Exception as e:
        print(f"error: {e}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
