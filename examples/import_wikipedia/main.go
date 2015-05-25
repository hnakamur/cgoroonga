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

func addArticle(table *grn.Table, title, text, timestamp string) error {
	recordID, _, err := table.AddRecord(title)
	if err != nil {
		return err
	}

	err = recordID.SetString(table.Column("text"), text)
	if err != nil {
		return err
	}

	updatedAt, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return err
	}
	return recordID.SetTime(table.Column("updated_at"), updatedAt)
}

func run(dbFilename, wikiArticlesXmlBzip2Filename string) error {
	err := grn.Init()
	if err != nil {
		return err
	}
	defer grn.Terminate()

	ctx, err := grn.NewContext()
	if err != nil {
		return err
	}
	defer ctx.Close()

	db, err := ctx.OpenOrCreateDB(dbFilename)
	if err != nil {
		return err
	}
	defer db.Close()

	table, err := db.OpenOrCreateTable("Articles", "",
		grn.OBJ_TABLE_HASH_KEY|grn.OBJ_PERSISTENT, grn.DB_SHORT_TEXT)
	if err != nil {
		return err
	}

	_, err = table.OpenOrCreateColumn("text", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, grn.DB_TEXT)
	if err != nil {
		return err
	}

	_, err = table.OpenOrCreateColumn("updated_at", "",
		grn.OBJ_PERSISTENT|grn.OBJ_COLUMN_SCALAR, grn.DB_TIME)
	if err != nil {
		return err
	}

	var file *os.File
	file, err = os.Open(wikiArticlesXmlBzip2Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := xml.NewDecoder(bzip2.NewReader(file))
	for {
		var t xml.Token
		t, err = decoder.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		if se, ok := t.(xml.StartElement); ok {
			if se.Name.Local == "page" {
				var p Page
				decoder.DecodeElement(&p, &se)

				err = addArticle(table, p.Title, p.Text, p.Timestamp)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
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
