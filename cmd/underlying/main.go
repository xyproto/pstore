package main

import (
	"fmt"
	"github.com/xyproto/pstore"
	"github.com/xyproto/simplehstore"
)

func main() {
	perm, err := pstore.New()
	if err != nil {
		fmt.Println("Could not open database")
		return
	}
	ustate := perm.UserState()

	// A bit of checking is needed, since the database backend is interchangeable
	if pustate, ok := ustate.(*pstore.UserState); ok {
		if host, ok := pustate.Host().(*simplehstore.Host); ok {
			db := host.Database()
			fmt.Printf("PostgreSQL database: %v (%T)\n", db, db)
		}
	} else {
		fmt.Println("Not using the PostgreSQL database backend")
	}
}
