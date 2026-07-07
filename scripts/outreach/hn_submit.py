"""Submit a link post to Hacker News via session login (unofficial — no write API)."""

from __future__ import annotations

import os
import re
import sys

import requests

HN_BASE = "https://news.ycombinator.com"


def _extract_form_fields(html: str) -> dict[str, str]:
    fields: dict[str, str] = {}
    for name in ("fnid", "hmac", "goto"):
        m = re.search(rf'name="{name}"\s+value="([^"]*)"', html)
        if m:
            fields[name] = m.group(1)
    return fields


def _login(session: requests.Session, username: str, password: str) -> None:
    r = session.get(f"{HN_BASE}/login", timeout=30)
    r.raise_for_status()
    fields = _extract_form_fields(r.text)
    if "fnid" not in fields:
        raise RuntimeError("HN login page missing fnid token")
    payload = {
        "fnid": fields["fnid"],
        "fnop": "login-page",
        "goto": fields.get("goto", "news"),
        "acct": username,
        "pw": password,
    }
    if "hmac" in fields:
        payload["hmac"] = fields["hmac"]
    r = session.post(f"{HN_BASE}/login", data=payload, timeout=30, allow_redirects=True)
    r.raise_for_status()
    if "logout" not in r.text.lower() and "Bad login" in r.text:
        raise RuntimeError("HN login failed — check HN_USERNAME / HN_PASSWORD")


def submit_link(title: str, url: str, *, dry_run: bool = False) -> dict[str, str]:
    username = os.environ.get("HN_USERNAME", "")
    password = os.environ.get("HN_PASSWORD", "")
    if dry_run:
        return {"status": "dry_run", "title": title, "url": url}
    if not username or not password:
        raise RuntimeError("HN_USERNAME and HN_PASSWORD required for live submit")

    session = requests.Session()
    session.headers.update(
        {
            "User-Agent": "disk-tool-outreach/1.0 (Show HN automation)",
        }
    )
    _login(session, username, password)

    r = session.get(f"{HN_BASE}/submit", timeout=30)
    r.raise_for_status()
    fields = _extract_form_fields(r.text)
    if "fnid" not in fields:
        raise RuntimeError("HN submit page missing fnid token")

    payload = {
        "fnid": fields["fnid"],
        "fnop": "submit-page",
        "title": title,
        "url": url,
        "text": "",
    }
    if "hmac" in fields:
        payload["hmac"] = fields["hmac"]

    r = session.post(f"{HN_BASE}/submit", data=payload, timeout=30, allow_redirects=True)
    r.raise_for_status()

    item_id = ""
    m = re.search(r"item\?id=(\d+)", r.url)
    if m:
        item_id = m.group(1)
    elif m := re.search(r"item\?id=(\d+)", r.text):
        item_id = m.group(1)

    result = {"status": "posted", "title": title, "url": url}
    if item_id:
        result["item_url"] = f"{HN_BASE}/item?id={item_id}"
    return result


def main() -> int:
    import argparse

    p = argparse.ArgumentParser(description="Submit Show HN link post")
    p.add_argument("--title", required=True)
    p.add_argument("--url", required=True)
    p.add_argument("--dry-run", action="store_true")
    args = p.parse_args()
    try:
        out = submit_link(args.title, args.url, dry_run=args.dry_run)
        print(out)
        return 0
    except Exception as e:
        print(f"error: {e}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
