# Codencil

Self-hosted markdown review platform — authors publish versioned documents; reviewers comment in the margin.

**Status:** Phase 2 complete (review path). Phase 3 (publish v2 + anchor migration) in progress. See [`agents/progress.md`](agents/progress.md).

**License:** [MIT](LICENSE)

**Contributing:** see [CONTRIBUTING.md](CONTRIBUTING.md) — one BUILD_ORDER task per PR, Docker-only verify.

**Dev requirement:** Docker + Docker Compose only (no local Go or Node needed).

## Layout

```
apps/api/          Go module (github.com/TheBlackHowling/codencil/apps/api)
apps/web/          Next.js frontend (Phase 0.5+)
db/migrations/     go-migrate SQL files
agents/            Planning & agent docs
```

```
Write in markdown. Review in the margin. Publish when it's ready.
```
