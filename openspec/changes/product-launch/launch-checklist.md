# Product launch checklist

Use after **v1.3.0** is on [Releases](https://github.com/jcuel/disk-tool/releases) with binaries and [Pages](https://jcuel.github.io/disk-tool/) is live.

## Repo polish (one-time)

- [ ] Homepage: `https://jcuel.github.io/disk-tool/`
- [ ] Description: *Local disk usage analyzer — overview-first scan, cleanup insights, support exports. Cross-platform Go + web UI.*
- [ ] Topics: `disk-usage`, `disk-analyzer`, `storage`, `cleanup`, `go`, `devtools`, `treesize-alternative`
- [ ] Release notes for v1.3.0 published

## Messaging pillars

- **Localhost-only** — scan data never leaves the machine
- **Overview-first drill-down** — TreeSize-style lazy scan
- **Safety zones** — review / maintenance / caution before delete
- **Try before install** — [live demo](https://jcuel.github.io/disk-tool/demo/)

## License (required in every post)

PolyForm Noncommercial 1.0.0 — free for personal and noncommercial use. Commercial use requires permission. See [LICENSE](https://github.com/jcuel/disk-tool/blob/dev/LICENSE).

## Show HN (draft)

**Title:** Show HN: disk-tool – local disk analyzer with lazy scan and cleanup insights

**Post:**

I built disk-tool because I wanted TreeSize-like overview drill-down without uploading anything to the cloud. It runs a Go backend on 127.0.0.1 with a web UI: pick a folder, see top consumers, drill into branches on demand, and export a support ticket for IT.

Safety zones flag what's safe to review vs OS-critical paths. There's a browser demo with sample data (no install): https://jcuel.github.io/disk-tool/demo/

Windows/Linux/macOS binaries: https://github.com/jcuel/disk-tool/releases

License: PolyForm Noncommercial — fine for personal use; commercial needs permission.

Feedback welcome on scan UX and cleanup heuristics.

## Reddit (check each sub's self-promo rules)

| Subreddit | Angle |
|-----------|--------|
| r/DataHoarder | Finding what's eating disk space, dev caches |
| r/sysadmin | Support ticket export, local-only auditing |
| r/selfhosted | Local web UI, no telemetry |
| r/golang | Go + embedded Vite SPA, lazy scanner design |

**Sample blurb:** Local disk-tool (Go): overview-first scan, drill-down, cleanup insights, GitHub Pages demo. Noncommercial license.

## Awesome lists

Many awesome lists require OSI-approved licenses. PolyForm NC may block inclusion — ask maintainers before opening a PR.

## GitHub Discussions (optional)

Enable Discussions → Announcements category → post release summary with demo + download links.
