package database

import (
	"database/sql"
	"log"
	"os"
	"ticketbeastar/pkg/migrations"
	"ticketbeastar/pkg/utils"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func OpenConnection(dsn string, l *log.Logger) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	if utils.GetStageStatus() == "dev" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	cliApp := &cli.App{
		Name: "bun",

		Commands: []*cli.Command{
			migrations.NewMigratorCommand(migrate.NewMigrator(db, migrations.Migrations)),
		},
	}

	if len(os.Args) >= 2 && os.Args[1] == "db" {
		defer CloseConnection(db)
		utils.Must(cliApp.Run(os.Args))
		os.Exit(0)
	}

	return db
}

func CloseConnection(db *bun.DB) {
	utils.Must(db.Close())
}
