package infra_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/internal/infra"
	"luny.dev/cherryauctions/pkg/closer"
	"luny.dev/cherryauctions/pkg/env"
)

func TestDatabaseCanConnect(t *testing.T) {
	dbUrl := env.Fatalenv("DATABASE_URL")
	db := infra.SetupDatabase(dbUrl)
	assert.NotNil(t, db)

	sqlDb, err := db.DB()
	assert.Nil(t, err)
	defer closer.CloseResources(sqlDb)
}
