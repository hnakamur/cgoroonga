package cgoroonga

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSetStringAndGetString(t *testing.T) {
	dirName, err := ioutil.TempDir("", "goroonga-TestSetStringAndGetString-")
	if err != nil {
		t.Errorf("failed to create a temporary directory with error: %s", err)
	}
	defer os.Remove(dirName)

	err = Init()
	if err != nil {
		t.Errorf("failed to initialize with error: %s", err)
	}
	defer Terminate()

	ctx, err := NewContext()
	if err != nil {
		t.Errorf("failed to create context with error: %s", err)
	}
	defer ctx.Close()

	dbPath := filepath.Join(dirName, "test.db")
	db, err := ctx.CreateDB(dbPath)
	if err != nil {
		t.Errorf("failed to create a database with error: %s", err)
	}
	defer db.Remove()

	table, err := db.CreateTable("Table1", "",
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}

	column, err := table.CreateColumn("column1", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, DB_TEXT)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}

	recordID, added, err := table.AddRecord("foo")
	if err != nil {
		t.Errorf("failed to add a record with error: %s", err)
	}
	if !added {
		t.Errorf("should be a new record")
	}

	value := "bar"
	err = recordID.SetString(column, value)
	if err != nil {
		t.Errorf("failed to set a value to record with error: %s", err)
	}
	actualValue := recordID.GetString(column)
	if actualValue != value {
		t.Errorf("record value mismatch: want %s, got %s", value,
			actualValue)
	}
}
