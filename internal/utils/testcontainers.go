package utils

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/mwojtyna/swift-api/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const pgImage = "postgres:17"
const migrationsFolder = "migrations"

type TestWithPostgresArgs struct {
	Pc  *postgres.PostgresContainer
	Env *config.Env
	Ctx context.Context
}

func TestWithPostgres(f func(TestWithPostgresArgs)) {
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalln("Failed to load env")
	}

	ctx := context.Background()

	// Find all 'up' migrations
	pattern := filepath.Join(env.ProjectRootPath, migrationsFolder, "*up*.sql")
	migrations, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalln("Failed to search for migrations")
	}

	pc, err := postgres.Run(ctx,
		pgImage,
		postgres.WithInitScripts(migrations...),
		postgres.WithDatabase(env.DB_NAME),
		postgres.WithUsername(env.DB_USER),
		postgres.WithPassword(env.DB_PASS),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(pc); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	port, err := pc.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalln("failed to get port")
	}
	env.DB_PORT = port.Port()

	f(TestWithPostgresArgs{
		Pc:  pc,
		Env: &env,
		Ctx: ctx,
	})
}
