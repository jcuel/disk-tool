"""Load outreach config and render templates with release variables."""

from __future__ import annotations

import re
from pathlib import Path
from typing import Any

import yaml

ROOT = Path(__file__).resolve().parents[2]
DEFAULT_CONFIG = ROOT / "config" / "outreach.yaml"

_PLACEHOLDER = re.compile(r"\{\{(\w+)\}\}")


def load_config(path: Path | None = None) -> dict[str, Any]:
    cfg_path = path or DEFAULT_CONFIG
    with cfg_path.open(encoding="utf-8") as f:
        return yaml.safe_load(f)


def build_context(version: str, config: dict[str, Any] | None = None) -> dict[str, str]:
    cfg = config or load_config()
    urls = cfg.get("urls", {})
    ver = version.lstrip("v")
    return {
        "version": ver,
        "demo_url": urls.get("demo", ""),
        "landing_url": urls.get("landing", ""),
        "release_url": f"{urls.get('release_base', '')}{ver}",
        "repo_url": urls.get("repo", ""),
    }


def render_string(text: str, context: dict[str, str]) -> str:
    def repl(match: re.Match[str]) -> str:
        key = match.group(1)
        return context.get(key, match.group(0))

    return _PLACEHOLDER.sub(repl, text)


def load_template(rel_path: str, context: dict[str, str]) -> str:
    path = ROOT / rel_path
    raw = path.read_text(encoding="utf-8")
    return render_string(raw, context).strip()


def append_license(body: str, config: dict[str, Any]) -> str:
    footer = (config.get("license_footer") or "").strip()
    if not footer:
        return body
    if footer in body:
        return body
    return f"{body.rstrip()}\n\n{footer}"


def render_hn(config: dict[str, Any], context: dict[str, str]) -> dict[str, str]:
    hn = config.get("hackernews", {})
    title = render_string(hn.get("title", ""), context)
    url = render_string(hn.get("url", context["demo_url"]), context)
    text = ""
    if tpl := hn.get("text_template"):
        text = append_license(load_template(tpl, context), config)
    return {"title": title, "url": url, "text": text}


def render_reddit_posts(config: dict[str, Any], context: dict[str, str]) -> list[dict[str, str]]:
    posts: list[dict[str, str]] = []
    for entry in config.get("reddit", []):
        kind = entry.get("kind", "link")
        title = render_string(entry.get("title", ""), context)
        subreddit = entry["subreddit"]
        body = ""
        if tpl := entry.get("body_template"):
            body = append_license(load_template(tpl, context), config)
        url = render_string(entry.get("url", context["demo_url"]), context) if kind == "link" else ""
        posts.append(
            {
                "subreddit": subreddit,
                "kind": kind,
                "title": title,
                "text": body if kind == "self" else "",
                "url": url if kind == "link" else "",
            }
        )
    return posts
