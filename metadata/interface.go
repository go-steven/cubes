package metadata

import (
	log "github.com/kdar/factorlog"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
)

var (
	logger *log.FactorLog = utils.SetGlobalLogger("")
)

func SetLogger(l *log.FactorLog) {
	logger = l
}
