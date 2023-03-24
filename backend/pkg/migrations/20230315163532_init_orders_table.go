package migrations

import (
	"context"
	"fmt"

	"ticketbeastar/pkg/models"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print("orders [up migration] ")
		_, err := db.NewCreateTable().Model((*models.Order)(nil)).IfNotExists().Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print("orders [down migration] ")
		_, err := db.NewDropTable().Model((*models.Order)(nil)).IfExists().Exec(ctx)
		return err
	})
}
