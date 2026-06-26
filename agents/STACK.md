# Stacked branches & PRs

Codencil uses a **stacked PR workflow** so agents can keep building without waiting for review/merge. The maintainer reviews and merges **in order** (bottom of stack first).

---

## Rules

| Rule | Detail |
|---|---|
| **One task = one branch = one PR** | e.g. `P0.2` only |
| **Branch from previous task** | `P0.2` branches from `feature/p0.1-*`, not `main` (unless `P0.1` is merged) |
| **PR targets previous branch** | First PR in stack targets `main`; each next PR targets the **previous feature branch** |
| **Draft PRs OK** | Open as **draft** until maintainer batch review |
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
 └── feature/p0.1-scaffold          PR #1 → base: main
      └── feature/p0.2-docker-stack PR #2 → base: feature/p0.1-scaffold
           └── feature/p0.3-migrate PR #3 → base: feature/p0.2-docker-stack
                └── ...
```

---

## Agent workflow (each task)

### 1. Find the base branch

| Situation | Branch from |
|---|---|
| First task in stack (or previous merged) | `main` |
| Previous task PR still open | `feature/p{prev}-*` (see [`progress.md`](progress.md) stack table) |
| Previous task merged | `main` (rebase your branch if needed) |

```bash
git fetch origin
git checkout main && git pull origin main
# OR:
git checkout feature/p0.1-scaffold && git pull origin feature/p0.1-scaffold
```

### 2. Create task branch

```bash
git checkout -b feature/p0.2-docker-stack
```

### 3. Implement one BUILD_ORDER task

See [`BUILD_ORDER.md`](BUILD_ORDER.md) deliverables + verify.

### 4. Push and open **draft** PR

**First in stack (base = main):**
```bash
git push -u origin feature/p0.1-scaffold
gh pr create --base main --draft \
  --title "P0.1: Repo scaffold" \
  --body-file .github/pull_request_template.md
```

**Stacked on previous branch:**
```bash
git push -u origin feature/p0.2-docker-stack
gh pr create --base feature/p0.1-scaffold --draft \
  --title "P0.2: Docker Compose dev stack" \
  --body-file .github/pull_request_template.md
```

### 5. Update tracking

In **`progress.md`**:
- Check off task when work is complete
- Add row to **Open stack** table (branch, PR base, PR URL, draft/open)

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
| Merge your own PRs | Leave draft/open for maintainer |
| Skip updating stack table in `progress.md` | Always record branch + PR link |
