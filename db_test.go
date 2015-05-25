package cgoroonga

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTableAndRemove(t *testing.T) {
	err := Init()
	if err != nil {
		t.Errorf("failed to initialize with error: %s", err)
	}
	defer Terminate()

	ctx, err := NewContext()
	if err != nil {
		t.Errorf("failed to create context with error: %s", err)
	}
	defer ctx.Close()

	dirName, err := ioutil.TempDir("", "goroonga-TestCreateTableAndRemove-")
	if err != nil {
		t.Errorf("failed to create a temporary directory with error: %s", err)
	}
	defer os.Remove(dirName)

	dbPath := filepath.Join(dirName, "test.db")
	db, err := ctx.CreateDB(dbPath)
	if err != nil {
		t.Errorf("failed to create a database with error: %s", err)
	}
	defer db.Remove()

	tableName := "Table1"
	tablePath := dbPath + ".Table1"
	table, err := db.CreateTable(tableName, tablePath,
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT,
		DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}
	defer func() {
		err := table.Remove()
		if err != nil {
			t.Errorf("failed to remove the table with error: %s", err)
		}
		if fileExists(tablePath) {
			t.Errorf("table file should be not exist")
		}
	}()

	if !fileExists(tablePath) {
		t.Errorf("table file should exist")
	}

	actualTableName := table.Name()
	if actualTableName != tableName {
		t.Errorf("table name mismatch: want %s, got %s", tableName,
			actualTableName)
	}

	actualTablePath := table.Path()
	if actualTablePath != tablePath {
		t.Errorf("table path mismatch: want %s, got %s", tablePath,
			actualTablePath)
	}
}

func TestOpenTableAndClose(t *testing.T) {
	err := Init()
	if err != nil {
		t.Errorf("failed to initialize with error: %s", err)
	}
	defer Terminate()

	ctx, err := NewContext()
	if err != nil {
		t.Errorf("failed to create context with error: %s", err)
	}
	defer ctx.Close()

	dirName, err := ioutil.TempDir("", "goroonga-TestOpenTableAndClose-")
	if err != nil {
		t.Errorf("failed to create a temporary directory with error: %s", err)
	}
	defer os.Remove(dirName)

	dbPath := filepath.Join(dirName, "test.db")
	db, err := ctx.CreateDB(dbPath)
	if err != nil {
		t.Errorf("failed to create a database with error: %s", err)
	}
	defer db.Remove()

	table, err := db.CreateTable("Table1", "",
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT,
		DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}

	table.Close()

	table, err = db.OpenTable("Table1")
	if err != nil {
		t.Errorf("failed to open the table with error: %s", err)
	}
}

func TestOpenNotExistingTable(t *testing.T) {
	err := Init()
	if err != nil {
		t.Errorf("failed to initialize with error: %s", err)
	}
	defer Terminate()

	ctx, err := NewContext()
	if err != nil {
		t.Errorf("failed to create context with error: %s", err)
	}
	defer ctx.Close()

	dirName, err := ioutil.TempDir("", "goroonga-TestOpenNotExistingTable-")
	if err != nil {
		t.Errorf("failed to create a temporary directory with error: %s", err)
	}
	defer os.Remove(dirName)

	dbPath := filepath.Join(dirName, "test.db")
	db, err := ctx.CreateDB(dbPath)
	if err != nil {
		t.Errorf("failed to create a database with error: %s", err)
	}
	defer db.Remove()

	_, err = db.OpenTable("Table1")
	if err != NotFoundError {
		t.Errorf("unexpected err from OpenTable, want %s, got %s.", NotFoundError, err)
	}
}

func TestOpenOrCreateTable(t *testing.T) {
	dirName, err := ioutil.TempDir("", "goroonga-TestOpenOrCreateTable-")
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

	table, err := db.OpenOrCreateTable("Table1", "",
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT,
		DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}
	if !fileExists(table.Path()) {
		t.Errorf("table file should exist")
	}
}
