## BUILD_ORDER task

<!-- e.g. P1.4 — required for implementation PRs -->

**Task ID:**

## Stack

**Branch:** `feature/p0.2-docker-stack`  
**PR base:** `feature/p0.1-scaffold` *(or `main` if first in stack)*

## Summary

<!-- What does this PR do? One BUILD_ORDER task only. -->

## Changes

<!-- Bullet list of meaningful changes -->

-

## Test plan

<!-- All commands via Docker Compose — check boxes you ran -->

- [ ] `docker compose build`
- [ ] `docker compose run --rm api go test ./...` (if API touched)
- [ ] `docker compose run --rm api go vet ./...` (if API touched)
- [ ] `docker compose run --rm web npm run build` (if web touched)
- [ ] Manual smoke (describe): 

## Progress tracking

- [ ] Updated `agents/progress.md` (task checkbox + session log)

## Notes

<!-- Blockers, follow-ups, or decisions that need review -->
