package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestShouldReturnCodeToXidMapping(t *testing.T) {
	db, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{"xid", "verifone_fleet_code"}
	sqlmock.ExpectQuery("select xid, verifone_fleet_code from fleets").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,vfone1"))

	set_db_object(db)
	mapping := get_fleet_mapping()

	if len(mapping) != 1 {
		t.Errorf("it was not expected to return map of length %d", len(mapping))
	}

	fleet_xid, ok := mapping["vfone1"]
	if !ok || fleet_xid != 1 {
		t.Errorf("It was expected to return id as %d for code %s. But some thing broke", 1, "vfone1")
	}

	if err = db.Close(); err != nil {
		t.Errorf("Error '%s' was not expected while closing the database", err)
	}
}
