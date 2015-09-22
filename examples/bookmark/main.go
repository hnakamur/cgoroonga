package main

import (
	"fmt"
	"strings"

	grn "github.com/hnakamur/cgoroonga"
)

func openOrCreateDB(ctx *grn.Context) (*grn.DB, error) {
	db, err := ctx.OpenOrCreateDB("/tmp/bookmark.db")
	if err != nil {
		return nil, err
	}

	shortTextType := ctx.At(grn.DB_SHORT_TEXT)
	defer shortTextType.Unlink()
	textType := ctx.At(grn.DB_TEXT)
	defer textType.Unlink()
	longTextType := ctx.At(grn.DB_LONG_TEXT)
	defer longTextType.Unlink()

	fmt.Println("creating Tag table")
	tagTable, err := db.OpenOrCreateTable("Tag", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, shortTextType)
	if err != nil {
		return nil, err
	}
	defer tagTable.Close()

	fmt.Println("creating Bookmark table")
	bookmarkTable, err := db.OpenOrCreateTable("Bookmark", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, shortTextType)
	if err != nil {
		return nil, err
	}
	defer bookmarkTable.Close()
	fmt.Println("creating title column")
	_, err = bookmarkTable.OpenOrCreateColumn("title", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, textType)
	if err != nil {
		return nil, err
	}
	fmt.Println("creating note column")
	_, err = bookmarkTable.OpenOrCreateColumn("note", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, longTextType)
	if err != nil {
		return nil, err
	}
	fmt.Println("creating tags column")
	_, err = bookmarkTable.OpenOrCreateColumn("tags", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_VECTOR, tagTable.AsObj())
	if err != nil {
		return nil, err
	}
	fmt.Println("created tables")

	keyColumn, err := bookmarkTable.OpenKeyColumn()
	if err != nil {
		return nil, err
	}
	defer keyColumn.Close()
	fmt.Printf("key type=%d\n", keyColumn.DataType())
	fmt.Printf("title type=%d\n", bookmarkTable.Column("title").DataType())
	fmt.Printf("note type=%d\n", bookmarkTable.Column("note").DataType())
	fmt.Printf("tags type=%d\n", bookmarkTable.Column("tags").DataType())

	_, err = tagTable.OpenOrCreateColumn("index_tags", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_INDEX, bookmarkTable.AsObj(), "tags")
	if err != nil {
		return nil, err
	}
	fmt.Println("created index_tags column")

	// url := "http://golang.org/pkg/"
	// title := "The Go Programming language API document"
	// // NOTE: Include an empty string to test an empty string too
	// tags := []string{"go", "", "api-docs"}

	// recordID, _, err := bookmarkTable.AddRecord(url)
	// if err != nil {
	// 	return nil, err
	// }
	// err = bookmarkTable.Column("title").SetString(recordID, title)
	// if err != nil {
	// 	return nil, err
	// }
	// err = bookmarkTable.Column("tags").SetStringArray(recordID, tags)
	// if err != nil {
	// 	return nil, err
	// }

	recordID := grn.ID(2)
	urlGot := keyColumn.GetString(recordID)
	fmt.Printf("urlGot=%s\n", urlGot)
	titleGot := bookmarkTable.Column("title").GetString(recordID)
	fmt.Printf("titleGot=%s\n", titleGot)
	tagsGot, err := bookmarkTable.Column("tags").GetStringArray(recordID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("tagsGot=%s\n", strings.Join(tagsGot, ","))

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
