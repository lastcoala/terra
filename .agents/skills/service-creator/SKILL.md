---
name: service-creator
description: Generates service implementations, mock generation, and tests for new domain entities.
---

# Service Creator Skill

This skill provides instructions on how to create the service layer for a new entity in this project, following the conventions established by the `UserService`.

## Prerequisites

- The domain entity must already exist in `internal/app/domain/` (use the `domain-creator` skill).
- The repository interface `I[Entity]Repo` must already exist in `internal/app/repo/repo.go` (use the `repo-creator` skill).

## Core Conventions

### 1. Service Implementation (`internal/app/service/[entity].go`)

- Create an `[Entity]Service` struct with a single unexported field `repo I[Entity]Repo`.
- Provide a constructor: `New[Entity]Service(repo repo.I[Entity]Repo) *[Entity]Service`.
- Add a private `error` method to wrap and log errors:
  ```go
  func (s [Entity]Service) error(ctx context.Context, err error, method string, params ...any) error {
      errF := fmt.Errorf("[Entity]Service.(%v)(%v) %w", method, params, err)
      slog.ErrorContext(ctx, errF.Error())
      return errF
  }
  ```
- Implement the standard CRUD methods below. Every method must call `s.error(ctx, ...)` on any failure and return a zero-value struct (not `nil`) on error:
  - `Insert[Entity](ctx context.Context, entity domain.[Entity], createdBy string) (domain.[Entity], error)` — call `entity.Validate(...)` first; if valid, call `s.repo.Insert[Entity]`.
  - `Get[Entity](ctx context.Context, id int) (domain.[Entity], error)`
  - `Get[Entities](ctx context.Context, offset, limit int) ([]domain.[Entity], error)` — return `nil` (not empty slice) on repo error.
  - `Update[Entity](ctx context.Context, entity domain.[Entity], updatedBy string) (domain.[Entity], error)` — fetch the existing record first with `Get[Entity]` to preserve any fields that should not be overwritten by the caller (e.g. `Password`); validate before calling the repo.
  - `Delete[Entity](ctx context.Context, id int, deletedBy string) error`
- Import only `"context"`, `"fmt"`, `"log/slog"`, the domain package (`github.com/lastcoala/terra/internal/app/domain`), and the repo package (`github.com/lastcoala/terra/internal/app/repo`).

### 2. Mock Generation (`internal/mocks/mock_I[Entity]Repo.go`)

The project uses `mockery` (configured in `.mockery.yaml`) to auto-generate mocks. After the repo interface is committed, run:

```sh
go generate ./...
```

or call mockery directly:

```sh
mockery
```

The generated mock lands in `internal/mocks/mock_I[Entity]Repo.go` as `MockI[Entity]Repo`. **Never edit generated mock files by hand.**

### 3. Service Unit Tests (`internal/app/service/[entity]_test.go`)

- Use `package service` (same package, no `_test` suffix).
- Declare `var testCtx = ctx.NewCtx("1234567890")` as a package-level test context from `github.com/lastcoala/terra/pkg/ctx`.
- Create a helper `new[Entity]Service(t *testing.T) (*[Entity]Service, *mocks.Mock[I[Entity]Repo])` that wires a fresh mock into the service:
  ```go
  func new[Entity]Service(t *testing.T) (*[Entity]Service, *mocks.MockI[Entity]Repo) {
      t.Helper()
      mockRepo := mocks.NewMockI[Entity]Repo(t)
      svc := New[Entity]Service(mockRepo)
      return svc, mockRepo
  }
  ```
- Create a `valid[Entity]() domain.[Entity]` helper that returns a fully-valid entity for reuse across tests.
- Use `mockRepo.EXPECT().[Method](...)` for setting expectations. Prefer:
  - `mock.Anything` when the exact argument value is irrelevant.
  - `mock.MatchedBy(func(...) bool { ... })` when the argument is derived and cannot be known in advance (e.g. a bcrypt hash).
- Each service method must be tested with at minimum:
  - `"success"` — happy path with `.Return(expected, nil)`.
  - A subtest for each validation or pre-condition failure that must **not** reach the repo (no `EXPECT()` set; `AssertExpectations` from testify will catch any unexpected call).
  - `"repo_error_is_propagated"` — mock returns an error; assert the service surfaces it.
- Use section separator comments between method groups:
  ```go
  // ---------------------------------------------------------------------------
  // MethodName
  // ---------------------------------------------------------------------------
  ```
- Use `require.NoError` for setup steps that must pass, and `assert.Error` / `assert.NoError` for the assertion under test.

## Example Usage

Refer to the `examples/` directory in this skill folder for a concrete service implementation and its tests.
