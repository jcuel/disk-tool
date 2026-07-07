"""Post to Reddit via OAuth2 script-app password grant."""

from __future__ import annotations

import os
import sys
import time
from typing import Any

import requests

REDDIT_AUTH = "https://www.reddit.com/api/v1/access_token"
REDDIT_API = "https://oauth.reddit.com"


def _user_agent() -> str:
    return os.environ.get(
        "REDDIT_USER_AGENT",
        "disk-tool-outreach/1.0 (by /u/disk-tool-bot)",
    )


def _get_token() -> str:
    client_id = os.environ.get("REDDIT_CLIENT_ID", "")
    client_secret = os.environ.get("REDDIT_CLIENT_SECRET", "")
    username = os.environ.get("REDDIT_USERNAME", "")
    password = os.environ.get("REDDIT_PASSWORD", "")
    if not all([client_id, client_secret, username, password]):
        raise RuntimeError(
            "REDDIT_CLIENT_ID, REDDIT_CLIENT_SECRET, REDDIT_USERNAME, REDDIT_PASSWORD required"
        )

    r = requests.post(
        REDDIT_AUTH,
        auth=(client_id, client_secret),
        data={"grant_type": "password", "username": username, "password": password},
        headers={"User-Agent": _user_agent()},
        timeout=30,
    )
    r.raise_for_status()
    data = r.json()
    token = data.get("access_token")
    if not token:
        raise RuntimeError(f"Reddit auth failed: {data}")
    return token


def _parse_submit_response(data: dict[str, Any]) -> dict[str, str]:
    errors = data.get("json", {}).get("errors") or data.get("errors")
    if errors:
        raise RuntimeError(f"Reddit submit errors: {errors}")
    j = data.get("json", {}).get("data", {})
    url = j.get("url", "")
    name = j.get("name", "")
    return {"status": "posted", "url": url, "fullname": name}


def submit_post(
    subreddit: str,
    title: str,
    *,
    kind: str = "link",
    url: str = "",
    text: str = "",
    dry_run: bool = False,
) -> dict[str, str]:
    if dry_run:
        return {
            "status": "dry_run",
            "subreddit": subreddit,
            "kind": kind,
            "title": title,
        }

    token = _get_token()
    payload: dict[str, str] = {
        "sr": subreddit,
        "title": title[:300],
        "kind": kind,
        "sendreplies": "true",
    }
    if kind == "link":
        payload["url"] = url
    else:
        payload["text"] = text

    r = requests.post(
        f"{REDDIT_API}/api/submit",
        data=payload,
        headers={
            "Authorization": f"Bearer {token}",
            "User-Agent": _user_agent(),
        },
        timeout=30,
    )
    r.raise_for_status()
    return _parse_submit_response(r.json())


def submit_posts(
    posts: list[dict[str, str]],
    *,
    stagger_seconds: int = 0,
    dry_run: bool = False,
) -> list[dict[str, str]]:
    results: list[dict[str, str]] = []
    for i, post in enumerate(posts):
        if i > 0 and stagger_seconds > 0 and not dry_run:
            time.sleep(stagger_seconds)
        kind = post.get("kind", "link")
        out = submit_post(
            post["subreddit"],
            post["title"],
            kind=kind,
            url=post.get("url", ""),
            text=post.get("text", ""),
            dry_run=dry_run,
        )
        out["subreddit"] = post["subreddit"]
        results.append(out)
    return results


def main() -> int:
    import argparse
    import json

    p = argparse.ArgumentParser(description="Submit Reddit post")
    p.add_argument("--subreddit", required=True)
    p.add_argument("--title", required=True)
    p.add_argument("--kind", choices=("link", "self"), default="link")
    p.add_argument("--url", default="")
    p.add_argument("--text", default="")
    p.add_argument("--dry-run", action="store_true")
    args = p.parse_args()
    try:
        out = submit_post(
            args.subreddit,
            args.title,
            kind=args.kind,
            url=args.url,
            text=args.text,
            dry_run=args.dry_run,
        )
        print(json.dumps(out))
        return 0
    except Exception as e:
        print(f"error: {e}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
