# Stacked branches & PRs

Codencil uses a **stacked PR workflow** so agents can keep building without waiting for review/merge. The maintainer reviews and merges **in order** (bottom of stack first).

---

## Resuming work (start of every session)

**Before picking a task or creating a branch**, check for open PRs:

```bash
gh pr list --state open
```

Also read the **Open stack** table in [`progress.md`](progress.md).

| Open PRs? | Branch from | PR base | Next task |
|---|---|---|---|
| **None** | `main` (pull latest) | `main` | First unchecked task in `progress.md` |
| **One or more** | **Tip of stack** — branch for the highest completed task still open (last row in stack table) | That same branch | Next unchecked task after the tip task |

**Tip of stack** = the branch you stack the *next* PR on. Example: if P0.1 and P0.2 PRs are open, P0.2 is the tip; start P0.3 from `feature/p0.2-docker-stack`.

```bash
git fetch origin

# No open PRs:
git checkout main && git pull origin main
git checkout -b feature/p0.1-scaffold

# Open PRs — resume from tip (example: P0.1 is tip, starting P0.2):
git checkout feature/p0.1-scaffold && git pull origin feature/p0.1-scaffold
git checkout -b feature/p0.2-docker-stack
```

Do **not** assume `main` has the latest work while a stack is open. Do **not** start a new task from an older branch in the middle of the stack.

---

## Rules

| Rule | Detail |
|---|---|
| **One task = one branch = one PR** | e.g. `P0.2` only |
| **Branch from previous task** | `P0.2` branches from `feature/p0.1-*`, not `main` (unless `P0.1` is merged) |
| **PR targets previous branch** | First PR in stack targets `main`; each next PR targets the **previous feature branch** |
| **Open PRs normally** | Do **not** use draft PRs — open ready for review |
| **Do not merge** | Agents open PRs; maintainer merges after review |
| **Merge order** | Bottom → top: merge `P0.1` → `main`, then rebase/retarget `P0.2` → `main`, etc. |

---

## Branch naming

```
feature/p0.1-scaffold
feature/p0.2-docker-stack
feature/p0.3-migrate
feature/p0.4-api-skeleton
feature/p0.5-web-skeleton
feature/p1.1-documents-migration
...
```

Pattern: `feature/{task-id-lowercase}-{short-slug}`

---

## Stack example (Phase 0)

```
main
 └── feature/p0.1-scaffold          PR → base: main
      └── feature/p0.2-docker-stack PR → base: feature/p0.1-scaffold
           └── feature/p0.3-migrate PR → base: feature/p0.2-docker-stack
                └── ...
```

*(Non-BUILD_ORDER docs PRs may sit below the first task PR — see stack table in `progress.md`.)*

---

## Agent workflow (each task)

### 1. Find the base branch

See **[Resuming work](#resuming-work-start-of-every-session)** above first.

| Situation | Branch from | PR base |
|---|---|---|
| No open PRs | `main` | `main` |
| Stack open — starting next task | **Tip branch** (highest open task in stack table) | Tip branch |
| Previous tip merged to `main` mid-stack | Rebase onto `main`; retarget PR to `main` | `main` |

### 2. Create task branch

```bash
git checkout -b feature/p0.2-docker-stack
```

### 3. Implement one BUILD_ORDER task

See [`BUILD_ORDER.md`](BUILD_ORDER.md) deliverables + verify.

### 4. Push and open PR

**First in stack (base = main):**
```bash
git push -u origin feature/p0.1-scaffold
gh pr create --base main \
  --title "P0.1: Repo scaffold" \
  --body-file .github/pull_request_template.md
```

**Stacked on previous branch:**
```bash
git push -u origin feature/p0.2-docker-stack
gh pr create --base feature/p0.1-scaffold \
  --title "P0.2: Docker Compose dev stack" \
  --body-file .github/pull_request_template.md
```

### 5. Update tracking

In **`progress.md`**:
- Check off task when work is complete
- Add row to **Open stack** table (branch, PR base, PR URL, status)

---

## Maintainer: mass review

1. Open PRs bottom of stack → top
2. Review each (CI should run per PR)
3. Merge **P0.1** → `main` first
4. For **P0.2**: click "Edit" on PR → change base to `main`, rebase branch onto `main` if GitHub prompts
5. Repeat until stack is merged

Alternatively: merge entire stack locally with rebase — maintainer choice.

---

## When previous PR is merged mid-stack

If `P0.1` merges while `P0.2` and `P0.3` PRs are open:

```bash
git checkout feature/p0.2-docker-stack
git fetch origin
git rebase origin/main
git push --force-with-lease
# Update PR base to main in GitHub UI
```

---

## CI on stacked PRs

CI runs on each PR against its **base branch**. Stacked PRs validate incremental changes. After retargeting to `main`, CI re-runs against full diff to `main`.

---

## Anti-patterns

| Don't | Do |
|---|---|
| PR directly to `main` for P0.2 while P0.1 PR open | PR base = `feature/p0.1-scaffold` |
| One branch with P0.1 + P0.2 + P0.3 | One branch per task |
| Open PRs as **draft** | Open ready for review |
| Merge your own PRs | Leave open for maintainer |
| Skip updating stack table in `progress.md` | Always record branch + PR link |
| Start from `main` while stack is open | Resume from **tip** branch |
