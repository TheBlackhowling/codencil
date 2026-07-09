# Codencil

Self-hosted markdown review platform — authors publish versioned documents; reviewers comment in the margin.

**Status:** MVP feature-complete (Phases 0–5). See [`agents/progress.md`](agents/progress.md).

**License:** [MIT](LICENSE)

**Contributing:** see [CONTRIBUTING.md](CONTRIBUTING.md) — one BUILD_ORDER task per PR, Docker-only verify.

**Requirement:** Docker + Docker Compose only (no local Go or Node needed).

## Quickstart (development)

```bash
git clone https://github.com/TheBlackHowling/codencil.git
cd codencil
cp .env.example .env
docker compose up -d postgres
docker compose run --rm migrate -path /migrations -database "$DATABASE_URL" up
docker compose up -d api web
```

- API: http://localhost:8080/health  
- Web: http://localhost:3000  
- Smoke test: `./scripts/phase1-smoke.ps1` (PowerShell; stack must be running)

Dev auth uses `X-Dev-User-Id` header (defaults to `dev-user`).

## Self-host (production)

Production compose builds optimized API + web images and expects **OIDC** (`AUTH_MODE=oidc`).

```bash
cp .env.example .env
# Edit .env — set POSTGRES_PASSWORD, DATABASE_URL, OIDC_ISSUER, OIDC_CLIENT_ID, NEXT_PUBLIC_API_URL

docker compose -f docker-compose.prod.yml up -d --build
```

### Environment variables

| Variable | Dev | Prod | Description |
|---|---|---|---|
| `DATABASE_URL` | required | required | Postgres connection string |
| `AUTH_MODE` | `dev` | `oidc` | `dev` = header stub; `oidc` = Bearer JWT |
| `ORG_ID` | `dev` | your org | Tenant id on users/documents |
| `OIDC_ISSUER` | — | required | IdP issuer URL |
| `OIDC_CLIENT_ID` | — | required | OAuth client id for API token validation |
| `NEXT_PUBLIC_API_URL` | required | required | Public API URL for the web app |
| `POSTGRES_PASSWORD` | optional | required | Postgres password (prod compose) |

### Architecture

```
┌─────────┐     ┌─────────┐     ┌──────────┐
│  Web    │────▶│   API   │────▶│ Postgres │
│ Next.js │     │   Go    │     │          │
└─────────┘     └─────────┘     └──────────┘
     │                │
     │  markdown      │  documents, versions,
     │  preview       │  anchors, threads, users
     └────────────────┘
```

**Flow:** Author edits draft → publish v1/v2… → reviewer selects text on a **published version** → margin threads → re-publish remaps anchors (shifted/orphaned).

### Roles

| Role | Capabilities |
|---|---|
| **owner** | Edit draft, publish, read, comment |
| **reviewer** | Read published versions, comment, resolve threads |
| **viewer** | Read only |

## Layout

```
apps/api/          Go API (chi + TypRow + Postgres)
apps/web/          Next.js frontend
db/migrations/     go-migrate SQL files
agents/            Planning & agent docs
docker-compose.yml       Development stack
docker-compose.prod.yml  Production stack
```

```
Write in markdown. Review in the margin. Publish when it's ready.
```

## Acknowledgment

Preview rendering is inspired by [Markdown Viewer](https://github.com/ThisIs-Developer/Markdown-Viewer) (Apache-2.0).
