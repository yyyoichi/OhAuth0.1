package database

import (
	"context"
	"log"
	"testing"
)

func TestDb(t *testing.T) {
	ctx := context.Background()
	db, _ := NewDatabase()
	c, _ := db.GetServieClientByID(ctx, "500")
	c.Name = "rename"
	log.Println(db.serviceClientByID["500"].Name)
}
