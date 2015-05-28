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

var db *grn.DB
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

const (
	defaultPageSize   = 10
	maxPageSize       = 100
	viewablePageCount = 9
)

func getIndex(c *gin.Context) {
	c.Request.ParseForm()
	q := c.Request.Form.Get("q")
	timespan := c.Request.Form.Get("timespan")
	var err error
	var limitCount int = defaultPageSize
	var page int = 1
	var numPages int = 1
	var matchedCount uint
	viewablePages := []int{}
	results := []Result{}
	if q != "" {
		page, err = formIntValue(c, "page", 1)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		limitCount, err = formIntValue(c, "limit", defaultPageSize)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		startDate := calcStartDate(timespan)

		table, err := db.OpenTable("Articles")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer table.Close()

		expr, defaultColumns, err := buildSearchExpr(table, q, startDate)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer expr.Close()
		defer defaultColumns.Close()

		res, err := table.Select(expr, nil, grn.OP_OR)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer res.Close()

		matchedCount, err = res.RecordCount()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		normalizePageAndLimit(int(matchedCount), &page, &limitCount)
		offset := (page - 1) * limitCount

		sorted, err := res.Sort("-updated_at", offset, limitCount)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer sorted.Close()

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

		keyColumn, err := sorted.OpenColumn("_key")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		textColumn, err := sorted.OpenColumn("text")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		updatedAtColumn, err := sorted.OpenColumn("updated_at")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		tc, err := sorted.OpenTableCursor("", "", 0, limitCount, grn.CURSOR_ASCENDING)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer tc.Close()

		flags := grn.SNIP_COPY_TAG | grn.SNIP_SKIP_LEADING_SPACES
		textMaxLen := 200
		snippet := expr.Snippet(flags, textMaxLen, 1, true,
			[][]string{
				[]string{"<b>", "</b>"},
			})
		for {
			id, hasNext := tc.Next()
			if !hasNext {
				break
			}

			title := id.GetString(keyColumn)
			text := id.GetString(textColumn)
			updatedAt := id.GetTime(updatedAtColumn)

			snipResults, err := snippet.Exec(text)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			if len(snipResults) > 0 {
				text = "..." + snipResults[0] + "..."

			} else {
				if utf8.RuneCountInString(text) >= textMaxLen {
					r := []rune(text)
					text = string(r[:textMaxLen]) + "..."
				}
				text = html.EscapeString(text)
			}

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

func normalizePageAndLimit(matchedCount int, page, limit *int) {
	if *limit < defaultPageSize {
		*limit = defaultPageSize
	} else if *limit > maxPageSize {
		*limit = maxPageSize
	}

	pageSize := *limit
	maxPage := (int(matchedCount)-1)/pageSize + 1
	if *page < 1 {
		*page = 1
	} else if *page > maxPage {
		*page = maxPage
	}
}

func buildSearchExpr(table *grn.Table, q string, startDate time.Time) (expr, defaultColumns *grn.Expr, err error) {
	defaultColumns, err = table.CreateQuery("")
	if err != nil {
		return
	}
	err = defaultColumns.Parse("_key,text", nil, grn.OP_MATCH, grn.OP_AND,
		grn.EXPR_SYNTAX_SCRIPT)
	if err != nil {
		return
	}

	expr, err = table.CreateQuery("")
	if err != nil {
		return
	}

	err = expr.Parse(q, defaultColumns, grn.OP_MATCH, grn.OP_AND,
		grn.EXPR_SYNTAX_QUERY|grn.EXPR_ALLOW_PRAGMA|grn.EXPR_ALLOW_COLUMN)
	if err != nil {
		return
	}
	if !startDate.IsZero() {
		usecStartTime := startDate.UnixNano() / 1000
		filter := "updated_at >= " + strconv.FormatInt(usecStartTime, 10)
		err = expr.Parse(filter, defaultColumns, grn.OP_MATCH, grn.OP_AND,
			grn.EXPR_SYNTAX_SCRIPT)
		if err != nil {
			return
		}

		if q != "" {
			err = expr.AppendOp(grn.OP_AND, 2)
			if err != nil {
				return
			}
		}
	}
	return
}

func main() {
	flag.Parse()

	var err error
	err = grn.Init()
	if err != nil {
		panic(err)
	}
	defer grn.Terminate()

	ctx, err := grn.NewContext()
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	db, err = ctx.OpenDB(dbFilename)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", getIndex)
	r.Run(listenAddress)
}
