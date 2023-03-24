package migrations

import (
	"context"
	"fmt"

	"ticketbeastar/pkg/models"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print("users [up migration] ")
		_, err := db.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print("users [down migration] ")
		_, err := db.NewDropTable().Model((*models.User)(nil)).IfExists().Exec(ctx)
		return err
	})
}
