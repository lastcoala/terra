package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeYAML creates a temporary YAML config file and returns its path.
// The file is automatically removed when the test ends.
func writeYAML(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

const baseYAML = `
rest:
  server:
    host: ":8080"
  db:
    datastore: "postgres://terra:terra123@localhost:5432/terra?sslmode=disable"
    nConn: 1
`

func TestLoadConfig(t *testing.T) {
	t.Run("loads values from yaml file", func(t *testing.T) {
		path := writeYAML(t, baseYAML)

		cfg, err := LoadConfig(path, "TEST")

		require.NoError(t, err)
		assert.Equal(t, ":8080", cfg.Rest.Server.Host)
		assert.Equal(t, "postgres://terra:terra123@localhost:5432/terra?sslmode=disable", cfg.Rest.Db.DataStore)
		assert.Equal(t, 1, cfg.Rest.Db.NumberConn)
	})

	t.Run("env vars override yaml values", func(t *testing.T) {
		path := writeYAML(t, baseYAML)

		t.Setenv("APP_REST_SERVER_HOST", ":9090")
		t.Setenv("APP_REST_DB_NCONN", "5")

		cfg, err := LoadConfig(path, "APP")

		require.NoError(t, err)
		assert.Equal(t, ":9090", cfg.Rest.Server.Host)
		assert.Equal(t, 5, cfg.Rest.Db.NumberConn)
		// datastore not overridden — should still come from file
		assert.Equal(t, "postgres://terra:terra123@localhost:5432/terra?sslmode=disable", cfg.Rest.Db.DataStore)
	})

	t.Run("env vars with different prefix are ignored", func(t *testing.T) {
		path := writeYAML(t, baseYAML)

		t.Setenv("OTHER_REST_SERVER_HOST", ":7777")

		cfg, err := LoadConfig(path, "APP")

		require.NoError(t, err)
		// The OTHER_ prefix should not affect values loaded under APP_ prefix
		assert.Equal(t, ":8080", cfg.Rest.Server.Host)
	})

	t.Run("env vars only — no yaml key conflict", func(t *testing.T) {
		path := writeYAML(t, baseYAML)

		t.Setenv("APP_REST_DB_DATASTORE", "postgres://new-host:5432/newdb?sslmode=disable")

		cfg, err := LoadConfig(path, "APP")

		require.NoError(t, err)
		assert.Equal(t, "postgres://new-host:5432/newdb?sslmode=disable", cfg.Rest.Db.DataStore)
	})

	t.Run("returns error for missing config file", func(t *testing.T) {
		missingPath := filepath.Join(t.TempDir(), "nonexistent.yaml")

		_, err := LoadConfig(missingPath, "APP")

		assert.ErrorContains(t, err, "loading config file")
	})

	t.Run("empty env prefix loads yaml only", func(t *testing.T) {
		path := writeYAML(t, baseYAML)

		// With prefix "", the constructed prefix becomes "_"; effectively no
		// meaningful env vars will match, so all values come from the file.
		cfg, err := LoadConfig(path, "")

		require.NoError(t, err)
		assert.Equal(t, ":8080", cfg.Rest.Server.Host)
		assert.Equal(t, 1, cfg.Rest.Db.NumberConn)
	})
}
