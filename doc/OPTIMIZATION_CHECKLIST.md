# Optimization Checklist

## Phase 1 — Critical correctness
- [x] Fix `encryptCmd.ValidArgsFunction` wiring in `cmd/encrypt.go`
- [x] Add/adjust command autocomplete test coverage for encrypt command
- [x] Verify: `go test ./cmd -run 'Test.*Encrypt.*' -v`

## Phase 2 — Deterministic behavior
- [x] Make `GetSelectedPDFs()` return stable order (cursor/index-sorted)
- [x] Add test asserting deterministic selected PDF order
- [x] Verify: `go test ./cmd/ui/multiSelect -v`

## Phase 3 — Duplication cleanup
- [x] Consolidate password prompt path(s) in `cmd/program/program.go`
- [x] Consolidate shared Lip Gloss styles via `cmd/styles/styles.go`
- [x] Keep existing UX behavior unchanged (interactive/non-interactive)

## Phase 4 — Test utilities
- [x] Extract shared PDF/file test helpers to a common test helper location
- [x] Refactor existing tests to consume shared helpers
- [x] Verify: `go test ./... -v`

## Phase 5 — Optional hardening
- [ ] Evaluate reducing mutable state in `FileUtils` for cleaner APIs
- [ ] Add/update docs/comments only where behavior is non-obvious

## Validation gates
- [x] `gofmt` clean
- [x] `staticcheck ./...`
- [x] `gosec ./...` (security-sensitive paths)
- [x] `go test ./... -v`
