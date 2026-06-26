# Codencil — Implementation Progress

> **Agents:** Update this file at the end of every session.  
> **Source of truth for "what's next"** — see `BUILD_ORDER.md` for task definitions.

**Last updated:** 2026-06-10  
**Current phase:** 0 (in progress)  
**Next task:** **P0.4** — API skeleton  
**Stack policy:** Stacked PRs (ready for review, not draft) — see [`STACK.md`](STACK.md)

---

## Open stack (branches & PRs)

| Task | Branch | PR base | PR | Status |
|---|---|---|---|---|
| P0.3 | `feature/p0.3-migrate` | `main` | *(opening)* | open |

*Agents: add a row when opening each PR. Remove or mark merged after maintainer merge.*

**Resume rule:** at session start, run `gh pr list --state open`. No open PRs → `main`. Otherwise → checkout **tip** branch.

---

## Phase 0 — Project shell

- [x] **P0.1** Repo scaffold
- [x] **P0.2** Docker Compose dev stack (Go + Node in containers)
- [x] **P0.3** go-migrate wiring
- [ ] **P0.4** API skeleton
- [ ] **P0.5** Web skeleton

## Phase 1 — Read path

- [ ] **P1.1** Migration: documents + versions
- [ ] **P1.2** TypRow models + store
- [ ] **P1.3** `internal/publish` diff scaffold
- [ ] **P1.4** HTTP: document CRUD + publish v1
- [ ] **P1.5** Web: markdown preview page
- [ ] **P1.6** Phase 1 smoke doc

## Phase 2 — Review path

- [ ] **P2.1** Migration: threads, comments, anchors
- [ ] **P2.2** TypRow models for review entities
- [ ] **P2.3** HTTP: comment API
- [ ] **P2.4** Web: text selection → anchor
- [ ] **P2.5** Web: thread panel UI

## Phase 3 — Publish v2 + anchor migration

- [ ] **P3.1** Anchor remap logic
- [ ] **P3.2** Wire publish v2+
- [ ] **P3.3** HTTP + UI version selector
- [ ] **P3.4** Orphaned anchor UX

## Phase 4 — Auth & roles

- [ ] **P4.1** Users table + dev auth middleware
- [ ] **P4.2** Document membership / roles
- [ ] **P4.3** OIDC (optional MVP+)
- [ ] **P4.4** Remove or gate dev auth in prod

## Phase 5 — Self-host polish

- [ ] **P5.1** Production compose + CI
- [ ] **P5.2** README self-host guide
- [ ] **P5.3** GitHub publish — repo live; CI on PRs; branch protection optional

---

## Session log

### 2026-06-10 — P0.3 go-migrate wiring

- Seed migration `000001_init` (no-op); removed migrate `tools` profile
- Added `migrate-up`, `migrate-down`, `migrate-reset` to Makefile and `scripts/dev.ps1`
- CI runs migrate up / down 1 / up when migration files exist
- **Next agent:** **P0.4** API skeleton (chi, `/health`)

### 2026-06-10 — P0.2 Docker Compose dev stack

- Added `docker-compose.yml` (postgres, migrate, api, web)
- Added `apps/api/Dockerfile`, `apps/web/Dockerfile`, `Makefile`, `scripts/dev.ps1`, `.env.example`
- Verified: `docker compose build`, `docker compose run --rm api go version`, postgres healthy
- **Next agent:** **P0.3** from `main` after P0.2 merges (or stack tip if PR open)

### 2026-06-10 — PR workflow updates

- No draft PRs — open ready for review
- Resume rule: check open PRs; none → `main`, else branch from stack tip
- Updated STACK.md, AGENTS.md, CONTRIBUTING.md, PR template

### 2026-06-10 — P0.1 scaffold (stacked)

- Branch `feature/p0.1-scaffold` from `feature/docs-stacked-pr-workflow`
- Added `apps/api/go.mod`, module layout (`internal/*`, `cmd/codencil`, `db/migrations`, `apps/web` placeholders)
- PR #2 stacked on docs PR #1
- **Next agent:** **P0.2** on branch from `feature/p0.1-scaffold`

### 2026-06-10 — Stacked PR workflow

- Added `agents/STACK.md`; CONTRIBUTING/AGENTS/PR template updated for stacked PRs
- Agents branch from previous task; maintainer merges bottom → top

### 2026-06-10 — Contributing

- Added `CONTRIBUTING.md`, PR template, Docker-based CI workflow, `CODEOWNERS`

### 2026-06-10 — GitHub

- Published planning docs to https://github.com/TheBlackhowling/codencil (MIT, public)
- **Next agent:** start **P0.1**

### 2026-06-10 — Stack

- **HTTP router:** chi v5 (≥ v5.2.4); not gin/echo

### 2026-06-10 — Planning (continued)

- **Docker-only dev:** no local Go/Node; all build/test/run via Compose
- **P0.2** expanded to full dev stack (Dockerfiles, api/web/migrate/postgres services)
- P5.1 refocused to production compose + CI (dev compose is P0.2)

### 2026-06-10 — Planning

- Created agent doc set in `agents/` (`DECISIONS.md`, `CONTEXT.md`, `AGENTS.md`, `BUILD_ORDER.md`, `progress.md`)
- Added `.cursor/rules/codencil-agents.mdc` pointing to `agents/`
- Deferred formal approval workflow
- **Next agent:** start **P0.1**

---

## Blockers

*(none)*
