# Codencil — Implementation Progress

> **Agents:** Update this file at the end of every session.  
> **Source of truth for "what's next"** — see `BUILD_ORDER.md` for task definitions.

**Last updated:** 2026-06-10  
**Current phase:** 0 (in progress)  
**Next task:** **P0.2** — Docker Compose dev stack  
**Stack policy:** Stacked draft PRs — see [`STACK.md`](STACK.md)

---

## Open stack (branches & PRs)

| Task | Branch | PR base | PR | Status |
|---|---|---|---|---|
| docs (stack workflow) | `feature/docs-stacked-pr-workflow` | `main` | [#1](https://github.com/TheBlackhowling/codencil/pull/1) | draft |
| P0.1 | `feature/p0.1-scaffold` | `feature/docs-stacked-pr-workflow` | [#2](https://github.com/TheBlackhowling/codencil/pull/2) | draft |

*Agents: add a row when opening each draft PR. Remove or mark merged after maintainer merge.*

---

## Phase 0 — Project shell

- [x] **P0.1** Repo scaffold
- [ ] **P0.2** Docker Compose dev stack (Go + Node in containers)
- [ ] **P0.3** go-migrate wiring
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

### 2026-06-10 — P0.1 scaffold (stacked)

- Branch `feature/p0.1-scaffold` from `feature/docs-stacked-pr-workflow`
- Added `apps/api/go.mod`, module layout (`internal/*`, `cmd/codencil`, `db/migrations`, `apps/web` placeholders)
- Draft PR stacked on docs PR #1
- **Next agent:** **P0.2** on branch from `feature/p0.1-scaffold`

### 2026-06-10 — Stacked PR workflow

- Added `agents/STACK.md`; CONTRIBUTING/AGENTS/PR template updated for draft stacked PRs
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
