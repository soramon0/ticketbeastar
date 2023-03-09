package migrations

import (
	"ticketbeastar/pkg/utils"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

func init() {
	utils.Must(Migrations.DiscoverCaller())
}
