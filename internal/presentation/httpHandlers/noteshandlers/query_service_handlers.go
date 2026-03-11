package noteshandlers

import (
	"net/http"
	"notes_backend/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueryServiceHandler struct {
	repo repository.INoteQueryService
}

func NewQueryServiceHandler(repo repository.INoteQueryService) *QueryServiceHandler {
	return &QueryServiceHandler{repo}
}

// GetAllPublic godoc
// @Summary      GetAllPublic
// @Description  Получение всех публичных заметок с постраничной пагинацией
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        page  path      int  true  "Номер страницы (начиная с 1)"
// @Success      200   {object}  GetAllPublicResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/notes/public/{page} [get]
func (h *QueryServiceHandler) GetAllPublic(c *gin.Context) {
	ctx := c.Request.Context()

	pageParam := c.Param("page")
	pageNum, err := strconv.ParseInt(pageParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notes, err := h.repo.GetAllPublic(ctx, pageNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := GetAllPublicResponse{
		Data: make([]repository.GetAllPublicBO, len(*notes)),
	}

	copy(resp.Data, *notes)

	c.JSON(http.StatusOK, resp)
}

type GetAllPublicResponse struct {
	Data []repository.GetAllPublicBO `json:"data"`
}
