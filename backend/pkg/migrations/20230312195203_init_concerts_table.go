package migrations

import (
	"context"
	"fmt"

	"ticketbeastar/pkg/models"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print("concerts [up migration] ")
		_, err := db.NewCreateTable().Model((*models.Concert)(nil)).IfNotExists().Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print("concerts [down migration] ")
		_, err := db.NewDropTable().Model((*models.Concert)(nil)).IfExists().Exec(ctx)
		return err
	})
}
