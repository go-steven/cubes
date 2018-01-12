package engine

import (
	log "github.com/kdar/factorlog"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
)

var (
	logger *log.FactorLog = utils.SetGlobalLogger("")
)

func SetLogger(l *log.FactorLog) {
	logger = l

	metadata.SetLogger(logger)
	source.SetLogger(logger)
}
