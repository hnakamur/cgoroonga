package main

import (
	"compress/bzip2"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	grn "github.com/hnakamur/cgoroonga"
)

type Page struct {
	Title     string `xml:"title"`
	Timestamp string `xml:"revision>timestamp"`
	Text      string `xml:"revision>text"`
}

func addArticle(ctx *grn.Ctx, table, textColumn, updatedAtColumn *grn.Obj, title, text, timestamp string) error {
	recordID, _, err := ctx.TableAdd(table, title)
	if err != nil {
		return err
	}

	var textValue grn.Obj
	grn.TextInit(&textValue, 0)
	ctx.TextPut(&textValue, text)
	err = ctx.ObjSetValue(textColumn, recordID, &textValue, grn.OBJ_SET)
	if err != nil {
		return err
	}

	updatedAt, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return err
	}
	var updatedAtValue grn.Obj
	grn.TimeInit(&updatedAtValue, 0)
	ctx.TimeSet(&updatedAtValue, updatedAt)
	return ctx.ObjSetValue(updatedAtColumn, recordID, &updatedAtValue, grn.OBJ_SET)
}

func closeFileDefer(err *error, f *os.File) {
	err2 := f.Close()
	if err2 != nil && *err == nil {
		*err = err2
	}
}

func run(dbFilename, wikiArticlesXmlBzip2Filename string) (err error) {
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

	keyType := ctx.At(grn.DB_SHORT_TEXT)
	table, err := ctx.TableOpenOrCreate("Articles", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, keyType, nil)
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, table)

	textColumn, err := ctx.ColumnOpenOrCreate(table, "text", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, ctx.At(grn.DB_TEXT))
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, textColumn)

	updatedAtColumn, err := ctx.ColumnOpenOrCreate(table, "updated_at", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, ctx.At(grn.DB_TIME))
	if err != nil {
		return
	}
	defer ctx.ObjUnlinkDefer(&err, updatedAtColumn)

	var file *os.File
	file, err = os.Open(wikiArticlesXmlBzip2Filename)
	if err != nil {
		return
	}
	defer closeFileDefer(&err, file)

	decoder := xml.NewDecoder(bzip2.NewReader(file))
	for {
		var t xml.Token
		t, err = decoder.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		if se, ok := t.(xml.StartElement); ok {
			if se.Name.Local == "page" {
				var p Page
				decoder.DecodeElement(&p, &se)

				err = addArticle(ctx, table, textColumn, updatedAtColumn,
					p.Title, p.Text, p.Timestamp)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

var dbFilenameVar string

func init() {
	flag.StringVar(&dbFilenameVar, "d", "wikipedia_ja.db", "database filename")
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s jawiki-20150422-pages-articles1.xml.bz2\n", os.Args[0])
		os.Exit(1)
	}
	wikiArticlesXmlBzip2Filename := flag.Arg(0)
	err := run(dbFilenameVar, wikiArticlesXmlBzip2Filename)
	if err != nil {
		panic(err)
	}
}
