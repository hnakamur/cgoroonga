package cgoroonga

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	column.Close()

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

type table1 struct {
	key       string
	content   string
	updatedAt time.Time
}

func addTable1Record(t *Table, record table1) error {
	recordID, _, err := t.AddRecord(record.key)
	if err != nil {
		return err
	}
	err = recordID.SetString(t.Column("content"), record.content)
	if err != nil {
		return err
	}
	return recordID.SetTime(t.Column("updated_at"), record.updatedAt)
}

func mustParseRFC3339Time(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSelect(t *testing.T) {
	tempDir, ctx, db := setupTestDB(t, "goroonga-TestSelect-")
	defer tearDownTestDB(t, tempDir, ctx, db)

	table, err := db.CreateTable("Table1", "",
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, DB_SHORT_TEXT)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}

	_, err = table.CreateColumn("content", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, DB_TEXT)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}

	_, err = table.CreateColumn("updated_at", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, DB_TIME)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}

	err = addTable1Record(table, table1{
		key: "key1", content: "content1",
		updatedAt: mustParseRFC3339Time("2015-05-24T12:34:56+09:00"),
	})
	if err != nil {
		t.Errorf("failed to add a record with error: %s", err)
	}
	err = addTable1Record(table, table1{
		key: "key2", content: "content2",
		updatedAt: mustParseRFC3339Time("2015-05-23T10:30:50+09:00"),
	})
	if err != nil {
		t.Errorf("failed to add a record with error: %s", err)
	}

	expr, err := table.CreateQuery("")
	if err != nil {
		t.Errorf("failed to create an expression with error: %s", err)
	}
	err = expr.Parse("_key:@key1", nil, OP_MATCH, OP_AND,
		EXPR_SYNTAX_QUERY|EXPR_ALLOW_PRAGMA|EXPR_ALLOW_COLUMN)
	if err != nil {
		t.Errorf("failed to parse the expression with error: %s", err)
	}

	records, err := table.Select(expr, nil, OP_OR)
	if err != nil {
		t.Errorf("failed to select the table with error: %s", err)
	}

	count, err := records.RecordCount()
	if err != nil {
		t.Errorf("failed to get a record count: %s", err)
	}
	if count != 1 {
		t.Errorf("record count mismatch: want %s, got %s", 1, count)
	}
}
