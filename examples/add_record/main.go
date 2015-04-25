package main

import (
	"fmt"

	grn "github.com/hnakamur/cgoroonga"
)

func handler(ctx *grn.Ctx, db *grn.Obj) (err error) {
	keyType := ctx.At(grn.DB_SHORT_TEXT)
	table, err := ctx.TableOpenOrCreate("table1", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, keyType, nil)
	if err != nil {
		return
	}
	fmt.Printf("table=%x\n", table)
	defer func() {
		err2 := ctx.ObjClose(table)
		if err2 != nil && err == nil {
			err = err2
		}
	}()

	columnType := ctx.At(grn.DB_TEXT)
	column, err := ctx.ColumnOpenOrCreate(table, "col", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, columnType)
	if err != nil {
		return
	}
	fmt.Printf("column=%x\n", column)
	defer func() {
		err2 := ctx.ObjClose(column)
		if err2 != nil && err == nil {
			err = err2
		}
	}()

	recordID, added, err := ctx.RecordAdd(table, "rec1")
	if err != nil {
		return
	}
	fmt.Printf("recordID=%d, added=%v\n", recordID, added)

	var value grn.Obj
	grn.TextInit(&value, 0)
	ctx.TextPut(&value, "groonga world")
	fmt.Printf("value=%v\n", value)
	err = ctx.ObjSetValue(column, recordID, &value, grn.OBJ_SET)
	if err != nil {
		return
	}
	fmt.Printf("SetValue done\n")

	var bulk grn.Obj
	grn.TextInit(&bulk, 0)
	grn.BulkRewind(&bulk)
	ctx.ObjGetValue(column, recordID, &bulk)
	fmt.Printf("bulk=%s\n", grn.BulkHead(&bulk))
	return
}

func run() error {
	return grn.WithCtx(0, func(ctx *grn.Ctx) error {
		return ctx.WithDB("hello.db", nil, handler)
	})
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
