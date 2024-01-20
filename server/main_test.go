package main

import (
	"globe-and-citizen/layer8/server/config"
	"testing"
)

func TestPostgresConnection(t *testing.T) {
	config.InitDB()
	t.Log(*config.DB)
}
