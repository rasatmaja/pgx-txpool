package pgxtxpool

import "testing"

func TestConfig(t *testing.T) {
	t.Run("should not panic and return config", func(t *testing.T) {

		defer func() {
			if r := recover(); r != nil {
				t.Log(r)
				t.FailNow()
			}
		}()

		cfg := config{}
		options := []Option{
			SetHost("localhost", "5432"),
			SetCredential("postgres", "postgres"),
			SetDatabase("postgres"),
			WithSSLMode("disable"),
			WithMaxConns(10),
			WithMaxIdleConns("30s"),
			WithMaxConnLifetime("5m"),
		}

		for _, opt := range options {
			opt(&cfg)
		}

		pgxCfg := cfg.ParseToPGXConfig()
		if pgxCfg == nil {
			t.Log("config is nil")
			t.FailNow()
		}
	})

	t.Run("should panic when parse to pgxpool config", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.FailNow()
			}
		}()

		cfg := config{}
		cfg.dsn.RawQuery = "sslmode=XXXX"
		cfg.ParseToPGXConfig()
	})
}
