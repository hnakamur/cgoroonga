package main

import (
	"flag"
	"fmt"
	"os"

	grn "github.com/hnakamur/cgoroonga"
)

func queryArticles(ctx *grn.Ctx, query string) (err error) {
	table := ctx.Get("Articles")
	defer ctx.ObjUnlinkDefer(&err, table)

	cond, v, err := ctx.ExprCreateForQuery(table)
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, cond)
	defer ctx.ObjUnlinkDefer(&err, v)

	flags := grn.EXPR_SYNTAX_QUERY | grn.EXPR_ALLOW_PRAGMA | grn.EXPR_ALLOW_COLUMN
	err = ctx.ExprParse(cond, query, nil, grn.OP_MATCH, grn.OP_AND, flags)
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

	var textColumn *grn.Obj
	textColumn, err = ctx.ObjColumn(res, "text")
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, textColumn)

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
		ctx.ObjGetValue(textColumn, id, &buf)
		text := grn.BulkHead(&buf)

		r := []rune(text)
		if len(r) >= 200 {
			text = string(r[:200]) + "…"
		}

		fmt.Printf("id=%d, key=%s, text=%s\n", id, key, text)
	}
	ctx.TableCursorClose(tc)

	return
}

func run(dbFilename, query string) (err error) {
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
	db, err = ctx.DBOpenOrCreate(dbFilename, nil)
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, db)

	err = queryArticles(ctx, query)
	return
}

var dbFilename string

func init() {
	flag.StringVar(&dbFilename, "d", "wikipedia_ja.db", "database filename")
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s [-d dbFilename] query\n", os.Args[0])
		os.Exit(1)
	}
	query := flag.Arg(0)
	err := run(dbFilename, query)
	if err != nil {
		panic(err)
	}
}
