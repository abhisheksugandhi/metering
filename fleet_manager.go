package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var dbObject *sql.DB

func init_fleet_manager() {
	dbConfig = get_db_config()
	var err error
	dbObject, err = sql.Open(dbConfig.Adapter, fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbConfig.UserName, dbConfig.Database))
	if err != nil {
		panic(err)
	}
}

func get_fleet_mapping() map[string]int {
	mapping := make(map[string]int)
	rows, err := dbObject.Query("select xid, verifone_fleet_code from fleets")
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var verifone_fleet_code string
		var fleet_xid int
		rows.Scan(&fleet_xid, &verifone_fleet_code)
		if len(verifone_fleet_code) > 0 {
			mapping[verifone_fleet_code] = fleet_xid
		}
	}
	return mapping
}
