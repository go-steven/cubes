package engine

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"strings"
)

type ExportEngine struct {
	tmpdir string
}

func NewExportEngine() *ExportEngine {
	engine := &ExportEngine{
		tmpdir: "./",
	}
	return engine
}

func (e *ExportEngine) ExportTable(dbname, table string) (source.Rows, error) {
	logger.Infof("sqlite dbname: %s", dbname)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer db.Close()

	sqlite := source.NewSqlite(dbname, db)
	fields, err := sqlite.GetTableFields(table)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if len(fields) == 0 {
		return nil, errors.New(fmt.Sprintf("Table[%s] not found.", table))
	}
	logger.Infof("fields : %v", fields)
	sql := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s", strings.Join(fields, ", "), table, fields[0])
	return sqlite.Query(sql, fields)
}
