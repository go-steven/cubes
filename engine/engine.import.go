package engine

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
)

type ImportEngine struct {
}

func NewImportEngine() *ImportEngine {
	return &ImportEngine{}
}

func (e *ImportEngine) ImportCsvFile(csvFile string, dbname, table string) error {
	logger.Infof("sqlite dbname: %s", dbname)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer db.Close()

	sqlite := source.NewSqlite(dbname, db)

	tags := make(map[string][]*metadata.TagMapping)
	return create_table_from_csv(sqlite, table, csvFile, true, tags)
}

func (e *ImportEngine) ImportJsonFile(jsonFile string, dbname, table string) error {
	logger.Infof("sqlite dbname: %s", dbname)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer db.Close()

	sqlite := source.NewSqlite(dbname, db)

	tags := make(map[string][]*metadata.TagMapping)
	return create_table_from_json(sqlite, table, jsonFile, tags)
}
