# Codencil — Implementation Progress

> **Agents:** Update this file at the end of every session.  
> **Source of truth for "what's next"** — see `BUILD_ORDER.md` for task definitions.

**Last updated:** 2026-07-09  
**Current phase:** Phase 2 complete (stack open for merge)  
**Next task:** **P3.1** — Anchor remap logic  
**Stack policy:** Stacked PRs (ready for review, not draft) — see [`STACK.md`](STACK.md)

---

## Open stack (branches & PRs)

| Task | Branch | PR base | PR | Status |
|---|---|---|---|---|
| P2.0 | `feature/p2.0-docs-status` | `main` | [#13](https://github.com/TheBlackhowling/codencil/pull/13) | open |
| P2.1 | `feature/p2.1-review-migration` | `feature/p2.0-docs-status` | [#14](https://github.com/TheBlackhowling/codencil/pull/14) | open |
| P2.2 | `feature/p2.2-review-store` | `feature/p2.1-review-migration` | [#15](https://github.com/TheBlackhowling/codencil/pull/15) | open |
| P2.3 | `feature/p2.3-comment-api` | `feature/p2.2-review-store` | *(pending)* | open |
| P2.4 | `feature/p2.4-web-selection` | `feature/p2.3-comment-api` | *(pending)* | open |

*Agents: add a row when opening each PR. Remove or mark merged after maintainer merge.*

**Resume rule:** at session start, run `gh pr list --state open`. No open PRs → `main`. Otherwise → checkout **tip** branch.

---

## Phase 0 — Project shell

- [x] **P0.1** Repo scaffold
- [x] **P0.2** Docker Compose dev stack (Go + Node in containers)
- [x] **P0.3** go-migrate wiring
- [x] **P0.4** API skeleton
- [x] **P0.5** Web skeleton

## Phase 1 — Read path

- [x] **P1.1** Migration: documents + versions
- [x] **P1.2** TypRow models + store
- [x] **P1.3** `internal/publish` diff scaffold
- [x] **P1.4** HTTP: document CRUD + publish v1
- [x] **P1.5** Web: markdown preview page
- [x] **P1.6** Phase 1 smoke doc

## Phase 2 — Review path

- [x] **P2.0** Docs status sync (Phase 1 merged → Phase 2 kickoff)
- [x] **P2.1** Migration: threads, comments, anchors
- [x] **P2.2** TypRow models for review entities
- [x] **P2.3** HTTP: comment API
- [x] **P2.4** Web: text selection → anchor
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

### 2026-07-09 — P2.0 docs status sync

- Synced `progress.md`, `README.md`, `CONTEXT.md` after Phase 1 stack merged (#7–#12)
- Open stack table reset for Phase 2
- **Next agent:** P2.1 from tip of stack after P2.0 PR opens

### 2026-06-10 — Phase 1 read path (stack #7–#12)

- P1.1: `000002_documents` migration
- P1.2: TypRow document store + integration tests
- P1.3: `internal/publish` line diff scaffold
- P1.4: document CRUD + publish HTTP API
- P1.5: Next.js markdown preview route
- P1.6: `scripts/phase1-smoke.ps1`
- **Next agent:** merge stack bottom → top, then **P2.1**

### 2026-06-10 — P0.5 Web skeleton

- Next.js 15 App Router home page with NEXT_PUBLIC_API_URL display
- Web dev service: npm install + next dev in compose

---

## Blockers

*(none)*
