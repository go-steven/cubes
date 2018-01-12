package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.xibao100.com/skyline/skyline/cubes/engine"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
	"net/http"
)

type CubesRptRequest struct {
	Tpl    string `form:"tpl" binding:"required"`
	TplCfg string `form:"tpl_cfg"`
}

func CubesRptHandler(c *gin.Context) {
	var request CubesRptRequest
	if err := c.BindWith(&request, binding.FormMultipart); err != nil {
		c.JSON(http.StatusOK, APIError{Code: BADREQUEST_ERROR, Msg: err.Error()})
		return
	}
	//logger.Infof("request = %v", Json(request))

	engine.SetLogger(logger)
	rptEngine := engine.NewReportEngine()
	defer rptEngine.Cleanup()

	storesLimit, err := metadata.NewStoresLimitFromStr(StoresLimitYaml, metadata.TPL_YAML)
	if err != nil {
		c.JSON(http.StatusOK, APIError{Code: BADREQUEST_ERROR, Msg: err.Error()})
		return
	}
	rptEngine.SetStoresLimit(storesLimit)

	rptResult, err := rptEngine.ExecuteTplConfig(request.Tpl, request.TplCfg, nil)
	if err != nil {
		c.JSON(http.StatusOK, APIError{Code: BADREQUEST_ERROR, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, rptResult)
}
