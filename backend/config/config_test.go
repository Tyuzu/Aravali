package config

import "testing"

func TestInitConfig_AllowsExplicitLocalPostgresValuesInProduction(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
	t.Setenv("REDIS_URL", "localhost:6379")
	t.Setenv("TLS_CERT_PATH", "./cert.pem")
	t.Setenv("TLS_KEY_PATH", "./key.pem")
	t.Setenv("TERMINATE_TLS_AT_LB", "false")
	t.Setenv("POSTGRES_URL", "postgres://apeman:ningning@localhost:5432/eventdb")
	t.Setenv("POSTGRES_USER", "apeman")
	t.Setenv("POSTGRES_PASSWORD", "ningning")
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DB", "eventdb")

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() returned unexpected error: %v", err)
	}

	if cfg.DatabaseURL == "" {
		t.Fatal("DatabaseURL should not be empty")
	}
	if cfg.Env != "production" {
		t.Fatalf("Env = %q, want production", cfg.Env)
	}
}
