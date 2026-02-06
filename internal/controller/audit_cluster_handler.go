package controller

import (
	"strconv"

	"github.com/SisyphusSQ/golib/models/vo/base_vo"
	"github.com/labstack/echo/v4"

	"go-starter/internal/service/example_srv"
)

type AuditClusterController struct {
	AuditClusterService example_srv.AuditClusterService
}

func InitAuditClusterController(e *echo.Echo, clusterService example_srv.AuditClusterService) {
	controller := &AuditClusterController{
		AuditClusterService: clusterService,
	}
	g := e.Group("/auditCluster")
	g.GET("/:id", controller.GetByID)
}

// GetByID godoc
func (a *AuditClusterController) GetByID(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return base_vo.CommErrResp(c, err)
	}
	id := int64(idParam)
	data, err := a.AuditClusterService.GetByID(c.Request().Context(), id)
	if err != nil {
		return base_vo.CommErrResp(c, err)
	}
	return base_vo.CommSuccResp(c, data)
}
