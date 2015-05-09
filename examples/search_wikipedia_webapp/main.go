package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	grn "github.com/hnakamur/cgoroonga"
)

func staticFileHandler(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

func formIntValue(r *http.Request, key string, defaultValue int) (int, error) {
	strValue := r.FormValue(key)
	intValue := defaultValue
	if strValue != "" {
		var err error
		intValue, err = strconv.Atoi(strValue)
		if err != nil {
			return defaultValue,
				fmt.Errorf("int parameter expected, but got \"%s\" for parameter \"%s\"",
					strValue, key)
		}
	}
	return intValue, nil
}

func formDateValue(r *http.Request, key string, defaultValue time.Time) (time.Time, error) {
	strValue := r.FormValue(key)
	timeValue := defaultValue
	if strValue != "" {
		var err error
		timeValue, err = time.Parse("2006-1-2", strValue)
		if err != nil {
			return defaultValue,
				fmt.Errorf("date parameter expected, but got \"%s\" for parameter \"%s\"",
					strValue, key)
		}
	}
	return timeValue, nil
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")

	offset, err := formIntValue(r, "offset", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limitCount, err := formIntValue(r, "limit", 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := formDateValue(r, "sd", time.Unix(0, 0))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bw := bufio.NewWriter(w)

	table := ctx.Get("Articles")
	defer ctx.ObjUnlinkDefer(&err, table)

	var res *grn.Obj
	if q != "" || !startDate.IsZero() {
		cond, v, err := ctx.ExprCreateForQuery(table)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer ctx.ObjUnlinkDefer(&err, cond)
		defer ctx.ObjUnlinkDefer(&err, v)

		if q != "" {
			flags := grn.EXPR_SYNTAX_QUERY | grn.EXPR_ALLOW_PRAGMA | grn.EXPR_ALLOW_COLUMN
			err = ctx.ExprParse(cond, q, nil, grn.OP_MATCH, grn.OP_AND, flags)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if !startDate.IsZero() {
			usecStartTime := startDate.UnixNano() / 1000
			filter := "updated_at >= " + strconv.FormatInt(usecStartTime, 10)
			err = ctx.ExprParse(cond, filter, nil, grn.OP_MATCH, grn.OP_AND,
				grn.EXPR_SYNTAX_SCRIPT)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if q != "" {
				err = ctx.ExprAppendOp(cond, grn.OP_AND, 2)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		res, err = ctx.TableSelect(table, cond, nil, grn.OP_OR)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer ctx.ObjUnlinkDefer(&err, res)
	} else {
		res = table
	}

	count, err := ctx.TableSize(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resultCount uint = count
	if uint(limitCount) < count {
		resultCount = uint(limitCount)
	}

	var keyColumn *grn.Obj
	keyColumn, err = ctx.ObjColumn(res, "_key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.ObjUnlinkDefer(&err, keyColumn)

	var textColumn *grn.Obj
	textColumn, err = ctx.ObjColumn(res, "text")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.ObjUnlinkDefer(&err, textColumn)

	var updatedAtColumn *grn.Obj
	updatedAtColumn, err = ctx.ObjColumn(res, "updated_at")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.ObjUnlinkDefer(&err, updatedAtColumn)

	tc, err := ctx.TableCursorOpen(res, "", "", offset, limitCount, grn.CURSOR_ASCENDING)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bw.WriteString(fmt.Sprintf(`{"matchedCount":%d,`, count))
	bw.WriteString(fmt.Sprintf(`"resultCount":%d,`, resultCount))
	bw.WriteString(`"results":[`)
	var buf grn.Obj
	defer ctx.ObjUnlinkDefer(&err, &buf)
	var jsonBuf []byte
	first := true
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

		grn.TimeInit(&buf, 0)
		ctx.ObjGetValue(updatedAtColumn, id, &buf)
		updatedAt := grn.TimeValue(&buf)

		r := []rune(text)
		if len(r) >= 200 {
			text = string(r[:200]) + "…"
		}

		if first {
			first = false
		} else {
			bw.WriteRune(',')
			bw.WriteRune('\n')
		}
		bw.WriteString(`{"title":`)
		jsonBuf, err = json.Marshal(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bw.Write(jsonBuf)
		bw.WriteString(`,"text":`)
		jsonBuf, err = json.Marshal(text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bw.Write(jsonBuf)
		bw.WriteString(`,"updated_at":`)
		bw.WriteString(strconv.FormatInt(updatedAt.Unix(), 10))
		bw.WriteString(`}`)
	}
	ctx.TableCursorClose(tc)
	bw.WriteString(`]}`)
	bw.Flush()

	return
}

var ctx *grn.Ctx
var dbFilename string
var listenAddress string

func init() {
	flag.StringVar(&dbFilename, "d", "wikipedia_ja.db", "database filename")
	flag.StringVar(&listenAddress, "l", ":8080", "listen address (address:port)")
}

func main() {
	flag.Parse()

	var err error
	err = grn.Init()
	if err != nil {
		panic(err)
	}
	defer grn.FinDefer(&err)

	ctx, err = grn.CtxOpen(0)
	if err != nil {
		panic(err)
	}
	defer ctx.CloseDefer(&err)

	var db *grn.Obj
	db, err = ctx.DBOpenOrCreate(dbFilename, nil)
	if err != nil {
		panic(err)
	}
	defer ctx.ObjUnlinkDefer(&err, db)

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/", staticFileHandler("public/index.html"))
	http.HandleFunc("/js/mithril.js", staticFileHandler("public/js/mithril.js"))
	http.HandleFunc("/js/observable.js", staticFileHandler("public/js/observable.js"))
	http.ListenAndServe(listenAddress, nil)
}
