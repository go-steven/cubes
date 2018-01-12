package source

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"regexp"
	"strings"
	"sync"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
)

type Sqlite struct {
	db     *sql.DB
	dbname string
	m      *sync.Mutex
}

func NewSqlite(dbname string, db *sql.DB) *Sqlite {
	return &Sqlite{
		db:     db,
		dbname: dbname,
		m:      new(sync.Mutex),
	}
}
func (s *Sqlite) Cleanup(removeFile bool) error {
	if s.db != nil {
		s.m.Lock()
		s.db.Close()
		s.m.Unlock()
	}
	if s.dbname != "" && removeFile {
		os.Remove(s.dbname)
	}
	return nil
}

func (s *Sqlite) CreateTable(table string, fields [][]string, dropExist bool) error {
	table = strings.TrimSpace(table)
	if table == "" {
		return errors.New("ERROR: Empty table name")
	}
	if len(fields) <= 0 {
		return errors.New(fmt.Sprintf("No fields for table: %s", table))
	}
	if dropExist {
		if err := s.DropTable(table); err != nil {
			logger.Error(err)
			return err
		}
	}
	var sqlBuffer bytes.Buffer

	sqlBuffer.WriteString(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (`, table))
	var v0, v1 string
	cnt := 0
	for _, v := range fields {
		v0 = strings.TrimSpace(v[0])
		v1 = strings.TrimSpace(v[1])
		if v0 == "" || v1 == "" {
			return errors.New("ERROR: Have empty field.")
		}

		if cnt > 0 {
			sqlBuffer.WriteString(", ")
		}
		sqlBuffer.WriteString(fmt.Sprintf("`%s` %s", v0, v1))

		cnt++
	}
	sqlBuffer.WriteString(");")
	sqlStmt := sqlBuffer.String()
	logger.Infof("%s", sqlStmt)

	s.m.Lock()
	_, err := s.db.Exec(sqlStmt)
	s.m.Unlock()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Sqlite) DropTable(table string) error {
	s.m.Lock()
	_, err := s.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s ", table))
	s.m.Unlock()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Sqlite) Query(sql string, fields []string) (Rows, error) {
	fmt.Printf("Sqlite run SQL: %s, fields:%v\n", sql, fields)
	if len(fields) == 0 {
		return nil, errors.New("ERROR: no return fields given.")
	}
	s.m.Lock()
	rows, err := s.db.Query(sql)
	//fmt.Println(Json(rows))
	s.m.Unlock()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	dst := []interface{}{}
	for _, _ = range fields {
		var s string
		dst = append(dst, &s)
	}
	ret := Rows{}
	for rows.Next() {
		err = rows.Scan(dst...)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		retRow := make(Row)
		for k, v := range fields {
			if v == "" {
				continue
			}
			if v[0] == '`' && v[len(v) - 1] == '`' {
				v = v[1:len(v)-2]
			}
			
			val, _ := dst[k].(*string)
			retRow[v] = *val
		}

		ret = append(ret, retRow)
	}
	logger.Infof("SQL result: %d", len(ret))
	return ret, nil
}

func (s *Sqlite) InsertTable(table string, data Rows) error {
	if len(data) == 0 {
		return nil
	}
	var sqlBuffer bytes.Buffer

	sqlBuffer.WriteString(fmt.Sprintf(`INSERT INTO %s (`, table))
	cnt := 0
	fields := []string{}
	for k, _ := range data[0] {
		fields = append(fields, k)
		if cnt == 0 {
			sqlBuffer.WriteString(fmt.Sprintf("`%s`", k))
		} else {
			sqlBuffer.WriteString(fmt.Sprintf(", `%s`", k))
		}
		cnt += 1
	}
	sqlBuffer.WriteString(") ")
	logger.Infof("INSERT INTO %s, rows:%d, fields: %s", table, len(data), utils.Json(fields))

	cnt = 0
	var sub_sql bytes.Buffer
	var wg sync.WaitGroup
	for _, row := range data {
		if cnt > 0 {
			sub_sql.WriteString(" UNION ALL ")
		}
		sub_sql.WriteString(" SELECT ")
		for k, v := range fields {
			if k == 0 {
				sub_sql.WriteString(fmt.Sprintf("'%s' AS `%s`", escape_sqlite_special_chars(row[v]), v))
			} else {
				sub_sql.WriteString(fmt.Sprintf(", '%s' AS `%s`", escape_sqlite_special_chars(row[v]), v))
			}
		}
		//logger.Infof("sub_sql : %s", sub_sql)
		cnt++
		if cnt >= 200 {
			var query bytes.Buffer
			query.WriteString(" ")
			query.Write(sqlBuffer.Bytes())
			query.WriteString(" ")
			query.Write(sub_sql.Bytes())
			query.WriteString(" ")

			//logger.Infof("query : %s", query)
			wg.Add(1)
			go func(sql bytes.Buffer) {
				defer wg.Done()

				if err := s.insert_into_db(sql); err != nil {
					logger.Error(err)
				}
			}(query)

			cnt = 0
			sub_sql.Reset()
		}
	}

	if cnt > 0 {
		var query bytes.Buffer
		query.WriteString(" ")
		query.Write(sqlBuffer.Bytes())
		query.WriteString(" ")
		query.Write(sub_sql.Bytes())
		query.WriteString(" ")
		wg.Add(1)
		go func(sql bytes.Buffer) {
			defer wg.Done()

			if err := s.insert_into_db(sql); err != nil {
				logger.Error(err)
			}
		}(query)
	}
	wg.Wait()

	return nil
}
func (s *Sqlite) insert_into_db(query bytes.Buffer) error {
	//logger.Info("insert_into_db.")
	s.m.Lock()
	defer s.m.Unlock()

	//logger.Info("tx begin.")
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	//logger.Info(query)
	stmt, err := tx.Prepare(query.String())
	if err != nil {
		logger.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		logger.Error(err)
		return err
	}
	tx.Commit()
	logger.Info("tx commit.")
	return nil
}

func (s *Sqlite) GetTableFields(table string) ([]string, error) {
	tableInfo, err := s.Query(fmt.Sprintf("select sql from sqlite_master WHERE type = 'table' AND name = '%s'", table), []string{
		"sql",
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if len(tableInfo) == 0 {
		return nil, errors.New(fmt.Sprintf("Table[%s] not found.", table))
	}
	reg := regexp.MustCompile(`\((.*)\)`)
	sql := reg.FindString(tableInfo[0]["sql"])
	if sql == "" {
		return nil, errors.New(fmt.Sprintf("Table[%s] not found.", table))
	}
	dmlFields := strings.Split(sql[1:len(sql)-1], ", ")
	//fmt.Printf("dmlFields:%v\n", dmlFields)

	fields := []string{}
	for _, v := range dmlFields {
		vals := strings.Split(v, " ")
		if len(vals) > 0 {
			fields = append(fields, strings.ToLower(vals[0]))
		}
	}

	return fields, nil
}

var sqlite_special_chars = map[string]string{
	"/": "//",
	"'": "''",
	"[": "/[",
	"]": "/]",
	"%": "/%",
	"&": "/&",
	"_": "/_",
	"(": "/(",
	")": "/)",
}

//
func escape_sqlite_special_chars(sql string) string {
	var (
		ret bytes.Buffer
		str string
	)
	for _, s := range sql {
		str = string(s)
		v, ok := sqlite_special_chars[str]
		if ok {
			ret.WriteString(v)
		} else {
			ret.WriteString(str)
		}
	}
	return ret.String()
}
