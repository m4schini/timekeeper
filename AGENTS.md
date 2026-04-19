# AGENTS.md

This file provides guidance to AI coding agents (Claude and others) working in this repository.

---

## Project Overview

**Type:** Web application (frontend / fullstack)  
**Primary language / stack:** Go

---

## Development Workflow

### Branching

We use **feature branches + pull requests** following the [Conventional Branch 1.0.0](https://conventional-branch.github.io/) specification. Never commit directly to `main`.

**Format:**
```
<type>/<description>
```

**Allowed types:**

| Type | Use for |
|------|---------|
| `feature/` or `feat/` | New features |
| `bugfix/` or `fix/` | Bug fixes |
| `hotfix/` | Urgent / production fixes |
| `release/` | Release preparation |
| `chore/` | Non-code tasks (deps, docs, tooling) |

**Rules:**
- Use only **lowercase letters, numbers, and hyphens** (`a-z`, `0-9`, `-`). No uppercase, underscores, or spaces.
- Dots (`.`) are permitted only in `release/` descriptions for version numbers (e.g. `release/v1.2.0`).
- No consecutive, leading, or trailing hyphens or dots.
- Keep descriptions concise and purpose-driven.
- If a ticket number exists, include it: `feat/issue-123-add-pagination`.
- Keep branches short-lived and focused on a single concern.
- Open a PR to merge into `main`; all CI checks must pass before merging.

**Valid examples:**
```
feat/add-pagination
fix/nil-pointer-in-db-pool
hotfix/security-patch
release/v1.2.0
chore/update-dependencies
feat/issue-123-add-pagination
```

**Invalid examples:**
```
Feature/Add-Login        ❌  uppercase
feature/new--login       ❌  consecutive hyphens
feature/-new-login       ❌  leading hyphen
fix/header_bug           ❌  underscore
unknown/some-task        ❌  unknown type
```

### CI/CD
- **GitHub Actions** runs on every push to `main`. All checks must pass.
- Do **not** modify any files under `.github/` or any CI/CD configuration.

---

## Testing

**Framework:** `go test` (stdlib)  
**Policy:** TDD — write tests first, then implementation.

- Before writing any implementation code, write a failing test that defines the expected behavior.
- Run tests frequently during development:
  ```bash
  go test ./...
  ```
- Run tests with the race detector before finalizing any change:
  ```bash
  go test -race ./...
  ```
- Test files live alongside the code they test (`foo.go` → `foo_test.go`).
- Use table-driven tests wherever multiple input/output cases exist.
- Do **not** delete or skip existing tests. If a test appears wrong, leave a comment — do not remove it.
- Aim for clear test names that describe behavior: `TestUserService_CreateUser_ReturnsErrorOnDuplicate`.

---

## Code Style & Conventions

- Follow standard Go idioms and the [Effective Go](https://go.dev/doc/effective_go) guidelines.
- Review all changed code based on [Go Review Comments](https://go.dev/wiki/CodeReviewComments)
- Run `gofmt` (or `goimports`) before committing — all code must be properly formatted.
- Run `go vet ./...` and ensure there are no warnings.
- Keep interfaces small; prefer defining them at the point of use (consumer side), not the producer side.
- Return errors explicitly — never panic in library or handler code.
- Define explicit errors types that implement the error interface (Avoid fmt.Errorf and errors.New)
- Wrap errors with context using `fmt.Errorf("doing X: %w", err)`.
- Avoid global state; pass dependencies explicitly via constructors or function parameters.
- All exported symbols must have a doc comment.

---

## Licensing

- All code must be compatible with **AGPL-3.0-or-later**.
- Every new source file must include the appropriate SPDX license identifier at the top:
  ```go
  // SPDX-License-Identifier: AGPL-3.0-or-later
  ```

---

## Commits & Attribution

### Commit Message Format

All commit messages must follow the [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**

| Type | SemVer impact | Use for |
|------|--------------|---------|
| `feat` | MINOR | A new feature |
| `fix` | PATCH | A bug fix |
| `docs` | — | Documentation only |
| `style` | — | Formatting, whitespace — no logic change |
| `refactor` | — | Neither a fix nor a feature |
| `test` | — | Adding or updating tests |
| `chore` | — | Build process, tooling, dependency updates |
| `perf` | — | Performance improvements |
| `ci` | — | CI/CD configuration changes |
| `revert` | — | Reverting a previous commit |
| `build` | — | Build system changes |

**Rules:**
- The type and description are **required**; all other parts are optional.
- A scope may follow the type in parentheses to add context: `feat(auth):`.
- The description must immediately follow `<type>[scope]: ` — lowercase, imperative mood, no trailing period.
- The body, if present, must begin one blank line after the description and may span multiple paragraphs.
- Footers must begin one blank line after the body and use the format `Token: value` or `Token #value`. Tokens use `-` instead of spaces (except `BREAKING CHANGE`).
- Breaking changes must be indicated by appending `!` after the type/scope, or by including a `BREAKING CHANGE: <description>` footer — or both.
- If `!` is used, the `BREAKING CHANGE:` footer may be omitted; the description serves as the breaking change summary.
- `BREAKING CHANGE` correlates to a **MAJOR** SemVer bump and may appear on any commit type.
- When reverting, use the `revert` type and reference the SHA(s) being reverted in a `Refs:` footer.

**Examples:**

```
feat(handler): add pagination to list endpoints

fix(db): handle nil pointer when connection pool is exhausted

feat!: remove deprecated v1 API routes

BREAKING CHANGE: all /v1/* routes have been removed, use /v2/*

docs: update README with local dev setup instructions

revert: let us never again speak of the noodle incident

Refs: 676104e, a215868
```

---

### Signed-off-by
- Agents **must not** add `Signed-off-by` tags. Only humans can legally certify the Developer Certificate of Origin (DCO).
- The human submitter is responsible for reviewing all AI-generated code, verifying licensing compliance, and adding their own `Signed-off-by` before submitting.

### Assisted-by tag
When an AI agent contributes to a commit, include an `Assisted-by` tag in the commit message using this format:

```
Assisted-by: AGENT_NAME:MODEL_VERSION [TOOL1] [TOOL2]
```

- `AGENT_NAME` — name of the AI tool or framework (e.g. `Claude`)
- `MODEL_VERSION` — specific model version used (e.g. `claude-sonnet-4-6`)
- `[TOOL1] [TOOL2]` — optional specialized analysis tools used (e.g. `staticcheck`, `golangci-lint`)
- Basic tools like `git`, `gcc`, `make`, and editors should **not** be listed

Example:

```
Assisted-by: Claude:claude-sonnet-4-6 staticcheck golangci-lint
```

---

## What Agents Should and Should Not Do

### ✅ Agents may
- Read any source file to understand context.
- Create new `.go` source and `_test.go` files.
- Refactor code as long as all tests continue to pass.
- Add or update doc comments.
- Suggest improvements via comments prefixed with `// AGENT NOTE:`.

### ❌ Agents must not
- Modify **any** config or infrastructure files — this includes `.github/`, `Dockerfile`, `docker-compose*.yml`, `Makefile`, `go.mod`, `go.sum`, and any `*.yaml`/`*.toml` config files. These are **strictly off-limits**.
- Delete or skip existing tests.
- Add new dependencies to `go.mod` without explicitly flagging it for human review in the commit message.
- Store secrets, credentials, or hardcoded environment-specific values in source files.
- Push code that does not compile or leaves `go test ./...` in a failing state.
- Add `Signed-off-by` tags — this is reserved for human contributors certifying the DCO.

---

## Running the Project Locally

```bash
# Build
go build ./...

# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Format code
gofmt -w .

# Vet
go vet ./...
```

---

## Questions or Ambiguity

If a task is underspecified or involves a non-trivial architectural decision, **stop and ask** rather than guessing. Surface uncertainty early — a short clarifying question is always better than code that needs to be fully rewritten.