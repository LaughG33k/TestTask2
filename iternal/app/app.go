package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/LaughG33k/TestTask2/client/psql"
	"github.com/LaughG33k/TestTask2/iternal"
	"github.com/LaughG33k/TestTask2/iternal/handler"
	"github.com/LaughG33k/TestTask2/iternal/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/stdlib"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/yaml.v2"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
)

func Run() {

	ctx := context.Background()

	cfg := initCfg("./config.yaml")

	cfg.ReadTimeoutSec = int64(time.Duration(cfg.ReadTimeoutSec) * time.Millisecond)
	cfg.WriteTimeoutSec = int64(time.Duration(cfg.WriteTimeoutSec) * time.Millisecond)

	if err := startMigrations("file://migrations/psql/auth/", cfg.UserDbCfg); err != nil {
		log.Panic(err)
	}

	if err := startMigrations("file://migrations/psql/auth/", cfg.NewsDbCfg); err != nil {
		log.Panic(err)
	}

	fiberApp := fiber.New(fiber.Config{
		CaseSensitive: true,
		ReadTimeout:   time.Duration(cfg.ReadTimeoutSec) * time.Second,
		WriteTimeout:  time.Duration(cfg.WriteTimeoutSec) * time.Second,
		IdleTimeout:   1 * time.Minute,
	})

	newsdb, err := initDb(ctx, cfg.NewsDbCfg)

	if err != nil {
		log.Panic(err)
	}

	usersdb, err := initDb(ctx, cfg.UserDbCfg)

	if err != nil {
		log.Panic(err)
	}

	newsRepo := &repository.News{Db: newsdb}
	userRepo := &repository.User{Db: usersdb}
	rtRepo := &repository.RefreshToken{Db: usersdb}

	jwtWorker := iternal.NewJwtWorker("testTask")

	authHandler := handler.Auth{
		Fiber:         fiberApp,
		JwtWorker:     jwtWorker,
		Repo:          userRepo,
		RtRepo:        rtRepo,
		RequstTimeout: time.Duration(cfg.WriteTimeoutSec),
		Ctx:           context.Background(),
	}

	getNewsHandler := handler.GetNews{
		Fiber:         fiberApp,
		Repo:          newsRepo,
		RequstTimeout: time.Duration(cfg.WriteTimeoutSec),
		Ctx:           context.Background(),
	}

	editNewsHandler := handler.EditNews{
		FiberApp:      fiberApp,
		Repo:          newsRepo,
		Ctx:           context.Background(),
		RquestTiemout: time.Duration(cfg.WriteTimeoutSec),
		JwtWorker:     jwtWorker,
	}

	editNewsHandler.Handle()
	getNewsHandler.Handle()
	authHandler.Handle()

	fiberApp.Listen(cfg.Host + ":" + cfg.Port)

}

func initDb(ctx context.Context, cfg iternal.DbConfig) (*reform.DB, error) {
	psqlClient, err := psql.NewPool(ctx, 3, 2*time.Second, cfg)

	sqldb := stdlib.OpenDBFromPool(psqlClient)

	if err != nil {
		return nil, err
	}

	if err := sqldb.Ping(); err != nil {
		return nil, err
	}

	return reform.NewDB(sqldb, postgresql.Dialect, nil), nil
}

func startMigrations(path string, cfg iternal.DbConfig) error {
	migrations, err := migrate.New(
		path,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, "disable"),
	)

	if err != nil {
		return err
	}

	if err := migrations.Up(); err != nil {
		if err.Error() != "no change" {
			return err
		}
	}

	migrations.Close()

	return nil
}

func initCfg(path string) iternal.AppCfg {

	file, err := os.ReadFile(path)

	if err != nil {
		log.Panic(err)
	}

	cfg := iternal.AppCfg{}

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Panic(err)
	}

	return cfg
}
