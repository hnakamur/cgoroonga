package cgoroonga

import "testing"

func TestSnippet(t *testing.T) {
	tempDir, ctx, db := setupTestDB(t, "goroonga-TestSnippet-")
	defer tearDownTestDB(t, tempDir, ctx, db)

	shortTextType := ctx.At(DB_SHORT_TEXT)
	defer shortTextType.unlink()
	textType := ctx.At(DB_TEXT)
	defer textType.unlink()
	timeType := ctx.At(DB_TIME)
	defer timeType.unlink()

	table, err := db.CreateTable("Table1", "",
		OBJ_TABLE_HASH_KEY|OBJ_PERSISTENT, shortTextType)
	if err != nil {
		t.Errorf("failed to create a table with error: %s", err)
	}

	_, err = table.CreateColumn("content", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, textType)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}

	_, err = table.CreateColumn("updated_at", "",
		OBJ_PERSISTENT|OBJ_COLUMN_SCALAR, timeType)
	if err != nil {
		t.Errorf("failed to create a column with error: %s", err)
	}

	wants := []table1{
		table1{
			key: "key1",
			content: `The MIT License (MIT)

Copyright (c) <year> <copyright holders>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.`,
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
	err = expr.Parse("content:@permit", nil, OP_MATCH, OP_AND,
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

	flags := SNIP_COPY_TAG | SNIP_SKIP_LEADING_SPACES
	textMaxLen := 200
	snippet := expr.Snippet(flags, textMaxLen, 1, true,
		[][]string{
			[]string{"<b>", "</b>"},
		})
	i := 0
	for {
		id, hasNext := cursor.Next()
		if !hasNext {
			break
		}

		record := getTable1Record(table, id)
		snipResults, err := snippet.Exec(record.content)
		if err != nil {
			t.Errorf("failed to execute snippet: %s", err)
		}
		if len(snipResults) != 1 {
			t.Errorf("snip result count mismatch: want %d; got %d",
				1, len(snipResults))
		}
		want := `copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to <b>permit</b> persons to whom the Software is
furnished to do so, subject to the following conditions:

The ab`
		if snipResults[0] != want {
			t.Errorf("snip result mismatch: want %d; got %d",
				want, snipResults[0])
		}

		i++

	}
}
