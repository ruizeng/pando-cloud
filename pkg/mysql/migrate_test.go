package mysql

import (
	"testing"
)

func TestMigrate(t *testing.T) {
	err := MigrateDatabase("localhost", "3306", "PandoCloud", "root", "")
	if err != nil {
		t.Error(err)
	}
}
