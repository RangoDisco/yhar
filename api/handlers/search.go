package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/api/common"
	"github.com/rangodisco/yhar/api/dto/request"
	"github.com/rangodisco/yhar/api/services"
)

type SearchHandler struct {
	service *services.SearchService
}

func NewSearchHandler(service *services.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) Search(c *gin.Context) {
	var query request.SearchQuery
	ctx := c.Request.Context()

	err := c.BindQuery(&query)
	if err != nil {
		common.RespondWithError(c, 400, err, "Invalid search query")
		return
	}

	res, totalCount, err := h.service.Search(ctx, query)
	if err != nil {
		common.RespondWithError(c, 500, err, "Unable to perform search")
		return
	}

	paginated := common.BuildPaginatedResponse(res, 1, query.Pagination.Limit, totalCount)

	common.RespondWithData(c, 200, paginated)
}
