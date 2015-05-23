package cgoroonga

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateColumnAndRemove(t *testing.T) {
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

	dirName, err := ioutil.TempDir("", "goroonga-TestCreateColumnAndRemove-")
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
	tablePath := dbPath + "." + tableName
	table, err := db.CreateTable(tableName, tablePath,
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}
	defer table.Remove()

	columnName := "column1"
	columnPath := tablePath + "." + columnName
	column, err := table.CreateColumn(columnName, columnPath,
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, DB_TEXT)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}
	defer func() {
		err := column.Remove()
		if err != nil {
			t.Errorf("failed to remove the column with error: %s", err)
		}
		if fileExists(columnPath) {
			t.Errorf("column file should be not exist")
		}
	}()

	if !fileExists(columnPath) {
		t.Errorf("column file should exist")
	}

	actualColumnName := column.Name()
	if actualColumnName != columnName {
		t.Errorf("column name mismatch: want %s, got %s.", columnName,
			actualColumnName)
	}

	actualColumnPath := column.Path()
	if actualColumnPath != columnPath {
		t.Errorf("column path mismatch: want %s, got %s.", columnPath,
			actualColumnPath)
	}
}

func TestCreateColumnWithDefaultPathAndRemove(t *testing.T) {
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

	dirName, err := ioutil.TempDir("", "goroonga-TestCreateColumnAndRemove-")
	if err != nil {
		t.Errorf("failed to create a temporary directory with error: %s", err)
	}
	defer os.Remove(dirName)

	dbPath := filepath.Join(dirName, "test.db")
	db, err := ctx.CreateDB(dbPath)
	if err != nil {
		t.Errorf("failed to create a database with error: %s", err)
	}
	var tablePath, columnPath string
	defer func() {
		// NOTE: when you remove the database, tables and columns are
		// removed automatically.
		err := db.Remove()
		if err != nil {
			t.Errorf("failed to remove the database with error: %s", err)
		}

		if fileExists(tablePath) {
			t.Errorf("table file should not exist")
		}
		if fileExists(columnPath) {
			t.Errorf("column file should not exist")
		}
	}()

	table, err := db.CreateTable("Table1", "",
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}
	tablePath = table.Path()
	if !fileExists(tablePath) {
		t.Errorf("table file should exist")
	}

	column, err := table.CreateColumn("column1", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, DB_TEXT)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}
	columnPath = column.Path()
	if !fileExists(columnPath) {
		t.Errorf("column file should exist")
	}
}

func TestOpenColumnAndClose(t *testing.T) {
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

	dirName, err := ioutil.TempDir("", "goroonga-TestOpenColumnAndClose-")
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
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}

	column, err := table.CreateColumn("column1", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, DB_TEXT)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}
	err = column.Close()
	if err != nil {
		t.Errorf("failed to close the column with error: %s", err)
	}

	column, err = table.OpenColumn("column1")
	if err != nil {
		t.Errorf("failed to open a column with error: %s", err)
	}
}

func TestOpenNonExistingColumn(t *testing.T) {
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

	dirName, err := ioutil.TempDir("", "goroonga-TestOpenNonExistingColumn-")
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
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}

	_, err = table.OpenColumn("column1")
	if err != NotFoundError {
		t.Errorf("unexpected err from OpenColumn, want %s, got %s.", NotFoundError, err)
	}
}
