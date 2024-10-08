package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"go-starter/internal/models/resp"
	"go-starter/internal/service"
)

type AuditClusterController struct {
	AuditClusterService service.AuditClusterService
}

func InitAuditClusterController(e *echo.Echo, clusterService service.AuditClusterService) {
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
		return resp.CommErrResp(c, err)
	}
	id := int64(idParam)
	data, err := a.AuditClusterService.GetByID(c.Request().Context(), id)
	if err != nil {
		return resp.CommErrResp(c, err)
	}
	return resp.CommSuccResp(c, data)
}
