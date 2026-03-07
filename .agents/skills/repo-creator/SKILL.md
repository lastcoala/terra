---
name: repo-creator
description: Generates database migrations, GORM models, repository implementations, and tests for new domain entities.
---

# Repo Creator Skill

This skill provides instructions on how to create the repository layer for a new domain entity in the `internal/app/repo` package, along with the necessary database migrations in the `migration` folder.

## Core Conventions

When generating a new repository for an entity (e.g., `Product`), you must follow these rules:

### 1. Database Migrations
- Place migrations in the `migration/` directory using sequential prefixes (e.g., `000003_create_products.up.sql`).
- Every table MUST have a numeric primary key: `id SERIAL PRIMARY KEY`.
- Every table MUST include audit columns with standard constraints:
  - `created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()`
  - `updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()`
  - `created_by VARCHAR(255) NOT NULL`
  - `updated_by VARCHAR(255) NOT NULL`
  - `is_deleted BOOLEAN NOT NULL DEFAULT FALSE`

### 2. GORM Models (`internal/app/repo/[entity]_gorm_model.go`)
- Create an `[Entity]GormModel` struct that embeds `BaseModel` (from `base_model.go`).
- Use appropriate `gorm` tags (e.g., `gorm:"not null"`, `gorm:"uniqueIndex"`).
- Add mapping functions:
  - `New[Entity]GormModel(domain domain.[Entity]) [Entity]GormModel`
  - `(m [Entity]GormModel) ToDomain() domain.[Entity]`
- Add a slice type `type [Entity]GormModels [][Entity]GormModel` with `ToDomains()` method.
- Add `TableName() string` methods to both the model and the slice.

### 3. Repository Interfaces & Implementations (`internal/app/repo/repo.go` & `internal/app/repo/[entity]_gorm.go`)
- Append an `I[Entity]Repo` interface in `repo.go` with standard CRUD methods:
  - `Get[Entity](ctx context.Context, id int) (domain.[Entity], error)`
  - `Get[Entities](ctx context.Context, offset, limit int, filters ...filter.Filter) ([]domain.[Entity], error)`
  - `Insert[Entity](ctx context.Context, entity domain.[Entity], createdBy string) (domain.[Entity], error)`
  - `Update[Entity](ctx context.Context, entity domain.[Entity], updatedBy string) (domain.[Entity], error)`
  - `Delete[Entity](ctx context.Context, id int, deletedBy string) error`
- Create `[Entity]GormRepo` struct wrapping `*gorm.DB` in `[entity]_gorm.go`.
- Implement an `error(ctx context.Context, err error, method string, params ...any) error` function to log errors using `slog.ErrorContext` and return wrapped errors with `fmt.Errorf`.
- Ensure all queries use `r.db.WithContext(ctx)`.
- `Delete[Entity]` MUST be a soft delete: update `is_deleted = true` and `updated_by = deletedBy`. It should not physically delete records.

### 4. Unit Tests (`internal/app/repo/[entity]_gorm_test.go`)
- Tests must use a real PostgreSQL database with a transaction that rolls back at the end of each test scope:
  ```go
  tx := db.Begin()
  defer tx.Rollback()
  repo := New[Entity]GormRepo(tx)
  ```
- Use `ctx.NewCtx("some_id")` from `github.com/lastcoala/terra/pkg/ctx` for contexts.
- Use `github.com/stretchr/testify/assert` and `require` for assertions.

## Example Usage
Refer to the `examples/` directory in this skill folder for concrete implementations of a model, repo, migration, and test.

