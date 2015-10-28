package mysql

import (
	"testing"
)

func TestMigrate(t *testing.T) {
	err := MigrateDatabase("localhost", "3306", "", "root", "")
	if err != nil {
		t.Error(err)
	}
}
