package config

import "time"

type Config struct {
	Server       ServerConfig
	ControlPlane ControlPlaneConfig
	Victoria     VictoriaConfig
	K8s          K8sConfig
	Postgres     PostgresConfig
	Redis        RedisConfig
	Admin        AdminConfig
}

type ServerConfig struct {
	ListenAddr        string
	TLSCertFile       string
	TLSKeyFile        string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

type ControlPlaneConfig struct {
	BaseURL string
}

type VictoriaConfig struct {
	QueryBaseURL string
}

type K8sConfig struct {
	KubeconfigEncryptionKey string
}

type PostgresConfig struct {
	DBURL string
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

type AdminConfig struct {
	BootstrapTelegramBotToken string
	BootstrapTelegramChatID   string
}

func Load() *Config {
	requirePairedEnv("ADMIN_UI_TLS_CERT_FILE", "ADMIN_UI_TLS_KEY_FILE")

	return &Config{
		Server: ServerConfig{
			ListenAddr:        getEnv("ADMIN_UI_LISTEN_ADDR", ":8082"),
			TLSCertFile:       getEnv("ADMIN_UI_TLS_CERT_FILE", ""),
			TLSKeyFile:        getEnv("ADMIN_UI_TLS_KEY_FILE", ""),
			ReadHeaderTimeout: 5 * time.Second,
			ShutdownTimeout:   10 * time.Second,
		},
		ControlPlane: ControlPlaneConfig{
			BaseURL: getEnv("ADMIN_UI_CONTROLPLANE_URL", "http://127.0.0.1:8000"),
		},
		Victoria: VictoriaConfig{
			QueryBaseURL: getEnv("ADMIN_UI_VICTORIA_QUERY_BASE_URL", ""),
		},
		K8s: K8sConfig{
			KubeconfigEncryptionKey: getRequiredEnv("ADMIN_UI_KUBECONFIG_ENCRYPTION_KEY"),
		},
		Postgres: PostgresConfig{
			DBURL: getRequiredEnv("ADMIN_UI_PSQL_DB_URL"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("ADMIN_UI_REDIS_ADDR", "127.0.0.1:6379"),
			Username: getEnv("ADMIN_UI_REDIS_USERNAME", ""),
			Password: getEnv("ADMIN_UI_REDIS_PASSWORD", ""),
			DB:       getEnvInt("ADMIN_UI_REDIS_DB", 0),
		},
		Admin: AdminConfig{
			BootstrapTelegramBotToken: getEnv("ADMIN_UI_BOOTSTRAP_TELEGRAM_BOT_TOKEN", ""),
			BootstrapTelegramChatID:   getEnv("ADMIN_UI_BOOTSTRAP_TELEGRAM_CHAT_ID", ""),
		},
	}
}
