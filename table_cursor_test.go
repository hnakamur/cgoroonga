package cgoroonga

import "testing"

func TestOpenTableCursor(t *testing.T) {
	tempDir, ctx, db := setupTestDB(t, "goroonga-TestOpenTableCursor-")
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

	wants := []table1{
		table1{
			key: "key1", content: "content1",
			updatedAt: mustParseRFC3339Time("2015-05-24T12:34:56+09:00"),
		},
		table1{
			key: "key2", content: "content2",
			updatedAt: mustParseRFC3339Time("2015-05-23T10:30:50+09:00"),
		},
	}
	err = addTable1Record(table, wants[0])
	if err != nil {
		t.Errorf("failed to add a record with error: %s", err)
	}
	err = addTable1Record(table, wants[1])
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

	cursor, err := records.OpenTableCursor("", "", 0, -1, CURSOR_ASCENDING)
	if err != nil {
		t.Errorf("failed to open a table cursor: %s", err)
	}
	defer func() {
		err := cursor.Close()
		if err != nil {
			t.Errorf("failed to close the table cursor: %s", err)
		}
	}()

	i := 0
	for {
		id, hasNext := cursor.Next()
		if !hasNext {
			break
		}

		want := wants[i]
		got := getTable1Record(table, id)
		if got.key != want.key {
			t.Errorf("key mismatch for i=%s, got: %s; want: %s", want.key, got.key)
		}
		if got.content != want.content {
			t.Errorf("content mismatch for i=%s, got: %s; want: %s", want.content, got.content)
		}
		if got.updatedAt != want.updatedAt {
			t.Errorf("updatedAt mismatch for i=%s, got: %s; want: %s", want.updatedAt, got.updatedAt)
		}

		i++
	}
}
