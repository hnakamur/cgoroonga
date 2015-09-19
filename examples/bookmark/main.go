package main

import (
	grn "github.com/hnakamur/cgoroonga"
)

func openOrCreateDB(ctx *grn.Context) (*grn.DB, error) {
	db, err := ctx.OpenOrCreateDB("bookmark.db")
	if err != nil {
		return nil, err
	}

	shortTextType := ctx.At(grn.DB_SHORT_TEXT)
	defer shortTextType.Unlink()

	tagTable, err := db.CreateTable("Tag", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, shortTextType)
	if err != nil {
		return nil, err
	}

	bookmarkTable, err := db.CreateTable("Bookmark", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, shortTextType)
	if err != nil {
		return nil, err
	}
	_, err = bookmarkTable.CreateColumn("title", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, shortTextType)
	if err != nil {
		return nil, err
	}
	_, err = bookmarkTable.CreateColumn("tags", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_VECTOR, tagTable.AsObj())
	if err != nil {
		return nil, err
	}

	_, err = tagTable.CreateColumn("index_tags", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_INDEX, bookmarkTable.AsObj(), "tags")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	err := grn.Init()
	if err != nil {
		panic(err)
	}
	defer grn.Terminate()

	ctx, err := grn.NewContext()
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	db, err := openOrCreateDB(ctx)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
