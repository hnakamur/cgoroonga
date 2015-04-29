package main

import (
	"fmt"

	grn "github.com/hnakamur/cgoroonga"
)

type Site struct {
	URL   string
	Title string
}

func createSiteTable(ctx *grn.Ctx) (err error) {
	keyType := ctx.At(grn.DB_SHORT_TEXT)
	table, err := ctx.TableOpenOrCreate("Site", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, keyType, nil)
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, table)

	columnType := ctx.At(grn.DB_TEXT)
	column, err := ctx.ColumnOpenOrCreate(table, "title", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, columnType)
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, column)

	return
}

func insertSites(ctx *grn.Ctx, sites ...Site) (err error) {
	table := ctx.Get("Site")
	defer ctx.ObjUnlinkDefer(&err, table)

	var column *grn.Obj
	column, err = ctx.ObjColumn(table, "title")
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, column)

	var value grn.Obj
	defer ctx.ObjUnlinkDefer(&err, &value)

	for _, site := range sites {
		var recordID grn.ID
		recordID, _, err = ctx.TableAdd(table, site.URL)
		if err != nil {
			return
		}

		grn.TextInit(&value, 0)
		ctx.TextPut(&value, site.Title)
		err = ctx.ObjSetValue(column, recordID, &value, grn.OBJ_SET)
		if err != nil {
			return
		}
	}
	return
}

func selectSites(ctx *grn.Ctx) (err error) {
	table := ctx.Get("Site")
	defer ctx.ObjUnlinkDefer(&err, table)

	var keyColumn *grn.Obj
	keyColumn, err = ctx.ObjColumn(table, "_key")
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, keyColumn)

	var titleColumn *grn.Obj
	titleColumn, err = ctx.ObjColumn(table, "title")
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, titleColumn)

	res := table

	count, err := ctx.TableSize(res)
	if err != nil {
		return
	}
	fmt.Printf("record count=%d\n", count)

	tc, err := ctx.TableCursorOpen(res, "", "", 0, -1, grn.CURSOR_ASCENDING)
	if err != nil {
		return
	}

	var buf grn.Obj
	defer ctx.ObjUnlinkDefer(&err, &buf)
	for {
		id := ctx.TableCursorNext(tc)
		if id == grn.ID_NIL {
			break
		}
		grn.TextInit(&buf, 0)
		grn.BulkRewind(&buf)
		ctx.ObjGetValue(keyColumn, id, &buf)
		key := grn.BulkHead(&buf)

		grn.TextInit(&buf, 0)
		grn.BulkRewind(&buf)
		ctx.ObjGetValue(titleColumn, id, &buf)
		title := grn.BulkHead(&buf)

		fmt.Printf("id=%d, key=%s, title=%s\n", id, key, title)
	}
	ctx.TableCursorClose(tc)

	return
}

func selectComSites(ctx *grn.Ctx) (err error) {
	table := ctx.Get("Site")
	defer ctx.ObjUnlinkDefer(&err, table)

	cond, v, err := ctx.ExprCreateForQuery(table)
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, cond)
	defer ctx.ObjUnlinkDefer(&err, v)

	flags := grn.EXPR_SYNTAX_QUERY | grn.EXPR_ALLOW_PRAGMA | grn.EXPR_ALLOW_COLUMN
	err = ctx.ExprParse(cond, "_key:@.com", nil, grn.OP_MATCH, grn.OP_AND, flags)
	if err != nil {
		return
	}

	res, err := ctx.TableSelect(table, cond, nil, grn.OP_OR)
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, res)

	count, err := ctx.TableSize(res)
	if err != nil {
		return
	}
	fmt.Printf("record count=%d\n", count)

	var keyColumn *grn.Obj
	keyColumn, err = ctx.ObjColumn(res, "_key")
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, keyColumn)

	var titleColumn *grn.Obj
	titleColumn, err = ctx.ObjColumn(res, "title")
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, titleColumn)

	tc, err := ctx.TableCursorOpen(res, "", "", 0, -1, grn.CURSOR_ASCENDING)
	if err != nil {
		return
	}

	var buf grn.Obj
	defer ctx.ObjUnlinkDefer(&err, &buf)
	for {
		id := ctx.TableCursorNext(tc)
		if id == grn.ID_NIL {
			break
		}
		grn.TextInit(&buf, 0)
		grn.BulkRewind(&buf)
		ctx.ObjGetValue(keyColumn, id, &buf)
		key := grn.BulkHead(&buf)

		grn.TextInit(&buf, 0)
		grn.BulkRewind(&buf)
		ctx.ObjGetValue(titleColumn, id, &buf)
		title := grn.BulkHead(&buf)

		fmt.Printf("id=%d, key=%s, title=%s\n", id, key, title)
	}
	ctx.TableCursorClose(tc)

	return
}

func run() (err error) {
	err = grn.Init()
	if err != nil {
		return
	}
	defer grn.FinDefer(&err)

	ctx, err := grn.CtxOpen(0)
	if err != nil {
		return
	}
	defer ctx.CloseDefer(&err)

	var db *grn.Obj
	db, err = ctx.DBOpenOrCreate("site.db", nil)
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, db)

	err = createSiteTable(ctx)
	if err != nil {
		return
	}

	err = insertSites(ctx,
		Site{URL: "http://example.org/", Title: "This is test record 1!"},
		Site{URL: "http://example.net/", Title: "test record 2."},
		Site{URL: "http://example.com/", Title: "test test record three."},
		Site{URL: "http://example.net/afr", Title: "test record four."},
		Site{URL: "http://example.org/aba", Title: "test test test record five."},
		Site{URL: "http://example.com/rab", Title: "test test test test record six."},
		Site{URL: "http://example.net/atv", Title: "test test test record seven."},
		Site{URL: "http://example.org/gat", Title: "test test record eight."},
		Site{URL: "http://example.com/vdw", Title: "test test record nine."},
	)
	if err != nil {
		return
	}

	err = selectSites(ctx)
	if err != nil {
		return
	}

	err = selectComSites(ctx)
	return
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
