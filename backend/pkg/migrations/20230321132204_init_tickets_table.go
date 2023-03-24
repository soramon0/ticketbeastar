package migrations

import (
	"context"
	"fmt"

	"ticketbeastar/pkg/models"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print("tickets [up migration] ")
		_, err := db.NewCreateTable().Model((*models.Ticket)(nil)).IfNotExists().Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print("tickets [down migration] ")
		_, err := db.NewDropTable().Model((*models.Ticket)(nil)).IfExists().Exec(ctx)
		return err
	})
}
