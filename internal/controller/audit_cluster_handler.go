package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"go-starter/internal/models/resp"
	"go-starter/internal/service"
	"go-starter/utils"
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
// @Summary Get BusinessGroup By ID
// @Description Get BusinessGroup By ID
// @Tags BusinessGroup
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security ApiKeyAuth
// @Header 200 {string} Token "qwerty"
// @Success 200 {object} SimpleResponse{data=models.BusinessGroup} "BusinessGroup Info"
// @Failure 400,401,404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /business_groups/{id} [get]
func (a *AuditClusterController) GetByID(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), resp.ErrorResp(err))
	}
	id := int64(idParam)
	ac, err := a.AuditClusterService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), resp.ErrorResp(err))
	}
	return c.JSON(http.StatusOK, resp.SuccessResp(ac))
}
