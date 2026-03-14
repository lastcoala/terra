package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func main() {
	var project, module string

	cmd := &cli.Command{
		Name:  "terra-codegen",
		Usage: "Scaffold a new terra application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "project",
				Usage:       "the name of the project",
				Required:    true,
				Destination: &project,
			},
			&cli.StringFlag{
				Name:        "module",
				Usage:       "the name of the go module",
				Required:    true,
				Destination: &module,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return scaffold(project, module)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// scaffold creates the full directory/file structure for a new terra project.
func scaffold(project, module string) error {
	files := buildFiles(project, module)

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		fmt.Printf("created %s\n", path)
	}

	fmt.Println("\nDone! Your new terra project is ready.")
	fmt.Printf("  cd %s\n", project)
	fmt.Println("  go mod tidy")
	return nil
}

// writeFile creates all parent directories and writes data to path.
func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// r replaces the terra module path and project name placeholders inside a template.
// "terra" pkg-level imports stay as-is; only the module path changes.
func r(s, module, project string) string {
	s = strings.ReplaceAll(s, "github.com/lastcoala/terra", module)
	s = strings.ReplaceAll(s, "terra-local", project+"-local")
	s = strings.ReplaceAll(s, "terra-test", project+"-test")
	s = strings.ReplaceAll(s, "TERRA", strings.ToUpper(strings.ReplaceAll(project, "-", "_")))
	s = strings.ReplaceAll(s, "\"terra\"", "\""+project+"\"")
	s = strings.ReplaceAll(s, "database terra", "database "+project)
	return s
}

// buildFiles returns map[filepath]content for every file in the scaffold.
func buildFiles(project, module string) map[string]string {
	root := project + "/"

	files := map[string]string{
		root + "cmd/app/main.go":                        r(tmplCmdAppMain, module, project),
		root + "config/config.go":                       r(tmplConfigGo, module, project),
		root + "config/config.yaml":                     tmplConfigYaml,
		root + "deploy/local/docker-compose.yaml":       r(tmplDeployLocal, module, project),
		root + "deploy/test/docker-compose.yaml":        r(tmplDeployTest, module, project),
		root + "internal/app/domain/.gitkeep":           "",
		root + "internal/app/repo/base_model.go":        tmplRepoBaseModel,
		root + "internal/app/repo/gorm.go":              r(tmplRepoGorm, module, project),
		root + "internal/app/repo/repo.go":              r(tmplRepoRepo, module, project),
		root + "internal/app/rest/v1/helper_test.go":    r(tmplRestV1HelperTest, module, project),
		root + "internal/app/rest/v1/route.go":          r(tmplRestV1Route, module, project),
		root + "internal/app/rest/v1/util.go":           r(tmplRestV1Util, module, project),
		root + "internal/app/rest/rest.go":              r(tmplRestGo, module, project),
		root + "internal/app/service/service.go":        r(tmplServiceGo, module, project),
		root + "internal/app/app.go":                    r(tmplAppGo, module, project),
		root + "internal/mocks/.gitkeep":                "",
		root + "migration/000001_set_timezone.up.sql":   r(tmplMigrationUp, module, project),
		root + "migration/000001_set_timezone.down.sql": tmplMigrationDown,
		root + ".mockery.yaml":                          r(tmplMockeryYaml, module, project),
		root + "go.mod":                                 r(tmplGoMod, module, project),
		root + "go.sum":                                 "",
		root + "Makefile":                               tmplMakefile,
		root + "README.md":                              r(tmplReadme, module, project),
	}
	return files
}

// ---------------------------------------------------------------------------
// File templates
// ---------------------------------------------------------------------------

const tmplCmdAppMain = `package main

import (
	"context"
	"log"
	"os"

	"github.com/lastcoala/terra/config"
	"github.com/lastcoala/terra/internal/app"
	"github.com/lastcoala/terra/pkg/handler"
	"github.com/urfave/cli/v3"
)

func main() {
	var configPath string
	var rest bool

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "rest",
				Value:       false,
				Usage:       "run rest server",
				Destination: &rest,
			}, &cli.StringFlag{
				Name:        "config",
				Value:       "config/config.yaml",
				Usage:       "config file path",
				Destination: &configPath,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg, err := config.LoadConfig(configPath, "TERRA")
			if err != nil {
				return err
			}

			app := app.NewApp(cfg)
			registry := handler.NewRegistry()

			if rest {
				restHandler := app.CreateRestServer()
				registry.Register("REST", restHandler)
			}

			registry.StartAll()
			registry.StopAll()

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
`

const tmplConfigGo = `package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Rest RestConfig ` + "`koanf:\"rest\"`" + `
}

type RestConfig struct {
	Server ServerConfig   ` + "`koanf:\"server\"`" + `
	Db     DatabaseConfig ` + "`koanf:\"db\"`" + `
}

type ServerConfig struct {
	Host string ` + "`koanf:\"host\"`" + `
}

type DatabaseConfig struct {
	DataStore  string ` + "`koanf:\"datastore\"`" + `
	NumberConn int    ` + "`koanf:\"nconn\"`" + `
}

// LoadConfig loads configuration from a YAML file and merges it with environment
// variables that share the given prefix. Environment variables take priority over
// file values when both define the same key.
//
// Env vars are mapped by stripping the prefix, lowercasing, and replacing "_"
// with ".". For example, with prefix "APP", APP_REST_SERVER_HOST maps to
// rest.server.host.
func LoadConfig(configPath, envPrefix string) (Config, error) {
	k := koanf.New(".")

	// 1. Load the YAML file (lower priority).
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("loading config file %q: %w", configPath, err)
	}

	// 2. Overlay environment variables (higher priority).
	prefix := strings.ToUpper(envPrefix) + "_"
	err := k.Load(env.Provider(prefix, ".", func(s string) string {
		// Strip the prefix, lowercase, replace "_" with "." to form the koanf key.
		s = strings.TrimPrefix(s, prefix)
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil)
	if err != nil {
		return Config{}, fmt.Errorf("loading env vars with prefix %q: %w", envPrefix, err)
	}

	// 3. Unmarshal into Config struct.
	cfg := Config{}
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return Config{}, fmt.Errorf("unmarshaling config: %w", err)
	}

	return cfg, nil
}
`

const tmplConfigYaml = `rest:
  server:
    host: ":8080"
  db:
    datastore: "postgres://terra:terra123@localhost:5433/terra?sslmode=disable"
    nConn: 1
`

const tmplDeployLocal = `name: terra-local

services:
  postgres:
    restart: no
    image: postgres:16
    environment:
      - POSTGRES_USER=terra
      - POSTGRES_PASSWORD=terra123
      - POSTGRES_DB=terra
    ports:
      - 5432:5432
`

const tmplDeployTest = `name: terra-test

services:
  postgres:
    restart: no
    image: postgres:16
    environment:
      - POSTGRES_USER=terra
      - POSTGRES_PASSWORD=terra123
      - POSTGRES_DB=terra
    ports:
      - 5433:5432
`

const tmplRepoBaseModel = `package repo

import "time"

type BaseModel struct {
	Id        int ` + "`gorm:\"primaryKey;autoIncrement\"`" + `
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
	IsDeleted bool
}
`

const tmplRepoGorm = `package repo

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormRepo(path string, nConn int) (*gorm.DB, error) {
	var err error
	client, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: client}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
	if err != nil {
		return nil, err
	}

	dbX, err := db.DB()
	if err != nil {
		return nil, err
	}
	dbX.SetMaxIdleConns(nConn)
	dbX.SetMaxOpenConns(nConn)
	dbX.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}
`

const tmplRepoRepo = `package repo

import (
	"context"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/pkg/filter"
)

type IUserRepo interface {
	GetUser(ctx context.Context, id int) (domain.User, error)
	GetUsers(ctx context.Context, offset, limit int, filters ...filter.Filter) ([]domain.User, error)
	InsertUser(ctx context.Context, user domain.User, createdBy string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User, updatedBy string) (domain.User, error)
	DeleteUser(ctx context.Context, id int, deletedBy string) error
}
`

const tmplRestV1HelperTest = `package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func requestJsonTestHelper[T any](method string, data T, queryParam string) (echo.Context,
	*httptest.ResponseRecorder) {

	e := echo.New()
	var req *http.Request
	url := fmt.Sprintf("/%v", queryParam)

	if any(data) != nil {
		jsonData, _ := json.Marshal(data)
		req = httptest.NewRequest(method, url, bytes.NewReader(jsonData))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
`

const tmplRestV1Route = `package v1

import (
	echo "github.com/labstack/echo/v4"
	"github.com/lastcoala/terra/internal/app/service"
)

const (
	USER_ID = "userId"
)

func Route(group *echo.Group, userService service.IUserService) {

	userRoute := NewUserRoute(userService)
	group.POST("/user", userRoute.InsertUser)
	group.GET("/user", userRoute.GetUsers)
	group.GET("/user/:"+USER_ID, userRoute.GetUser)
	group.PUT("/user/:"+USER_ID, userRoute.UpdateUser)
	group.DELETE("/user/:"+USER_ID, userRoute.DeleteUser)
}
`

const tmplRestV1Util = `package v1

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func queryParamToOffsetLimit(c echo.Context, useDefault bool) (int, int) {
	var page, limit int

	err := echo.QueryParamsBinder(c).Int("page", &page).
		Int("limit", &limit).BindError()

	if err != nil {
		return 0, 10
	}

	if page < 1 {
		if useDefault {
			return 0, 10
		} else {
			return 0, 0
		}
	}

	return (page - 1) * limit, limit
}

func getUserId(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param(USER_ID))
	if err != nil {
		return 0, err
	}
	return id, nil
}
`

const tmplRestGo = `package rest

import (
	"github.com/labstack/echo/v4"
	v1 "github.com/lastcoala/terra/internal/app/rest/v1"
	"github.com/lastcoala/terra/internal/app/service"
	"github.com/lastcoala/terra/pkg/handler"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title TERRA API
// @version 1.0
// @description This is a sample REST API Server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /v1
type RestHandler struct {
	*handler.ServerHandler
	userService service.IUserService
}

func NewRest(host string, userService service.IUserService) *RestHandler {

	server := handler.NewServerHandler(host)

	return &RestHandler{
		ServerHandler: server,
		userService:   userService,
	}
}

func (h *RestHandler) Start() error {
	h.Echo = echo.New()
	h.Echo.HideBanner = true
	h.Route()
	h.StartServer()

	return nil
}

func (h *RestHandler) Route() {
	h.ServerHandler.Route()

	v1Route := h.Echo.Group("/v1")

	v1Route.GET("/swagger/*", echoSwagger.WrapHandler)
	v1Route.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	v1.Route(v1Route, h.userService)
}
`

const tmplServiceGo = `package service

import (
	"context"

	"github.com/lastcoala/terra/internal/app/domain"
)

type IUserService interface {
	InsertUser(ctx context.Context, user domain.User, createdBy string) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	GetUsers(ctx context.Context, offset, limit int) ([]domain.User, error)
	UpdateUser(ctx context.Context, user domain.User, updatedBy string) (domain.User, error)
	ChangePassword(ctx context.Context, id int, newPassword string, updatedBy string) (domain.User, error)
	DeleteUser(ctx context.Context, id int, deletedBy string) error
}
`

const tmplAppGo = `package app

import (
	"github.com/lastcoala/terra/config"
	"github.com/lastcoala/terra/internal/app/repo"
	"github.com/lastcoala/terra/internal/app/rest"
	"github.com/lastcoala/terra/internal/app/service"
	"github.com/lastcoala/terra/pkg/handler"
)

type App struct {
	cfg config.Config
}

func NewApp(cfg config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) CreateRestServer() handler.IHandler {
	gorm, err := repo.NewGormRepo(a.cfg.Rest.Db.DataStore, a.cfg.Rest.Db.NumberConn)
	if err != nil {
		panic(err)
	}

	userRepo := repo.NewUserGormRepo(gorm)
	userService := service.NewUserService(userRepo)

	restHandler := rest.NewRest(a.cfg.Rest.Server.Host, userService)

	return restHandler
}
`

const tmplMigrationUp = `alter database terra set timezone to 'Asia/Jakarta';
`

const tmplMigrationDown = `alter database terra set timezone to 'UTC'; 
`

const tmplMockeryYaml = `with-expecter: true
packages:
  github.com/lastcoala/terra/internal/app/repo:
    config:
      all: true
      dir: internal/mocks
      outpkg: mocks
  github.com/lastcoala/terra/internal/app/service:
    config:
      all: true
      dir: internal/mocks
      outpkg: mocks
`

const tmplGoMod = `module github.com/lastcoala/terra

go 1.25.0

require (
	github.com/knadh/koanf/parsers/yaml v1.1.0
	github.com/knadh/koanf/providers/env v1.1.0
	github.com/knadh/koanf/providers/file v1.2.1
	github.com/knadh/koanf/v2 v2.3.3
	github.com/labstack/echo/v4 v4.15.1
	github.com/lib/pq v1.11.2
	github.com/stretchr/testify v1.11.1
	github.com/swaggo/echo-swagger v1.4.1
	github.com/urfave/cli/v3 v3.7.0
	golang.org/x/crypto v0.48.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)
`

const tmplMakefile = `compose-up-local:
	docker-compose -f deploy/local/docker-compose.yaml up -d

compose-down-local:
	docker-compose -f deploy/local/docker-compose.yaml down

compose-up-test:
	docker-compose -f deploy/test/docker-compose.yaml up -d

compose-down-test:
	docker-compose -f deploy/test/docker-compose.yaml down

gen-mocks:
	touch internal/mocks/a.txt && rm internal/mocks/* && mockery --dir internal --output internal/mocks --all
`

const tmplReadme = `# terra

A terra application.
`
