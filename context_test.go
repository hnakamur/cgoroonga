package cgoroonga

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewContextAndClose(t *testing.T) {
	err := Init()
	if err != nil {
		t.Errorf("failed to initialize with error: %s", err)
	}
	defer func() {
		err := Terminate()
		if err != nil {
			t.Errorf("failed to initialize with error: %s", err)
		}
	}()

	ctx, err := NewContext()
	if err != nil {
		t.Errorf("failed to create context with error: %s", err)
	}
	defer func() {
		err := ctx.Close()
		if err != nil {
			t.Errorf("failed to close context with error: %s", err)
		}
	}()
}

func TestCreateDBAndRemove(t *testing.T) {
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

	dirName, err := ioutil.TempDir("", "goroonga-TestCreateDBAndRemove-")
	if err != nil {
		t.Errorf("failed to create a temporary directory with error: %s", err)
	}
	defer func() {
		err := os.Remove(dirName)
		if err != nil {
			t.Errorf("failed to remove the temporary directory with error: %s", err)
		}
	}()

	path := filepath.Join(dirName, "test.db")
	db, err := ctx.CreateDB(path)
	if err != nil {
		t.Errorf("failed to create a database with error: %s", err)
	}
	defer func() {
		err := db.Remove()
		if err != nil {
			t.Errorf("failed to remove the database with error: %s", err)
		}
		if fileExists(path) {
			t.Errorf("database file should be not exist")
		}
	}()

	if !fileExists(path) {
		t.Errorf("database file should exist")
	}
}

func TestOpenDBAndClose(t *testing.T) {
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

	dirName, err := ioutil.TempDir("", "goroonga-TestOpenDBAndClose-")
	if err != nil {
		t.Errorf("failed to create a temporary directory with error: %s", err)
	}
	defer os.Remove(dirName)

	path := filepath.Join(dirName, "test.db")
	db, err := ctx.CreateDB(path)
	if err != nil {
		t.Errorf("failed to create a database with error: %s", err)
	}

	db.Close()

	db, err = ctx.OpenDB(path)
	if err != nil {
		t.Errorf("failed to open the database with error: %s", err)
	}
	defer db.Remove()
}

func TestOpenNonExistentDB(t *testing.T) {
	dirName, err := ioutil.TempDir("", "goroonga-TestOpenNonExistentDB-")
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

	path := filepath.Join(dirName, "test.db")
	_, err = ctx.OpenDB(path)
	if err != InvalidArgumentError {
		t.Errorf("unexpected err from OpenDB, want: %s, got: %s", InvalidArgumentError, err)
	}
}

func fileExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}
