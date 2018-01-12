package source

import (
	"bytes"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"regexp"
	"strings"
)

const (
	WITH_NO_SPACE = "no-space"
	WITH_SPACES   = "has_spaces"
)

var forbid_keywords = map[string][]string{
	WITH_NO_SPACE: []string{
		";",
	},
	WITH_SPACES: []string{
		"ALTER",
		"AUTOCOMMIT",
		"BEGIN",
		"CALL",
		"COMMIT",
		"CREATE",
		"DELETE",
		"DENY",
		"DO",
		"DROP",
		"END",
		"EXPLAIN",
		"FLUSH",
		"FUNCTION",
		"GRANT",
		"HANDLER",
		"INDEX",
		"INSERT",
		"INTO",
		"LOAD",
		"LOCK",
		"PROCEDURE",
		"READ",
		"RENAME",
		"REPLACE",
		"REVOKE",
		"ROLLBACK",
		"SAVEPOINT",
		"SET",
		"SHOW",
		"START",
		"TABLE",
		"TABLES",
		"TRANSACTION",
		"TRUNCATE",
		"UNLOCK",
		"UPDATE",
		"VALUES",
		"VIEW",
		"WITH",
	},
}

func SqlForbidRegexp() *regexp.Regexp {
	var buffer bytes.Buffer
	buffer.WriteString(`.*`)
	for keywords_type, keywords := range forbid_keywords {
		for k, v := range keywords {
			if k > 0 {
				buffer.WriteString(`|`)
			}
			switch keywords_type {
			case WITH_NO_SPACE:
				buffer.WriteString(utils.UpperTrim(v))
			case WITH_SPACES:
				buffer.WriteString(fmt.Sprintf(`\\s%v\\s`, utils.UpperTrim(v)))
			}
		}
	}
	buffer.WriteString(`.*`)
	//logger.Infof("regexp: %v", buffer.String())
	reg, err := regexp.Compile(buffer.String())
	if err != nil {
		logger.Error(err)
		return nil
	}

	return reg
}

func SqlForbid(sql string, reg *regexp.Regexp) bool {
	sql = utils.UpperTrim(sql)
	//logger.Infof("SQL:%v", sql)
	if !strings.HasPrefix(sql, "SELECT") {
		logger.Error("SQL should begin with SELECT.")
		return true
	}

	if reg != nil && reg.MatchString(sql) {
		logger.Error("SQL has not-allowed keywords.")
		return true
	}

	return false
}
