package main

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gopkg.in/flosch/pongo2.v3"

	grn "github.com/hnakamur/cgoroonga"
)

var ctx *grn.Ctx
var dbFilename string
var listenAddress string

func init() {
	flag.StringVar(&dbFilename, "d", "wikipedia_ja.db", "database filename")
	flag.StringVar(&listenAddress, "l", ":8080", "listen address (address:port)")
}

func formIntValue(c *gin.Context, key string, defaultValue int) (int, error) {
	strValue := c.Request.FormValue(key)
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

func formDateValue(c *gin.Context, key string, defaultValue time.Time) (time.Time, error) {
	strValue := c.Request.FormValue(key)
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

type Result struct {
	URL       string
	Title     string
	Content   string
	UpdatedAt string
}

func titleToURL(title string) string {
	return fmt.Sprintf("https://ja.wikipedia.org/wiki/%s",
		url.QueryEscape(title))
}

func calcStartDate(timespan string) time.Time {
	var duration time.Duration
	switch timespan {
	case "week":
		duration = 7 * 24 * time.Hour
	case "month":
		duration = 31 * 24 * time.Hour
	case "year":
		duration = 365 * 24 * time.Hour
	default:
		return time.Unix(0, 0)
	}
	return time.Now().Add(-duration)
}

const viewablePageCount = 9

func getIndex(c *gin.Context) {
	c.Request.ParseForm()
	q := c.Request.Form.Get("q")
	timespan := c.Request.Form.Get("timespan")
	var err error
	var limitCount int = 10
	var page int = 1
	var numPages int = 1
	var matchedCount uint
	viewablePages := []int{}
	results := []Result{}
	var cond *grn.Obj
	if q != "" {
		page, err = formIntValue(c, "page", 1)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if page < 1 {
			page = 1
		}

		limitCount, err = formIntValue(c, "limit", 10)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		offset := (page - 1) * limitCount

		startDate := calcStartDate(timespan)

		table := ctx.Get("Articles")
		defer ctx.ObjUnlinkDefer(&err, table)

		var res *grn.Obj
		var v *grn.Obj
		cond, v, err = ctx.ExprCreateForQuery(table)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer ctx.ObjUnlinkDefer(&err, cond)
		defer ctx.ObjUnlinkDefer(&err, v)

		query := fmt.Sprintf("_key:@%s OR text:@%s", q, q)
		flags := grn.EXPR_SYNTAX_QUERY | grn.EXPR_ALLOW_PRAGMA | grn.EXPR_ALLOW_COLUMN
		err = ctx.ExprParse(cond, query, nil, grn.OP_MATCH, grn.OP_AND, flags)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if !startDate.IsZero() {
			usecStartTime := startDate.UnixNano() / 1000
			filter := "updated_at >= " + strconv.FormatInt(usecStartTime, 10)
			err = ctx.ExprParse(cond, filter, nil, grn.OP_MATCH, grn.OP_AND,
				grn.EXPR_SYNTAX_SCRIPT)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			if q != "" {
				err = ctx.ExprAppendOp(cond, grn.OP_AND, 2)
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
			}
		}

		res, err = ctx.TableSelect(table, cond, nil, grn.OP_OR)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer ctx.ObjUnlinkDefer(&err, res)

		matchedCount, err = ctx.TableSize(res)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		numPages = ((int(matchedCount) - 1) / limitCount) + 1
		if page > numPages {
			page = numPages
		}
		startPage := page - (viewablePageCount / 2)
		if startPage < 1 {
			startPage = 1
		}
		endPage := startPage + viewablePageCount
		if endPage > numPages+1 {
			endPage = numPages + 1
			if endPage-startPage < viewablePageCount {
				startPage = endPage - viewablePageCount
				if startPage < 1 {
					startPage = 1
				}
			}
		}
		for i := startPage; i < endPage; i++ {
			viewablePages = append(viewablePages, i)
		}

		var keyColumn *grn.Obj
		keyColumn, err = ctx.ObjColumn(res, "_key")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer ctx.ObjUnlinkDefer(&err, keyColumn)

		var textColumn *grn.Obj
		textColumn, err = ctx.ObjColumn(res, "text")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer ctx.ObjUnlinkDefer(&err, textColumn)

		var updatedAtColumn *grn.Obj
		updatedAtColumn, err = ctx.ObjColumn(res, "updated_at")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer ctx.ObjUnlinkDefer(&err, updatedAtColumn)

		tc, err := ctx.TableCursorOpen(res, "", "", offset, limitCount, grn.CURSOR_ASCENDING)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer ctx.TableCursorClose(tc)

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
			title := grn.BulkHead(&buf)

			grn.TextInit(&buf, 0)
			grn.BulkRewind(&buf)
			ctx.ObjGetValue(textColumn, id, &buf)
			text := grn.BulkHead(&buf)

			flags := grn.SNIP_COPY_TAG | grn.SNIP_SKIP_LEADING_SPACES
			textManLen := 200
			snippet := ctx.ExprSnippet(cond, flags, textManLen, 1, true,
				[][]string{
					[]string{"<b>", "</b>"},
				})
			snipResults, err := ctx.SnipExec(snippet, text)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			if len(snipResults) > 0 {
				text = "..." + snipResults[0] + "..."

			} else {
				if utf8.RuneCountInString(text) >= textManLen {
					r := []rune(text)
					text = string(r[:textManLen]) + "..."
				}
				text = html.EscapeString(text)
			}

			grn.TimeInit(&buf, 0)
			ctx.ObjGetValue(updatedAtColumn, id, &buf)
			updatedAt := grn.TimeValue(&buf)

			url := titleToURL(title)

			results = append(results, Result{
				URL:       url,
				Title:     title,
				Content:   text,
				UpdatedAt: updatedAt.Format("2006/01/02"),
			})
		}
	}

	url_ := fmt.Sprintf("%s?q=%s&timespan=%s&limit=%d",
		c.Request.URL.Path, url.QueryEscape(q), timespan, limitCount)
	tpl, err := pongo2.FromFile("templates/index.html")
	if err != nil {
		c.String(500, "Internal Server Error")
	}
	err = tpl.ExecuteWriter(pongo2.Context{
		"url":           url_,
		"q":             q,
		"timespan":      timespan,
		"matchedCount":  matchedCount,
		"results":       results,
		"page":          page,
		"numPages":      numPages,
		"viewablePages": viewablePages,
	}, c.Writer)
	if err != nil {
		c.String(500, "Internal Server Error")
	}
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

	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", getIndex)
	r.Run(listenAddress)
}
