package models_test

import (
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"
	"ticketbeastar/tests/concerts"

	"github.com/uptrace/bun"
)

func TestConcertModels(t *testing.T) {
	db := database.OpenConnection(utils.GetTestDatabaseURL(), utils.InitLogger())
	defer database.CloseConnection(db)

	tests := map[string]func(t *testing.T, cs models.ConcertService){
		"with a published_at date are published": func(t *testing.T, cs models.ConcertService) {
			publishedA := concerts.CreateConcert(t, db, nil, "", true)
			publishedB := concerts.CreateConcert(t, db, nil, "", true)
			unpublished := concerts.CreateConcert(t, db, &models.Concert{PublishedAt: bun.NullTime{}}, "", true)

			concerts, err := cs.FindPublished()
			if err != nil {
				t.Fatalf("Failed to retrieve published concerts %v", err)
			}

			if len(*concerts) != 2 {
				t.Fatalf("only two concerts should be retunred; got %d", len(*concerts))
			}

			for _, concert := range *concerts {
				if concert.Id == unpublished.Id {
					t.Fatalf("concerts should not have unpublished concert %v", concert)
				}

				if concert.Id != publishedA.Id && concert.Id != publishedB.Id {
					t.Fatalf("concerts do not contain published concerts %v", concerts)
				}
			}
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			concerts.SetupTable(t, db)
			defer concerts.TeardownTable(t, db)
			test(t, models.NewConcertService(db))
		})
	}
}
