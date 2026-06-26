# Codencil — Build Order (AI agent tasks)

Granular phases for **one agent session ≈ one task**. Each task has prerequisites, deliverables, and verification.

**Status tracking:** update checkboxes in `agents/progress.md` (not here — this file stays stable).

---

## Principles

0. **Docker-only execution** — the host has **no Go, Node, or migrate installed**. Build, test, migrate, and run **only** via Docker Compose (see `P0.2`).
1. **Backend before frontend** for each feature (API contract first)
2. **Pure logic before IO** — `internal/publish` tests before wiring publish HTTP
3. **Auth last among MVP features** — dev stub unblocks Phases 1–3
4. **Vertical demo early** — Phase 1 ends with visible preview in browser
5. **Do not parallelize** Phases 1–3 tasks out of order

---

## Docker commands (use for all verification)

Do **not** run `go`, `npm`, or `migrate` on the host. Use these patterns:

```bash
# Start stack
docker compose up -d postgres
docker compose run --rm migrate -path /migrations -database "$DATABASE_URL" up

# API tests / vet
docker compose run --rm api go test ./...
docker compose run --rm api go vet ./...

# Web build
docker compose run --rm web npm run build

# Logs / health
docker compose up -d api web
curl http://localhost:8080/health
```

Exact service names and env vars are defined in **P0.2**. Agents must not assume a local toolchain.

---

## Phase 0 — Project shell

| ID | Task | Depends | Deliverables | Verify |
|---|---|---|---|---|
| **P0.1** | Repo scaffold | — | Module layout, `.gitignore`, MIT `LICENSE`, root `README.md` | Files present; valid `apps/api/go.mod` (no host `go` required) |
| **P0.2** | **Docker Compose dev stack** | P0.1 | `docker-compose.yml`, `apps/api/Dockerfile`, `apps/web/Dockerfile`, migrate service, dev `api` + `web` services with volume mounts, `Makefile` or `scripts/dev.*` | `docker compose build` succeeds; `docker compose run --rm api go version`; postgres healthy |
| **P0.3** | go-migrate wiring | P0.2 | `db/migrations/` (can be empty seed), migrate service config | `docker compose run --rm migrate … up` / `down 1` / `up` |
| **P0.4** | API skeleton | P0.2 | chi v5 router, health `GET /health`, config from env | `docker compose up -d api` → `curl localhost:8080/health` → 200 |
| **P0.5** | Web skeleton | P0.2 | Next.js app in `apps/web`, env for API URL | `docker compose up -d web` → home page loads in browser |

**Phase 0 exit:** Full dev stack runs via **`docker compose` only**; Postgres migrates; empty API + web reachable. No product features yet.

---

## Phase 1 — Read path (document → publish v1 → preview)

| ID | Task | Depends | Deliverables | Verify |
|---|---|---|---|---|
| **P1.1** | Migration: documents + versions | P0.3 | `000001_documents.up.sql` — `documents`, `document_versions` | migrate up/down via compose |
| **P1.2** | TypRow models + store | P1.1 | Create/read/update draft; insert version on publish | `docker compose run --rm api go test ./internal/...` |
| **P1.3** | `internal/publish` diff scaffold | P0.2 | Package + test file; line-diff util (no anchors yet) | `docker compose run --rm api go test ./internal/publish/...` |
| **P1.4** | HTTP: document CRUD + publish v1 | P1.2, P0.4 | `POST/GET/PATCH /documents`, `POST .../publish` | curl/create/publish/get v1 |
| **P1.5** | Web: markdown preview page | P1.4, P0.5 | Route `/documents/[id]/versions/[v]` + react-markdown | Browser shows published md |
| **P1.6** | Phase 1 smoke doc | P1.5 | Seed script or curl script in `scripts/` | Script runs end-to-end |

**Phase 1 exit:** Author can create doc, edit draft, publish v1, view rendered preview. **No comments.**

---

## Phase 2 — Review path (anchors + threads)

