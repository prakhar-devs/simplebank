package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/prakhar-devs/simplebank/util"
)

var TestQueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {
	// we used "../.." because the app.env file is present two folders above the file where we are loading config
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	TestDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	TestQueries = New(TestDB)

	os.Exit(m.Run())
}
