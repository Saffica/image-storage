package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

type db struct {
	pool *pgxpool.Pool
}

const schemaTableName = "schema_version"

func New(dbUser, dbPassword, dbHost, dbName, pathToMigrations string, dbPort uint) (*db, error) {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	cfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.TODO(), cfg)
	if err != nil {
		return nil, err
	}

	db := &db{
		pool: pool,
	}

	err = db.migrateDatabase(pathToMigrations)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *db) Close() {
	r.pool.Close()
}

func (r *db) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *db) Exec(ctx context.Context, query string, arg ...interface{}) (pgconn.CommandTag, error) {
	return r.pool.Exec(ctx, query, arg...)
}

func (r *db) Query(ctx context.Context, query string, arg ...interface{}) (pgx.Rows, error) {
	return r.pool.Query(ctx, query, arg...)
}

func (r *db) QueryRow(ctx context.Context, query string, arg ...interface{}) pgx.Row {
	return r.pool.QueryRow(ctx, query, arg...)
}

func (r *db) migrateDatabase(pathToMigrations string) error {
	ctx := context.TODO()
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), schemaTableName)
	if err != nil {
		return err
	}

	version, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		return err
	}

	err = migrator.LoadMigrations(pathToMigrations)
	if err != nil {
		return err
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	newVersion, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Migration done. From version: %d to version %d\n", version, newVersion)
	return nil
}
