# Contributing to Codencil

Thanks for your interest in Codencil. This project is early-stage; most work is tracked in [`agents/BUILD_ORDER.md`](agents/BUILD_ORDER.md).

## Before you start

1. Read [`agents/CONTEXT.md`](agents/CONTEXT.md) and [`agents/AGENTS.md`](agents/AGENTS.md)
2. **Check open PRs:** `gh pr list --state open` and [`agents/progress.md`](agents/progress.md) stack table
3. **Resume:** no open PRs → branch from `main`; open stack → branch from **tip** (see [`agents/STACK.md`](agents/STACK.md))
4. Pick the next unchecked task in `progress.md`
5. **Docker required** — build, test, and migrate run via Docker Compose only (no local Go/Node needed)

## Branch workflow (stacked PRs)

Codencil uses **stacked branches** so work continues without waiting for merge. Full guide: [`agents/STACK.md`](agents/STACK.md).

| Rule | Detail |
|---|---|
| **Start session** | `gh pr list --state open` — no PRs → `main`; else branch from **stack tip** |
| **First task (empty stack)** | Branch from `main` → PR base **`main`** |
| **Next tasks** | Branch from **tip branch** → PR base **tip branch** |
| **PRs** | Open **ready for review** (not draft); maintainer merges in order |
| **Naming** | `feature/p0.1-scaffold`, `feature/p0.2-docker-stack`, … |
| **One task per PR** | No combining BUILD_ORDER tasks |

```bash
# Stacked example: P0.2 after P0.1
git checkout feature/p0.1-scaffold && git pull
git checkout -b feature/p0.2-docker-stack
git push -u origin feature/p0.2-docker-stack
gh pr create --base feature/p0.1-scaffold --title "P0.2: Docker Compose dev stack"
```

## Pull requests

1. Open a PR ready for review (see stack rules above for base branch; no `--draft`)
2. Fill out the [PR template](.github/pull_request_template.md) completely — include **PR base branch**
3. Include the **BUILD_ORDER task ID** in the title (e.g. `P1.4: document publish API`)
4. Update **`agents/progress.md`** (checkbox + **Open stack** table + session log)
5. Run verification from [`agents/AGENTS.md`](agents/AGENTS.md) via Docker before opening the PR
6. **Verify CI passes** — `gh pr checks <number>`; fix failures on the same branch before starting the next task
7. **Agents do not merge** — maintainer reviews and merges bottom-of-stack first

### PR title format (recommended)

```
P1.4: HTTP document CRUD and publish v1
```

## Verification (Docker only)

```bash
docker compose build

# After API exists:
docker compose run --rm api go test ./...
docker compose run --rm api go vet ./...

# After web exists:
docker compose run --rm web npm run build

# After migrations exist:
docker compose up -d postgres
docker compose run --rm migrate -path /migrations -database "$DATABASE_URL" up
```

CI runs the same checks on pull requests when `docker-compose.yml` and app code are present.

## What not to include in PRs

- Scope outside the claimed BUILD_ORDER task
- OIDC/auth before Phase 4 (unless the task explicitly covers it)
- Formal approval workflow features (deferred — see `agents/DECISIONS.md`)
- Changes to the TypRow library repo (file issues there instead)
- Secrets or `.env` files

## Decisions and architecture

- Product/architecture decisions live in [`agents/DECISIONS.md`](agents/DECISIONS.md)
- Do not change decisions without discussion in the PR
- Implementation order is defined in [`agents/BUILD_ORDER.md`](agents/BUILD_ORDER.md) — do not skip phases

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
