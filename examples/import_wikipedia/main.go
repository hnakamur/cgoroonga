package main

import (
	"compress/bzip2"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"

	grn "github.com/hnakamur/cgoroonga"
)

type Page struct {
	Title string `xml:"title"`
	Text  string `xml:"revision>text"`
}

func addArticle(ctx *grn.Ctx, table, column *grn.Obj, title, text string) error {
	recordID, _, err := ctx.RecordAdd(table, title)
	if err != nil {
		return err
	}

	var value grn.Obj
	grn.TextInit(&value, 0)
	ctx.TextPut(&value, text)
	return ctx.ObjSetValue(column, recordID, &value, grn.OBJ_SET)
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
	defer ctx.ObjCloseDefer(&err, db)

	keyType := ctx.At(grn.DB_SHORT_TEXT)
	table, err := ctx.TableOpenOrCreate("Articles", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, keyType, nil)
	if err != nil {
		return
	}
	defer ctx.ObjCloseDefer(&err, table)

	columnType := ctx.At(grn.DB_TEXT)
	column, err := ctx.ColumnOpenOrCreate(table, "text", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, columnType)
	if err != nil {
		return
	}
	defer ctx.ObjCloseDefer(&err, column)

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

				err = addArticle(ctx, table, column, p.Title, p.Text)
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
