package handler

import (
	"strconv"

	"epbms/internal/domain"
	"epbms/internal/middleware"
	"epbms/pkg/response"
	"github.com/gin-gonic/gin"
)

// PerformerHandler handles performer-related HTTP requests.
type PerformerHandler struct {
	performerSvc domain.PerformerService
}

// NewPerformerHandler creates a new PerformerHandler.
func NewPerformerHandler(performerSvc domain.PerformerService) *PerformerHandler {
	return &PerformerHandler{performerSvc: performerSvc}
}

// GetAll godoc
// GET /api/v1/performers
// Accessible by: all authenticated users
func (h *PerformerHandler) GetAll(c *gin.Context) {
	filter := domain.PerformerFilter{
		Category: c.Query("category"),
		Page:     parseIntQuery(c, "page", 1),
		PageSize: parseIntQuery(c, "page_size", 20),
	}

	performers, total, err := h.performerSvc.GetAll(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OKWithMeta(c, performers, response.Meta{
		Page:     filter.Page,
		PageSize: filter.PageSize,
		Total:    total,
	})
}

// GetByID godoc
// GET /api/v1/performers/:id
// Accessible by: all authenticated users
func (h *PerformerHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid performer id")
		return
	}

	performer, err := h.performerSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, performer)
}

// Create godoc
// POST /api/v1/performers
// Accessible by: admin only
func (h *PerformerHandler) Create(c *gin.Context) {
	var req domain.CreatePerformerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	callerID := middleware.GetCallerID(c)
	performer, err := h.performerSvc.Create(c.Request.Context(), callerID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, performer)
}

// Update godoc
// PUT /api/v1/performers/:id
// Accessible by: admin (any), performer (own profile only)
func (h *PerformerHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid performer id")
		return
	}

	var req domain.UpdatePerformerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	callerID := middleware.GetCallerID(c)
	callerRole := middleware.GetCallerRole(c)

	performer, err := h.performerSvc.Update(c.Request.Context(), id, callerID, callerRole, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, performer)
}

// Delete godoc
// DELETE /api/v1/performers/:id
// Accessible by: admin only
func (h *PerformerHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid performer id")
		return
	}

	if err := h.performerSvc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, gin.H{"message": "performer deleted successfully"})
}

// --- helpers ---

func parseIDParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 64)
	return uint(val), err
}

func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	raw := c.Query(key)
	if raw == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val < 1 {
		return defaultVal
	}
	return val
}
