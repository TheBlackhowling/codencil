# Codencil — Agent Guide

Single source of truth for **how AI agents implement Codencil**. Read with `BUILD_ORDER.md` before coding.

---

## Non-negotiables

- **Go backend** — not Python
- **Postgres** — source of truth for documents (not git mount, not localStorage)
- **TypRow** — `github.com/TheBlackHowling/typrow` for SQL access (not sqlc, not GORM)
- **go-migrate** — SQL migrations in `db/migrations/`
- **Markdown render in frontend** — API stores raw markdown strings only
- **One task per session** — pick exactly one task from `BUILD_ORDER.md`; finish or clearly partial with progress update
- **No scope creep** — if the task is "publish v1 API", do not add comments, auth, or Docker polish unless the task includes it
- **Docker-only host** — no Go, Node, npm, or migrate on the developer machine. **All** build, test, migrate, and run commands go through **Docker Compose** (see P0.2)

---

## Defaults (resolve open questions — do not re-debate)

| Question | Default |
|---|---|
| HTTP router | **chi** (`github.com/go-chi/chi/v5` ≥ **v5.2.4**) |
| chi middleware | **Do not use** `middleware.RedirectSlashes`; use `RequestID`, `Recoverer`, `Logger` |
| License | **MIT** |
| Frontend layout | **`apps/web`** (Next.js App Router) |
| Real-time / WebSockets | **Defer** — polling or refresh OK for MVP |
| Multi-tenant | **`org_id` on tables** from day one; single-org dev mode OK |
| Auth before Phase 4 | **Dev stub** — `X-Dev-User-Id` header or env default user |
| Approval workflow | **Out of scope** — see `DECISIONS.md` |
| Local toolchain | **None** — Docker Compose only (`P0.2`) |

---

## Repo layout (target)

```
codencil/
  agents/              ← planning & agent docs (this folder)
    CONTEXT.md
    AGENTS.md
    BUILD_ORDER.md
    DECISIONS.md
    progress.md
  .cursor/rules/     ← points here
  README.md
  docker-compose.yml
  Makefile             ← docker compose shortcuts (optional)
  db/migrations/
  apps/
    api/
      Dockerfile       ← dev + prod targets
      cmd/codencil/main.go
      internal/
        publish/     # pure Go — no DB imports in core logic
        models/
        store/
        http/
        auth/        # Phase 4+
    web/
      Dockerfile       ← dev + prod targets
```

---

## Implementation rules

### Backend

- **chi v5** on `net/http` — handlers are `http.HandlerFunc`; JSON via `encoding/json`
- Pin **`github.com/go-chi/chi/v5` ≥ v5.2.4**; run `govulncheck` in CI when available
- **Do not** use `middleware.RedirectSlashes` (open-redirect risk; not needed for JSON API)
- Configure **`middleware.RealIP`** only behind a trusted reverse proxy (Phase 5+)
- Absolute imports from module root (`github.com/TheBlackHowling/codencil/apps/api/...`)
- Business logic in `internal/` — handlers stay thin
- **`internal/publish`** must not import `database/sql` or TypRow — test with table-driven cases
- All SQL in TypRow model methods or small `store/` queries — no ORM
- Return appropriate HTTP status codes; JSON error shape consistent within project

### Frontend

- Use `react-markdown` + GFM plugins for preview (Phase 1+)
- API client: typed fetch to Go backend; no secret env in client bundle
- Margin comment UI: build on **published version** snapshot only (never anchor to live draft)

### Migrations

- Paired `.up.sql` / `.down.sql`
- Numbered: `000001_init.up.sql`
- One logical change per migration when possible

---

## Verification checklist (by layer)

**Prerequisite:** `docker compose build` succeeds (from **P0.2**).

**After API changes:**
```bash
docker compose run --rm api go test ./...
docker compose run --rm api go vet ./...
```

**After migration changes:**
```bash
docker compose up -d postgres
docker compose run --rm migrate -path /migrations -database "$DATABASE_URL" up
docker compose run --rm migrate -path /migrations -database "$DATABASE_URL" down 1
docker compose run --rm migrate -path /migrations -database "$DATABASE_URL" up
```

**After web changes:**
```bash
docker compose run --rm web npm run build
```

**End of Phase 1 smoke (via compose):**
1. `docker compose up -d`
2. Create document via API (`curl` to exposed api port)
3. Save draft markdown
4. Publish v1
5. Open web preview URL — rendered markdown matches snapshot

---

## Anti-patterns (agents fail here)

| Don't | Do instead |
|---|---|
| Build OIDC in Phase 1 | Dev auth stub; real auth Phase 4 |
| Anchor comments to draft | Anchors on published version snapshot |
| Port Markdown Viewer wholesale | react-markdown + fresh UI |
| Implement approval gates | Publish + resolve threads only |
| One giant PR / session | One `BUILD_ORDER` task |
| Open PR to `main` for P0.2 while P0.1 still open | Stack: PR base = previous feature branch |
| Merge your own PR | Leave open for maintainer |
| Open PR as draft | Open ready for review |
| Modify TypRow repo | Use published module; file issue if bug |
| Run `go test` / `npm` on host | `docker compose run --rm api go test ./...` |
| Assume local Go/Node installed | Docker-only; see P0.2 |
| Skip `agents/progress.md` update | Always hand off next agent |

---

## Pull requests (stacked)

Full workflow: [`agents/STACK.md`](STACK.md)

- **Start every session:** `gh pr list --state open` + [`progress.md`](progress.md) stack table
- **No open PRs** → branch from `main`, PR base `main`
- **Open PRs** → branch from **tip of stack** (last completed task branch), PR base = that branch
- **One `BUILD_ORDER` task per branch/PR** — open ready for review; **do not merge**
- Update **`agents/progress.md`** Open stack table + session log
- Use [PR template](../.github/pull_request_template.md); CI runs per PR

---

## Session prompt template

Copy into a new agent chat:

```
Working on Codencil at C:\source\codencil.

1. Read agents/CONTEXT.md, agents/STACK.md, agents/progress.md
2. Check open PRs: gh pr list --state open
   - None → branch from main
   - Open stack → branch from tip (last row in progress.md stack table)
3. Execute BUILD_ORDER task: [TASK ID]
4. Open PR ready for review (base = main or tip branch; no `--draft`)
5. Update agents/progress.md open stack table
6. Do not merge
```

---

## When to stop and ask the user

- Changing a decision in `DECISIONS.md`
- Adding dependencies (npm, Go modules beyond chi/typrow/pgx) — chi middleware from `go-chi/chi/v5/middleware` is OK
- Deviating from `BUILD_ORDER` task scope
- Schema changes that drop/rename columns after Phase 2 ships
- Anything involving work `typedb` repos or OSS approval
