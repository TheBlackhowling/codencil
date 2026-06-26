# Contributing to Codencil

Thanks for your interest in Codencil. This project is early-stage; most work is tracked in [`agents/BUILD_ORDER.md`](agents/BUILD_ORDER.md).

## Before you start

1. Read [`agents/CONTEXT.md`](agents/CONTEXT.md) and [`agents/AGENTS.md`](agents/AGENTS.md)
2. Check [`agents/progress.md`](agents/progress.md) for the next task
3. **Docker required** — build, test, and migrate run via Docker Compose only (no local Go/Node needed)

## Branch workflow

- Branch from **`main`**
- Name branches: `feature/p0.2-docker-stack`, `feature/p1.4-publish-api`, etc.
- **One BUILD_ORDER task per PR** (e.g. `P1.4` only — no drive-by refactors)
- Keep PRs small and reviewable

## Pull requests

1. Open a PR against **`main`**
2. Fill out the [PR template](.github/pull_request_template.md) completely
3. Include the **BUILD_ORDER task ID** in the title or description (e.g. `P1.4: document publish API`)
4. Update **`agents/progress.md`** in the same PR (checkbox + session log) when implementation work is done
5. Run verification from [`agents/AGENTS.md`](agents/AGENTS.md) via Docker before requesting review

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
