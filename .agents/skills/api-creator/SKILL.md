---
name: api-creator
description: Generates REST API route handlers, DTOs, route registration, and tests for new domain entities in the v1 REST layer.
---

# API Creator Skill

This skill provides instructions on how to create the REST API layer for a new entity in this project, following the conventions established by the `User` entity in `internal/app/rest/v1/`.

Each entity requires four files:

| File | Purpose |
|---|---|
| `internal/app/rest/v1/[entity]_dto.go` | Request / response DTOs and mapping helpers |
| `internal/app/rest/v1/[entity]_route.go` | Route handler struct with CRUD methods |
| `internal/app/rest/v1/route.go` | Register the new entity's routes (modify existing) |
| `internal/app/rest/v1/[entity]_test.go` | Unit tests for every route handler |

---

## Core Conventions

### 1. DTOs (`internal/app/rest/v1/[entity]_dto.go`)

- Use `package v1`.
- Define a **response DTO** named `[Entity]Resp` with exported fields and `json` tags matching the domain struct fields (excluding sensitive fields like `Password`).
- Define request DTOs per operation:
  - `Insert[Entity]Req` – for POST, includes all writable fields.
  - `Update[Entity]Req` – for PUT, includes updatable fields only.
  - Additional request types as needed (e.g., `Change[Entity]PasswordReq`).
- Each request DTO that creates a domain object must implement a `ToDomain() domain.[Entity]` method.
- Provide private mapping helpers `[entity]ToResp(e domain.[Entity]) [Entity]Resp` and `[entity]sToResp(es []domain.[Entity]) [][Entity]Resp`.

### 2. Route Handler (`internal/app/rest/v1/[entity]_route.go`)

- Use `package v1`.
- Define a `[Entity]Route` struct that holds an `[entity]Service service.I[Entity]Service` field.
- Provide a `New[Entity]Route(svc service.I[Entity]Service) *[Entity]Route` constructor.
- Implement the standard five CRUD methods on `*[Entity]Route`:
  - `Insert[Entity](c echo.Context) error`
  - `Get[Entity]s(c echo.Context) error`
  - `Get[Entity](c echo.Context) error`
  - `Update[Entity](c echo.Context) error`
  - `Delete[Entity](c echo.Context) error`
- Annotate every method with a full Swagger `godoc` comment block (`// @Summary`, `// @Tags`, `// @Security`, `// @Accept`, `// @Produce`, `// @Param`, `// @Success`, `// @Router`).
- Use `c.Bind(&req)` with `http.StatusBadRequest` on failure for POST/PUT handlers.
- Use `getUserId(c)` (from `util.go`) to parse the path param `userId` and return `http.StatusBadRequest` on failure.
- Use `queryParamToOffsetLimit(c, true)` (from `util.go`) for list endpoints; pass `true` to use defaults when params are absent.
- Use `NewResponseDto(msg, data, "[entity]")` for single-entity responses.
- Use `NewResponsesDto(msg, data, "[entity]s")` for list responses.
- Use `MSG_SUCCESS` constant for successful response messages.
- Return `http.StatusInternalServerError` on service errors and `http.StatusOK` on success.
- The hardcoded `createdBy`/`updatedBy`/`deletedBy` actor should be `"admin"` until authentication is wired in.

### 3. Route Registration (`internal/app/rest/v1/route.go`)

- `route.go` holds the `Route(group *echo.Group, ...)` function and the path parameter constants (e.g., `USER_ID = "userId"`).
- Add a new parameter to `Route` for the new entity's service interface.
- Add a constant for the new entity's path param ID (e.g., `PRODUCT_ID = "productId"`).
- Instantiate `New[Entity]Route(svc)` and register the five CRUD routes:
  ```go
  group.POST("/[entity]", entityRoute.Insert[Entity])
  group.GET("/[entity]", entityRoute.Get[Entity]s)
  group.GET("/[entity]/:"+ENTITY_ID, entityRoute.Get[Entity])
  group.PUT("/[entity]/:"+ENTITY_ID, entityRoute.Update[Entity])
  group.DELETE("/[entity]/:"+ENTITY_ID, entityRoute.Delete[Entity])
  ```
- Also update `internal/app/rest/rest.go` to pass the new service from `RestHandler` into `v1.Route(...)`.

### 4. Unit Tests (`internal/app/rest/v1/[entity]_test.go`)

- Use `package v1` (same package, white-box testing).
- Import: `"errors"`, `"net/http"`, `"testing"`, `"github.com/lastcoala/terra/internal/app/domain"`, `"github.com/lastcoala/terra/internal/mocks"`, `"github.com/stretchr/testify/assert"`.
- Use `requestJsonTestHelper[T]` (from `helper_test.go`) to create the Echo context and recorder.
- Use `mocks.NewMock[I][Entity]Service(t)` and `.EXPECT()` for mocking the service layer.
- Separate test functions per handler using section separator comments:
  ```go
  // ---------------------------------------------------------------------------
  // Insert[Entity]
  // ---------------------------------------------------------------------------
  ```
- Cover the following cases for each handler:

  | Handler | Cases |
  |---|---|
  | `Insert[Entity]` | success → `200`, service error → `500` |
  | `Get[Entity]s` | success with params → `200`, no params defaults → `200`, service error → `500` |
  | `Get[Entity]` | success → `200`, invalid id → `400`, service error → `500` |
  | `Update[Entity]` | success → `200`, invalid id → `400`, service error → `500` |
  | `Delete[Entity]` | success → `200`, invalid id → `400`, service error → `500` |

- For list tests, set query params both in the URL **and** on the context:
  ```go
  c, rec := requestJsonTestHelper("GET", struct{}{}, "/entity?page=1&limit=10")
  c.QueryParams().Set("page", "1")
  c.QueryParams().Set("limit", "10")
  ```
- Use `assert.JSONEq` to validate the exact JSON response body where possible.
- `mockUserService.EXPECT()` calls without a `.Return()` are not needed for error-path tests that return before hitting the service.

---

## Shared Utilities (do not recreate)

These helpers already exist in the package and must be reused:

| Symbol | File | Purpose |
|---|---|---|
| `getUserId(c)` | `util.go` | Parses `:userId` (or any `:*Id`) path param as int |
| `queryParamToOffsetLimit(c, useDefault)` | `util.go` | Converts `?page=&limit=` to `offset, limit` ints |
| `NewResponseDto(msg, data, key)` | `response_dto.go` | Single-item JSON response wrapper |
| `NewResponsesDto[T](msg, data, key)` | `response_dto.go` | Slice JSON response wrapper |
| `MSG_SUCCESS` | `response_dto.go` | `"success"` constant |
| `requestJsonTestHelper[T](...)` | `helper_test.go` | Creates an Echo context for tests |

> **Note:** `util.go` currently has `getUserId` hardcoded to the `USER_ID` param constant. For a different entity with a different path param name, either add a new helper (e.g., `getProductId`) or make the helper generic by accepting the param name as an argument.

---

## Example Usage

Refer to the `examples/` directory in this skill folder for a concrete implementation of the full REST layer for the `User` entity.
