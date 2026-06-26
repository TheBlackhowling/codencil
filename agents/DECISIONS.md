# Codencil — Planning & Architecture Decisions

> **Purpose:** Stable product & architecture decisions (changes rarely).  
> **For implementation order:** see [`BUILD_ORDER.md`](BUILD_ORDER.md)  
> **For agent rules:** see [`AGENTS.md`](AGENTS.md)  
> **Last updated:** 2026-06-10  
> **Status:** Planning complete; implementation not started (see [`progress.md`](progress.md))

---

## What this document is (and is not)

| In scope | Out of scope (separate thread) |
|---|---|
| **Codencil** — markdown document review platform | **TypRow** rename (`typedb` → `typrow`) — see [TypRow note](#typrow-separate-project) at bottom |

---

## One-line pitch

**Codencil** is a self-hosted markdown review platform: authors work in `.md` and publish versioned snapshots; reviewers comment in the margin like Word—not in raw markdown.

**Tagline direction:** *Write in markdown. Review in the margin. Publish when it's ready.*

---

## Publisher & repos

| Item | Decision |
|---|---|
| **GitHub org** | [TheBlackHowling](https://github.com/TheBlackHowling) |
| **Main repo (planned)** | `TheBlackHowling/codencil` — **greenfield, not a fork** |
| **License (planned)** | MIT or Apache-2.0 (TBD; Markdown Viewer is Apache-2.0 if any code borrowed) |
| **Business model** | Free OSS, self-hosted; optional public demo — **not charging** |

---

## Problem & audience

### Authors
- Work in markdown (often with AI agents)
- Need explicit **draft → publish** workflow (not live-edit chaos)
- Want reviewers to engage without learning markdown

### Reviewers / managers
- Word-like **inline comments in the margin**
- Comment **threads** with open / resolved state
- Resolved comments: **collapsed but visible** (Word-like)
- Do **not** want to comment inside raw `.md`

### Operators
- **Self-host first** (Docker + Postgres)
- **Docker-only dev** — no local Go/Node on maintainer machine; build/test/run via Compose
- SSO/OIDC for enterprise
- Optional public demo under TheBlackHowling later

---

## Inspiration vs fork — Markdown Viewer

**Inspired by:** [markdownviewer.pages.dev](https://markdownviewer.pages.dev) / [ThisIs-Developer/Markdown-Viewer](https://github.com/ThisIs-Developer/Markdown-Viewer) (Apache-2.0)

**Decision: do NOT fork or extend Markdown Viewer long-term.**

| Markdown Viewer | Codencil |
|---|---|
| 100% client-side | Server-backed |
| localStorage / URL share | Postgres source of truth |
| Editor + preview tabs | Publish versions + margin review UI |
| No auth, no comments | Auth, roles, comment threads |
| Vanilla JS static app | Go API + Next.js |

**What to reuse:**
- **Ideas:** GFM preview, split-pane layout, Mermaid/math rendering patterns
- **Optional:** small Apache-licensed snippets with README/NOTICE credit
- **Do not:** fork the repo or inherit its product identity (trademark rebrand required)

**Rendering:** keep markdown rendering in the **frontend** (`react-markdown` + remark/rehype). Backend stores raw markdown snapshots only.

---

## Core domain model

### Documents
- Stored in **Postgres** (not file mount as source of truth)
- Each document has a **draft** markdown body and metadata (title, org, etc.)

### Versions
- **Explicit publish:** draft → publish v1, v2, v3…
- Each published version is a **full markdown snapshot** (immutable once published)
- Draft can continue editing after publish

### Comment threads
- Threads belong to the **document** (stable identity across versions)
- Individual **comments** belong to a thread

### Anchors (per version)
- Each anchor is tied to a **specific published version**
- Fields (conceptual):
  - `document_id`, `version`, `anchor_id` (stable within document)
  - Line range: `start_line`, `end_line`
  - `quoted_text` — **required**; do not rely on line numbers alone
  - `anchor_status`: `active` | `shifted` | `orphaned`
- **Review state per anchor:** `open` | `resolved` (+ `resolved_by`, `resolved_at`, reopen tracking)

### Publish v(N) → v(N+1) anchor migration
1. Diff old snapshot → new snapshot
2. Remap line ranges using diff + fuzzy match on `quoted_text`
3. Copy anchors + review state forward to new version
4. Mark `shifted` / `orphaned` where text cannot be relocated

**Implementation note:** anchor remapping logic belongs in a pure Go `internal/publish` package with table-driven tests, independent of HTTP/DB where possible.

---

## Auth & roles

| Role | Capabilities (high level) |
|---|---|
| **Owner** | CRUD document, publish, manage access |
| **Reviewer** | View published versions, create/respond/resolve threads |
| **Viewer** | Read-only |

- **Auth required** for commenting and ownership
- **Self-host:** OIDC/SSO
- **Public demo (later):** OAuth

---

## Tech stack (decided)

| Layer | Choice | Notes |
|---|---|---|
| **Backend** | **Go** | User is stronger in Go than Python; good self-host story (single binary) |
| **HTTP router** | **chi** (`github.com/go-chi/chi/v5` ≥ v5.2.4) | `net/http` style; no `RedirectSlashes`; see `AGENTS.md` |
| **Database** | **Postgres** | Documents, versions, anchors, threads, users |
| **Migrations** | **go-migrate** | User preference (not goose) |
| **Data access** | **[TypRow](https://github.com/TheBlackHowling/typrow)** (`github.com/TheBlackHowling/typrow`) | Type-safe SQL library; SQL-first, not ORM |
| **Frontend** | **Next.js** | Margin-comment UX, markdown preview |
| **Deploy** | `docker compose` | `migrate` → `api` → `web` → `postgres` |
| **Dev environment** | **Docker Compose only** | Host has no Go/Node; see `BUILD_ORDER` **P0.2** |

**Explicitly rejected for Codencil backend:** Python/FastAPI (despite pcaicode familiarity), forking Markdown Viewer, TypeDB/Vaticle graph DB, **gin/echo** (using chi).

---

## Repo layout (planned)

```
codencil/
  db/
    migrations/           # go-migrate .up.sql / .down.sql
  apps/
    api/
      cmd/codencil/
      internal/
        domain/           # Document, Version, Anchor, Thread
        models/           # TypRow models + Query* methods
        publish/          # diff + anchor remapping (pure Go)
        auth/             # OIDC middleware
        http/             # handlers + routes
    web/                  # Next.js — preview + margin comments
  docker-compose.yml
```

---

## MVP build order (summary)

Full granular tasks: **`BUILD_ORDER.md`** (Phases 0–5, one agent task per row).

| Phase | Outcome |
|---|---|
| **0** | Repo shell + **Docker dev stack (P0.2)** — all tooling in containers |
| **1** | Create doc → publish v1 → preview in browser (**no comments**) |
| **2** | Margin comments + threads on published version |
| **3** | Publish v2+ with anchor remapping |
| **4** | Auth & roles (dev stub first, OIDC optional) |
| **5** | Self-host docker compose + README |

Start **`internal/publish` tests** in P1.3 (before anchor migration in P3).

---

## Agent documentation (repo layout)

```
codencil/
  agents/
    README.md        ← index
    CONTEXT.md       ← read first (orientation)
    AGENTS.md        ← rules, defaults, verification
    BUILD_ORDER.md   ← phased tasks (P0.1 … P5.3)
    DECISIONS.md     ← this file (why)
    progress.md      ← living status (update every session)
  .cursor/rules/     ← Cursor points to agents/
```

## Competitive landscape (no exact match)

| Product | Overlap | Codencil differentiator |
|---|---|---|
| **Draftmark** | Hosted MD + comments + versions | Self-host + SSO |
| **markview.io** | MD viewer for non-devs | No server comments/versions |
| **Markdown Viewer** | Client preview/editor | Server-backed review platform |
| CollabMD, Lumen, Commentary, md-redline, Crit | Partial | Combined: versions + margin + self-host |

**Codencil differentiators:** self-hosted + SSO, explicit publish versions, per-version anchor review state, manager-first Word-like UX, DB-stored docs.

---

## Naming

### Product name: **Codencil**

- Invented spelling evoking *codicil* (document amendments/versions)
- 0 GitHub/npm collisions at time of search
- Repo: `TheBlackHowling/codencil`

### Rejected names (do not reuse)

Redraft, Notedown, Upmark, UpNote, Markview, Viewmark, NoteUp, Passnote, Margentum, Codicil — various collision/trademark/SEO issues documented in planning sessions.

---

## Domains (planned, not purchased yet)

| Domain | Status | Plan |
|---|---|---|
| **codencil.dev** | Likely available | Primary — docs, API, `demo.codencil.dev` |
| **codencil.io** | Likely available | Marketing / manager-facing |
| **codencil.com** | Registered (~$6.4k premium) + historical malware subdomain reports | **Skip** unless vetted |
| **getcodencil.com** etc. | Likely available | Optional cheap `.com` redirect later |

**Registrar:** AWS Route 53 (user preference).  
**Until domains bought:** GitHub-only is fine; README can use placeholder URLs.

---

## URLs (target, when live)

```
Website:   https://codencil.dev
Demo:      https://demo.codencil.dev
Docs:      https://docs.codencil.dev
GitHub:    https://github.com/TheBlackHowling/codencil
```

---

## README acknowledgment (when shipping)

> Codencil is a self-hosted markdown review platform. Preview rendering is inspired by [Markdown Viewer](https://github.com/ThisIs-Developer/Markdown-Viewer) (Apache-2.0).

If any code is copied: add `NOTICE` file with Apache attribution.

---

## Schema sketch (first migration — draft)

Tables to plan (names TBD at implementation):

- `documents` — id, org_id, title, draft_markdown, timestamps
- `document_versions` — document_id, version, markdown, published_at, published_by
- `comment_threads` — id, document_id, created_at, …
- `comments` — id, thread_id, author_id, body, created_at, …
- `version_anchors` — document_id, version, anchor_id, start_line, end_line, quoted_text, anchor_status, review_state, thread_id, …

Use TypRow models with `Load`/`Insert`/`Update` for entities and `QueryAll[T]` for review-page JOIN read models.

---

## What is NOT started yet

- [ ] GitHub repo `TheBlackHowling/codencil`
- [ ] Initial migration + TypRow models
- [ ] `internal/publish` anchor remapping prototype
- [ ] Next.js margin-comment UI
- [ ] Docker compose self-host bundle
- [ ] Domain registration
- [ ] Architecture doc with full anchor migration pseudocode (optional follow-up)

---

## TypRow (separate project)

The **typedb → TypRow** rename is a **different initiative** (Go SQL library under `TheBlackHowling/typrow`).

**Decisions already made there:**
- New repo `typrow`; old `typedb` repo **frozen** for work OSS approval and internal tools
- `typedb-examples` → `typrow-examples` for new docs
- Codencil should depend on **`github.com/TheBlackHowling/typrow`**, not `typedb`

Do not conflate Codencil product work with TypRow library maintenance unless explicitly scoped.

---

## Agent bootstrap checklist

1. Read **`agents/CONTEXT.md`** → **`agents/AGENTS.md`** → **`agents/progress.md`**
2. Pick **one** task: the first unchecked item in `agents/progress.md`
3. Execute per **`agents/BUILD_ORDER.md`** (deliverables + verify column)
4. Update **`agents/progress.md`** + session log before ending
5. Do **not** buy domains, change TypRow, or implement deferred features unless asked

---

## Deferred (explicitly out of scope for now)

### Formal approval flow

**Decision:** Skip for initial implementation. The existing model is enough:

- Author **publish** when ready (explicit v1, v2, v3…)
- Reviewers **comment in margin** + **resolve threads**
- No submit/approve gates, no `proposed` vs `approved` version states, no resolve-all blocking

Revisit post-MVP only if a concrete compliance/sign-off requirement appears.

---

## Resolved defaults (for agents — see AGENTS.md)

| Question | Decision |
|---|---|
| HTTP router | **chi** (`github.com/go-chi/chi/v5` ≥ **v5.2.4**) |
| License | **MIT** |
| Frontend | **`apps/web`** (Next.js App Router) |
| WebSockets | Defer |
| Multi-tenant | **`org_id` on tables**; single-org dev OK |
| Auth timing | Dev stub Phases 1–3; OIDC Phase 4 |
| Local toolchain | **None — Docker Compose only** |

## Open questions (still TBD)

- Exact OIDC provider assumptions for self-host docs
- Dogfood doc choice (e.g. copy of internal LLM usage guide)
