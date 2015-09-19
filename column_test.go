package cgoroonga

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSetStringAndGetString(t *testing.T) {
	tempDir, ctx, db := setupTestDB(t, "goroonga-TestSetTimeAndGetTime-")
	defer tearDownTestDB(t, tempDir, ctx, db)

	shortTextType := ctx.At(DB_SHORT_TEXT)
	defer shortTextType.unlink()
	textType := ctx.At(DB_TEXT)
	defer textType.unlink()

	table, err := db.CreateTable("Table1", "",
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, shortTextType)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}

	column, err := table.CreateColumn("column1", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, textType)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}

	count, err := table.RecordCount()
	if err != nil {
		t.Errorf("failed to get a record count: %s", err)
	}
	if count != 0 {
		t.Errorf("record count mismatch: want %s, got %s", 0, count)
	}

	recordID, added, err := table.AddRecord("foo")
	if err != nil {
		t.Errorf("failed to add a record with error: %s", err)
	}
	if !added {
		t.Errorf("should be a new record")
	}

	count, err = table.RecordCount()
	if err != nil {
		t.Errorf("failed to get a record count: %s", err)
	}
	if count != 1 {
		t.Errorf("record count mismatch: want %s, got %s", 1, count)
	}

	value := "bar"
	err = column.SetString(recordID, value)
	if err != nil {
		t.Errorf("failed to set a value to record with error: %s", err)
	}
	actualValue := column.GetString(recordID)
	if actualValue != value {
		t.Errorf("record value mismatch: want %s, got %s", value,
			actualValue)
	}
}

func TestSetTimeAndGetTime(t *testing.T) {
	tempDir, ctx, db := setupTestDB(t, "goroonga-TestSetTimeAndGetTime-")
	defer tearDownTestDB(t, tempDir, ctx, db)

	shortTextType := ctx.At(DB_SHORT_TEXT)
	defer shortTextType.unlink()
	timeType := ctx.At(DB_TIME)
	defer timeType.unlink()

	table, err := db.CreateTable("Table1", "",
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, shortTextType)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}

	column, err := table.CreateColumn("column1", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, timeType)
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

	value := time.Unix(123456789, 987654321)
	err = column.SetTime(recordID, value)
	if err != nil {
		t.Errorf("failed to set a value to record with error: %s", err)
	}
	want := value.UnixNano() / 1000
	got := column.GetTime(recordID).UnixNano() / 1000
	if got != want {
		t.Errorf("record value mismatch: want %s, got %s", want, got)
	}
}

func setupTestDB(t *testing.T, tempDirPrefix string) (tempDir string, ctx *Context, db *DB) {
	tempDir, err := ioutil.TempDir("", tempDirPrefix)
	if err != nil {
		t.Errorf("failed to create a temporary directory with error: %s", err)
	}

	err = Init()
	if err != nil {
		t.Errorf("failed to initialize with error: %s", err)
	}

	ctx, err = NewContext()
	if err != nil {
		t.Errorf("failed to create context with error: %s", err)
	}

	dbPath := filepath.Join(tempDir, "test.db")
	db, err = ctx.CreateDB(dbPath)
	if err != nil {
		t.Errorf("failed to create a database with error: %s", err)
	}

	return
}

func tearDownTestDB(t *testing.T, tempDir string, ctx *Context, db *DB) {
	db.Remove()
	ctx.Close()
	Terminate()
	os.Remove(tempDir)
}
