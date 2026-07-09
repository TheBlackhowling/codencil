# Codencil — Implementation Progress

> **Agents:** Update this file at the end of every session.  
> **Source of truth for "what's next"** — see `BUILD_ORDER.md` for task definitions.

**Last updated:** 2026-07-09  
**Current phase:** MVP complete (Phases 0–5); stack open for merge  
**Next task:** *(none — MVP shipped)*  
**Stack policy:** Stacked PRs (ready for review, not draft) — see [`STACK.md`](STACK.md)

---

## Open stack (branches & PRs)

| Task | Branch | PR base | PR | Status |
|---|---|---|---|---|
| P3.4 | `feature/p3.4-anchor-status-ux` | `feature/p3.3-version-selector` | *(merged or open)* | — |
| P4.1 | `feature/p4.1-dev-auth` | `feature/p3.4-anchor-status-ux` | *(open PR)* | open |
| P4.2 | `feature/p4.2-document-roles` | `feature/p4.1-dev-auth` | *(open PR)* | open |
| P4.3 | `feature/p4.3-oidc` | `feature/p4.2-document-roles` | *(open PR)* | open |
| P4.4 | `feature/p4.4-auth-gate` | `feature/p4.3-oidc` | *(open PR)* | open |
| P5.1 | `feature/p5.1-prod-compose` | `feature/p4.4-auth-gate` | *(open PR)* | open |
| P5.2 | `feature/p5.2-self-host-readme` | `feature/p5.1-prod-compose` | *(open PR)* | open |
| P5.3 | `feature/p5.3-github-publish` | `feature/p5.2-self-host-readme` | *(open PR)* | open |

*Merge order: P3 stack (if not merged) → P4.1 → … → P5.3*

---

## Phase 0 — Project shell

- [x] **P0.1** – **P0.5** *(complete)*

## Phase 1 — Read path

- [x] **P1.1** – **P1.6** *(complete)*

## Phase 2 — Review path

- [x] **P2.0** – **P2.5** *(complete)*

## Phase 3 — Publish v2 + anchor migration

- [x] **P3.0** – **P3.4** *(complete)*

## Phase 4 — Auth & roles

- [x] **P4.1** Users table + dev auth middleware
- [x] **P4.2** Document membership / roles
- [x] **P4.3** OIDC (optional MVP+)
- [x] **P4.4** Remove or gate dev auth in prod

## Phase 5 — Self-host polish

- [x] **P5.1** Production compose + CI
- [x] **P5.2** README self-host guide
- [x] **P5.3** GitHub publish — repo live; CI on PRs

---

## Session log

### 2026-07-09 — Phases 4–5 (stack on P3.4 tip)

- P4.1–P4.4: users, roles, OIDC, AUTH_MODE gating
- P5.1–P5.3: prod compose, README self-host guide, MVP status
- **Stack tip:** `feature/p5.3-github-publish`

### 2026-07-09 — Phase 3 publish v2 + anchor migration

- P3.0–P3.4: remap, publish wiring, version selector, status badges

---

## Blockers

*(none)*
