package config

import "testing"

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want required DATABASE_URL error")
	}
}

func TestLoadReadsConfiguration(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://reader:reader@localhost/reader")
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("GRPC_ADDR", "127.0.0.1:9090")
	t.Setenv("SHUTDOWN_TIMEOUT", "3s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.GRPCAddr != "127.0.0.1:9090" {
		t.Errorf("GRPCAddr = %q", cfg.GRPCAddr)
	}
	if cfg.ShutdownTimeout.String() != "3s" {
		t.Errorf("ShutdownTimeout = %s", cfg.ShutdownTimeout)
	}
}
