package app

import (
	adminui "aurora-adminui"
	"aurora-adminui/infra/psql"
	redisinfra "aurora-adminui/infra/redis"
	"aurora-adminui/internal/config"
	"context"
	"crypto/tls"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Application struct {
	cfg         *config.Config
	server      *http.Server
	db          *pgxpool.Pool
	redisClient *redis.Client
}

// New wires infrastructure, runs bootstrap migrations, and builds the admin UI HTTP server.
func New(ctx context.Context, cfg *config.Config) (*Application, error) {
	controlPlaneURL, err := url.Parse(cfg.ControlPlane.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid ADMIN_UI_CONTROLPLANE_URL: %w", err)
	}

	db, err := psql.NewPool(ctx, cfg.Postgres.DBURL)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	redisClient, err := redisinfra.NewClient(ctx, redisinfra.Config{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	if err := runMigrations(ctx, db); err != nil {
		db.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("run admin migrations: %w", err)
	}

	distFS, err := fs.Sub(adminui.EmbeddedDist, "dist")
	if err != nil {
		db.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("embed dist fs: %w", err)
	}

	modules, err := NewModules(db, redisClient, cfg)
	if err != nil {
		db.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("initialize modules: %w", err)
	}
	if err := modules.AdminAuthHandler.BootstrapIfNeeded(ctx); err != nil {
		db.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("bootstrap admin token: %w", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	RegisterRoutes(router, modules, controlPlaneURL, distFS)

	server := &http.Server{
		Addr:              cfg.Server.ListenAddr,
		Handler:           router,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}
	if cfg.Server.TLSCertFile != "" && cfg.Server.TLSKeyFile != "" {
		server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	return &Application{
		cfg:         cfg,
		server:      server,
		db:          db,
		redisClient: redisClient,
	}, nil
}

// Run starts the HTTP server and blocks until the process receives a shutdown signal.
func (a *Application) Run() error {
	servingHTTPS := a.cfg.Server.TLSCertFile != "" && a.cfg.Server.TLSKeyFile != ""
	if servingHTTPS {
		log.Printf("admin ui listening with https on %s", a.cfg.Server.ListenAddr)
	} else {
		log.Printf("admin ui listening with http on %s", a.cfg.Server.ListenAddr)
	}

	errCh := make(chan error, 1)
	go func() {
		if servingHTTPS {
			errCh <- a.server.ListenAndServeTLS(a.cfg.Server.TLSCertFile, a.cfg.Server.TLSKeyFile)
			return
		}
		errCh <- a.server.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("shutting down admin ui after %s", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("admin ui server failed: %w", err)
		}
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		_ = a.server.Close()
	}
	_ = a.redisClient.Close()
	a.db.Close()
	return nil
}

// runMigrations replays the embedded SQL bootstrap files in lexical order.
func runMigrations(ctx context.Context, db *pgxpool.Pool) error {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read admin migrations: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".up.sql") {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	for _, name := range names {
		sqlBytes, err := migrationFiles.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read admin migration %s: %w", name, err)
		}
		if _, err := db.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("exec admin migration %s: %w", name, err)
		}
	}
	return nil
}
