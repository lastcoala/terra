---
name: domain-creator
description: Generates domain entity structs, validation methods, and tests for new domain entities.
---

# Domain Creator Skill

This skill provides instructions on how to create the domain layer for a new entity in this project, following the conventions established by the `User` entity.

## Core Conventions

When generating a new domain entity (e.g., `Product`), you must follow these rules:

### 1. Domain Entity (`internal/app/domain/[entity].go`)

- Create a `[Entity]` struct with an `Id int` field (JSON tag `"id"`) and any entity-specific fields.
- Add a private `error(err error, method string, params ...any) error` method that wraps errors with context:
  ```go
  func (e Entity) error(err error, method string, params ...any) error {
      return fmt.Errorf("Entity.(%v)(%v) %w", method, params, err)
  }
  ```
- Implement field-level validation methods named `Validate[Field]() error` for each field that has business rules. Each method must:
  - Use the private `error()` helper with the method name.
  - Return `nil` on success and a wrapped error on failure.
- Implement a composite `Validate() error` method that calls all individual validators.
- Entity-specific business logic methods (e.g., `EncryptPassword`, `CheckPassword`) follow the same error-wrapping pattern.

### 2. Domain Unit Tests (`internal/app/domain/[entity]_test.go`)

- Use `package domain` (same package, no `_test` suffix) so unexported types are accessible.
- Test every exported method using table-driven tests where applicable.
- Import only `"testing"`, `"github.com/stretchr/testify/assert"`, and `"github.com/stretchr/testify/require"` plus any domain-specific packages (e.g., `golang.org/x/crypto/bcrypt`).
- Follow the section separator comment style:
  ```go
  // ---------------------------------------------------------------------------
  // MethodName
  // ---------------------------------------------------------------------------
  ```
- Use `assert.Error` / `assert.NoError` for table-driven cases. Use `require.NoError` before operations that must succeed for the rest of the test to be meaningful.

## Example Usage

Refer to the `examples/` directory in this skill folder for a concrete implementation of a domain entity and its tests.
