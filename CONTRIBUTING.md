# Contributing to Pinchtab

## Setup

```bash
git clone https://github.com/pinchtab/pinchtab.git
cd pinchtab
./dev doctor
```

Requires **Go 1.25+** and **Google Chrome**.

## Development Workflow

1. Make your changes
2. Run `./dev check`
3. Run `./dev test`
4. Commit — pre-commit hook runs checks automatically
5. Push: `git pull --rebase && git push`

### Creating a Pull Request

**Important:** When creating a PR, please keep the **"Allow edits from maintainers"** checkbox **enabled** (it's on by default). This lets us:

- Apply small fixes directly
- Resolve merge conflicts automatically
- Rebase and update your branch without asking

This significantly speeds up the merge process. Thank you! 🙏

## Checks and Tests

```bash
./dev doctor        # setup environment and hooks
./dev check         # format, vet, build, lint
./dev format dashboard
./dev test          # unit + E2E tests
./dev test unit     # unit tests only
./dev e2e           # E2E tests (curl + CLI)
./dev e2e curl      # curl tests only
./dev e2e cli       # CLI tests only
```

If you want the raw commands instead of `dev`:

```bash
go test ./... -count=1 -v
./dev e2e
```

## Style

- Adhere to **SOLID** principles, specifically using interfaces for dependency inversion.
- Handle all error returns explicitly.
- Lowercase error strings, wrap with `%w`.
- Tests should live in the same package as the source code.
- No new dependencies without significant technical justification.
- Comments when necessary

## Hooks

`./dev doctor` offers to install the git hooks. They enforce formatting and checks before commit.
You can also install them directly with:

```bash
./scripts/install-hooks.sh
```
