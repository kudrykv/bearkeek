package bearkeek

import (
	"database/sql"
	"strings"

	"github.com/mattn/go-sqlite3"
)

// nolint:gochecknoinits
func init() {
	sql.Register("sqlite_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("utf8lower", strings.ToLower, true)
		},
	})
}
