# Codencil — Implementation Progress

> **Agents:** Update this file at the end of every session.  
> **Source of truth for "what's next"** — see `BUILD_ORDER.md` for task definitions.

**Last updated:** 2026-07-09  
**Current phase:** Phase 3 complete (stack open for merge)  
**Next task:** **P4.1** — Users table + dev auth middleware  
**Stack policy:** Stacked PRs (ready for review, not draft) — see [`STACK.md`](STACK.md)

---

## Open stack (branches & PRs)

| Task | Branch | PR base | PR | Status |
|---|---|---|---|---|
| P3.0 | `feature/p3.0-docs-status` | `main` | *(pending)* | open |

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
- [x] **P2.5** Web: thread panel UI

## Phase 3 — Publish v2 + anchor migration

- [x] **P3.0** Docs status sync (Phase 2 merged → Phase 3 kickoff)
- [x] **P3.1** Anchor remap logic
- [x] **P3.2** Wire publish v2+
- [x] **P3.3** HTTP + UI version selector
- [x] **P3.4** Orphaned anchor UX

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

### 2026-07-09 — Phase 3 publish v2 + anchor migration (stack starting)

- P3.0–P3.4: remap logic, publish wiring, version selector, status badges
- **Next agent:** merge stack #19+ bottom → top, then **P4.1**

### 2026-07-09 — Phase 2 merged; Phase 3 kickoff

- Phase 2 stack #13–#18 merged to `main`
- No open PRs; local repo synced to `main`
- **Next agent:** P3.1 from tip after P3.0 PR opens

### 2026-07-09 — Phase 2 review path (stack #13–#18)

- P2.0–P2.5: full review path (migration → store → API → web)
- **Merged:** #13 → #18

### 2026-06-10 — Phase 1 read path (stack #7–#12)

- P1.1–P1.6: document CRUD, publish v1, preview, smoke script
- **Merged:** #7 → #12

---

## Blockers

*(none)*
