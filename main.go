package cgoroonga

/*

func (t *Table) CreateQuery(str string) (*Query, error) {
	return nil, nil
}

func (t *Table) OpenCursor(min, max string, offset uint64, limit int64, flags int) (*TableCursor, error) {
	return nil, nil
}

func (t *Table) RecordCount() (uint64, error) {
	return 0, nil
}

type Query struct {
	context *Context
	table   *Table
}

func (q *Query) Exec() (*Records, error) {
	return nil, nil
}

func (q *Query) Close() error {
	return nil
}

type TableCursor struct {
	context  *Context
	table    *Table
	recordID int
}

func (c *TableCursor) Next() (RecordID, bool) {
	return 0, false
}

func (c *TableCursor) Close() error {
	return nil
}

type Records Table

type RecordID int

func (r RecordID) SetString(column *Column, s string) error {
	return nil
}

func (r RecordID) GetString(column *Column) (string, error) {
	return nil
}

func (r RecordID) SetTime(column *Column, t time.Time) error {
	return nil
}

func (r RecordID) GetTime(column *Column) (time.Time, error) {
	return nil
}

func FormatTimeForQuery(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()/1000, 10)
}
*/
