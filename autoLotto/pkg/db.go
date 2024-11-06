package pkg

import (
	"autoLotto/ent"
)

var DbClient *ent.Client

func ConnectPostgres() error {
	var err error
	DbClient, err = ent.Open("postgres", "host=localhost port=5432 user=auto_lotto dbname=auto_lotto password=1234 sslmode=disable")
	return err
}
