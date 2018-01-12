package engine

import (
	"github.com/XiBao/dbpool"
	"github.com/bububa/goconfig/config"
	"github.com/bububa/mymysql/autorc"
	_ "github.com/bububa/mymysql/thrsafe"
)

const (
	_CONFIG_FILE = "/var/code/go/config.cfg"
)

func get_default_mysql() (*autorc.Conn, error) {
	cfg, err := config.ReadDefault(_CONFIG_FILE)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	hostMaster, _ := cfg.String("masterdb", "host")
	userMaster, _ := cfg.String("masterdb", "user")
	passwdMaster, _ := cfg.String("masterdb", "passwd")
	dbnameMaster, _ := cfg.String("masterdb", "dbname")

	mysqlConfigMaster := &dbpool.MySQLConfig{
		Host:   hostMaster,
		User:   userMaster,
		Passwd: passwdMaster,
		DbName: dbnameMaster,
	}

	mdb := autorc.New("tcp", "", mysqlConfigMaster.Host, mysqlConfigMaster.User, mysqlConfigMaster.Passwd, mysqlConfigMaster.DbName)
	mdb.Register("set names utf8")

	return mdb, nil
}
