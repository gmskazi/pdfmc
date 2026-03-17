# AGENTS.md

## Scope And Stack
- Repository type: Go CLI app (`github.com/gmskazi/pdfmc`).
- Entry point: `main.go` -> `cmd.Execute()`.
- Main packages: `cmd/`, `cmd/program/`, `cmd/pdf/`, `cmd/utils/`, `cmd/ui/*`.
- CLI framework: Cobra.
- TUI components: Bubble Tea + Bubbles + Lip Gloss.
- PDF engine: `pdfcpu`.
- Go version: `1.24.x` (see `go.mod` and CI workflow).

## Rules Files Check
- Cursor rules: not found (`.cursor/rules/` and `.cursorrules` absent).
- Copilot rules: not found (`.github/copilot-instructions.md` absent).
- If these files are added later, treat them as high-priority repo instructions.

## Build Commands
- Build binary: `go build .`
- Run app directly: `go run .`
- Build via Make: `make build`
- Default make target: `make` (aliases `all` -> `build`).
- Clean binary: `make clean`

## Lint And Security Commands
- Full lint target: `make lint`
- Static analysis only: `staticcheck ./...`
- Security scan only: `gosec ./...`
- Install local tools if missing:
  - `go install honnef.co/go/tools/cmd/staticcheck@latest`
  - `go install github.com/securego/gosec/v2/cmd/gosec@latest`

## Test Commands
- All tests (verbose): `go test ./... -v`
- Via make: `make test`
- Coverage profile + HTML report: `make coverage`
- Package-only tests:
  - `go test ./cmd/... -v`
  - `go test ./cmd/program -v`
  - `go test ./cmd/utils -v`

## Single-Test Commands (Preferred For Agents)
- Run one test function in one package:
  - `go test ./cmd/program -run '^TestNewProgram$' -v`
- Run one test in root command package:
  - `go test ./cmd -run '^TestEncryptCommand$' -v`
- Run one test in a nested package:
  - `go test ./cmd/utils -run '^TestCheckProvidedArgs$' -v`
- Run a subset by regex:
  - `go test ./cmd -run 'Test(Encrypt|Decrypt)Command' -v`
- Re-run without cache during debugging:
  - `go test ./cmd/program -run '^Test_getFlagValue$' -count=1 -v`

## Fast Local Agent Workflow
- 1) Change minimal files.
- 2) Run package-level tests for touched package.
- 3) Run targeted single test(s) for changed logic.
- 4) Run `staticcheck ./...` before finalizing.
- 5) Run `gosec ./...` when touching file I/O, paths, encryption, or command input.
- 6) Run full `go test ./... -v` before handoff.

## Commands Agents Should Avoid By Default
- Avoid `make install` in automation (uses `sudo` and writes `/usr/local/bin`).
- Avoid `make gif` unless explicitly requested (depends on `vhs`, rewrites media assets).
- Avoid destructive cleanup outside repo-scoped files.

## Code Organization Conventions
- Keep CLI command wiring in `cmd/*.go`.
- Keep orchestration flow in `cmd/program/program.go`.
- Keep PDF operations in `cmd/pdf/pdf.go`.
- Keep filesystem/path helpers in `cmd/utils/file_utils.go`.
- Keep TUI models and interactions in `cmd/ui/*` packages.
- Keep visual styling constants in `cmd/styles/styles.go`.

## Imports And Formatting
- Always use `gofmt` formatting (tabs, canonical spacing).
- Keep import groups in standard Go style:
  - stdlib first
  - blank line
  - third-party/internal imports
- Prefer short, explicit import aliases only when needed (example: `textInputs`).
- Do not reorder imports in a non-gofmt way.

## Naming Conventions
- Exported identifiers: `PascalCase` (`NewProgram`, `ExecuteMerge`).
- Unexported identifiers: `camelCase` (`getPassword`, `pdfExtension`).
- Receiver names are short and meaningful (`p *Program`, `f *FileUtils`).
- Constants for command keywords are centralized (`merge`, `encrypt`, `decrypt`).
- Test names follow `TestXxx` and are descriptive.

## Types And Struct Patterns
- Prefer small structs with clear ownership (`Program`, `FileUtils`, `PDFProcessor`).
- Keep constructor functions (`NewX`) returning initialized pointers.
- Keep flags/state grouped when logical (example: `MergeFlags`).
- Use explicit fields over `interface{}` where possible.

## Error Handling Conventions
- Return errors upward; avoid swallowing errors.
- In command handlers, render user-facing errors via style helpers:
  - `cmd.PrintErrln(styles.ErrorStyle.Render(err.Error()))`
- Use early returns to reduce nesting.
- Wrap errors with `%w` when adding context.
- Keep human-readable error messages direct and actionable.

## CLI And UX Conventions
- For command args, use `cobra.ArbitraryArgs` where current commands do.
- Keep command flags defined in command `init()`.
- Keep autocomplete integration via `ValidArgsFunction` in command setup.
- Interactive flow is controlled via `FileUtils.CheckProvidedArgs()` + TUI helpers.
- Non-interactive path should work from CLI args alone when feasible.

## File And Path Handling
- Use `filepath` helpers (`Join`, `Ext`) for cross-platform paths.
- Validate whether path is file vs directory before processing.
- Filter PDFs case-insensitively via lowercase suffix checks.
- Avoid hardcoded separators.

## Testing Guidelines
- Existing tests are mostly table-driven; follow that pattern.
- Use `testify/assert` for assertions.
- Prefer `t.TempDir()` for filesystem tests.
- Many tests call `os.Chdir`; avoid `t.Parallel()` in tests that mutate cwd.
- For CLI tests, set both `rootCmd.SetOut` and `rootCmd.SetErr` to capture output.
- Validate both command output text and filesystem side effects.

## Security And Validation Notes
- Encryption/decryption uses `pdfcpu` AES configuration; preserve current API behavior.
- Validate output PDFs when modifying PDF pipeline logic.
- Keep security scan (`gosec`) clean when touching filesystem or crypto code.

## Change Scope Discipline For Agents
- Prefer minimal, package-local changes.
- Do not refactor unrelated files during focused bug fixes.
- Keep public behavior and CLI flags stable unless requested.
- Update tests with behavior changes in the same PR.

## Pre-Handoff Checklist
- `gofmt`-clean code.
- Relevant package tests passing.
- Single targeted tests for changed paths passing.
- `staticcheck ./...` passes.
- `go test ./... -v` passes.
- `gosec ./...` checked for security-sensitive changes.

## Useful Paths
- `main.go`
- `cmd/root.go`
- `cmd/merge.go`
- `cmd/encrypt.go`
- `cmd/decrypt.go`
- `cmd/program/program.go`
- `cmd/pdf/pdf.go`
- `cmd/utils/file_utils.go`
- `.github/workflows/ci.yml`
- `Makefile`
