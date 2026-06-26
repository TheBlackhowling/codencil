# Codencil — Context Guide (for humans & AI agents)

## Purpose

Orientation doc for new sessions. **Read this first**, then follow links based on what you're doing.

**Current phase:** Pre-implementation planning → Phase 0 scaffold  
**Product:** Self-hosted markdown review (margin comments, versioned publish)

---

## Read order (new agent session)

| Order | File | When |
|---|---|---|
| 1 | **`CONTEXT.md`** (this file) | Always — orientation |
| 2 | **`AGENTS.md`** | Before writing code — rules & defaults |
| 3 | **`DECISIONS.md`** | When unsure *why* — architecture & product |
| 4 | **`BUILD_ORDER.md`** | Before implementing — *what* to do next |
| 5 | **`agents/progress.md`** | Check/update task status |

**Do not** start a task marked 🔒 blocked or ✅ done in `agents/progress.md`.

---

## Doc map

| File | Role |
|---|---|
| [`README.md`](./README.md) | Index (this folder) |
| `DECISIONS.md` | Stable product/architecture decisions (change rarely) |
| `BUILD_ORDER.md` | Phased implementation plan (granular agent tasks) |
| `AGENTS.md` | How to work in this repo; verification; anti-patterns |
| `progress.md` | Living checklist — **update at end of every session** |
| `../README.md` | Public-facing project intro (create at Phase 0) |
| `../CONTRIBUTING.md` | Branch workflow, PR rules, verification |

---

## Stack (quick reference)

- **Go API** — chi, Postgres, go-migrate, [TypRow](https://github.com/TheBlackHowling/typrow)
- **Next.js web** — preview + margin comments; markdown rendered client-side
- **Deploy** — docker compose (Postgres + migrate + api + web)

See `DECISIONS.md` for domain model (documents, versions, anchors, threads).

---

## Out of scope (do not expand unless asked)

- Formal approval workflow (`proposed` / `approved` gates)
- TypRow library changes (separate repo: `c:\source\typrow`)
- Forking Markdown Viewer
- Domain purchase / GitHub org setup (unless task says so)
- OIDC/SSO before Phase 4 (use dev auth stub earlier)
- Installing Go/Node on host for build or test (use Docker Compose)

---

## End-of-session handoff

Every agent session that writes code **must**:

1. Update **`agents/progress.md`** — mark task done/in-progress, note blockers
2. Leave a one-paragraph **session note** at bottom of `agents/progress.md` (what changed, what's next)
3. Run verification from **`AGENTS.md`** for the phase touched

---

## Related repos

| Repo | Relationship |
|---|---|
| `TheBlackHowling/typrow` | DB access library — depend on, don't modify in Codencil tasks |
| `TheBlackHowling/typedb` | Legacy — do not use |
| `TheBlackHowling/codencil` | Target GitHub remote (may match this folder later) |
