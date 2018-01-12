package source

// common query interface
type Row map[string]string
type Rows []Row
type QueryInterface interface {
	Query(sql string, fields []string) (Rows, error)
}