| ID | Task | Depends | Deliverables | Verify |
|---|---|---|---|---|
| **P2.1** | Migration: threads, comments, anchors | P1.1 | `000002_review.up.sql` | migrate up/down |
| **P2.2** | TypRow models for review entities | P2.1 | Store: create thread, anchor on version, list by version | `docker compose run --rm api go test ./internal/...` |
| **P2.3** | HTTP: comment API | P2.2, P1.4 | Create anchor+thread, add comment, list, resolve | curl tests |
| **P2.4** | Web: text selection → anchor | P2.3, P1.5 | Select text in preview; POST anchor; show in margin | manual UI test |
| **P2.5** | Web: thread panel UI | P2.4 | Reply, resolve, collapsed resolved state | manual UI test |

**Phase 2 exit:** Reviewer can comment on **published version** with Word-like margin UX.

---

## Phase 3 — Publish v2 + anchor migration

| ID | Task | Depends | Deliverables | Verify |
|---|---|---|---|---|
| **P3.1** | Anchor remap logic | P1.3, P2.1 | `RemapAnchors(old, new, anchors)` + table tests | `docker compose run --rm api go test ./internal/publish/...` |
| **P3.2** | Wire publish v2+ | P3.1, P2.2 | Publish increments version; copies/remaps anchors | test: edit line → shifted |
| **P3.3** | HTTP + UI version selector | P3.2 | List versions; view older; anchors per version | manual: v1 vs v2 anchors |
| **P3.4** | Orphaned anchor UX | P3.3 | UI shows orphaned/shifted badges | delete quoted text → orphaned |

**Phase 3 exit:** Full draft → publish → review → revise → publish loop works.

---

## Phase 4 — Auth & roles

| ID | Task | Depends | Deliverables | Verify |
|---|---|---|---|---|
| **P4.1** | Users table + dev auth middleware | P3.2 | Map `X-Dev-User-Id` to user row | requests attributed |
| **P4.2** | Document membership / roles | P4.1 | owner / reviewer / viewer on document | forbidden for wrong role |
| **P4.3** | OIDC (optional MVP+) | P4.2 | go-oidc behind env flag | login flow when enabled |
| **P4.4** | Remove or gate dev auth in prod | P4.3 | Document `AUTH_MODE=dev\|oidc` | prod compose uses oidc |

**Phase 4 exit:** Real users; reviewers cannot publish; viewers read-only.

---

## Phase 5 — Self-host polish

| ID | Task | Depends | Deliverables | Verify |
|---|---|---|---|---|
| **P5.1** | Production compose + CI | P4.1 min | Prod Docker targets, compose profiles, GitHub Actions using compose | CI green without host Go/Node |
| **P5.2** | README self-host guide | P5.1 | env vars, quickstart (**Docker only**), architecture diagram | human can follow with Docker installed only |
| **P5.3** | GitHub publish | P5.2 | Push to `TheBlackHowling/codencil` | CI green |

**Phase 5 exit:** OSS-ready self-host bundle.

---

## Dependency graph (simplified)

```
P0.1 scaffold
  └─► P0.2 docker dev stack (required before any go/npm on host)
        └─► P0.3 migrate ─► P0.4 api ─► P0.5 web
  └─► P1.1 schema ─► P1.2 store ─► P1.4 HTTP ─► P1.5 web preview
         └─► P1.3 publish (parallel with P1.2, after P0.2)
  └─► P2.* review (after P1.4)
  └─► P3.* anchor migration (after P2.*, needs P1.3)
  └─► P4.* auth (after P3.*)
  └─► P5.* polish (prod compose + CI; dev compose done in P0.2)
```

---

## Task sizing guide

| Too big for one session | Right size |
|---|---|
| "Build Phase 2" | P2.3 HTTP comment API only |
| "Full stack auth" | P4.1 dev auth middleware only |
| "Complete frontend" | P1.5 preview page only |
| "Implement Codencil" | Next unchecked task in `agents/progress.md` |

---

## API surface (reference — implement incrementally)

**Phase 1:**
- `GET /health`
- `POST /documents` — create
- `GET /documents/{id}` — draft + metadata
- `PATCH /documents/{id}` — update draft markdown/title
- `POST /documents/{id}/publish` — publish draft → new version
- `GET /documents/{id}/versions/{v}` — version snapshot

**Phase 2:** add
- `GET /documents/{id}/versions/{v}/anchors` — list with threads
- `POST /documents/{id}/versions/{v}/anchors` — create anchor + thread
- `POST /threads/{id}/comments`
- `POST /anchors/{id}/resolve` / `reopen`

(Exact paths may adjust — update this section if they do, one line per change.)
