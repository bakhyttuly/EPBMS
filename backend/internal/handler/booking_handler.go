package handler

import (
	"epbms/internal/domain"
	"epbms/internal/middleware"
	"epbms/pkg/response"
	"github.com/gin-gonic/gin"
)

// BookingHandler handles booking-related HTTP requests.
type BookingHandler struct {
	bookingSvc domain.BookingService
}

// NewBookingHandler creates a new BookingHandler.
func NewBookingHandler(bookingSvc domain.BookingService) *BookingHandler {
	return &BookingHandler{bookingSvc: bookingSvc}
}

// GetAll godoc
// GET /api/v1/bookings
// Role visibility:
//   - ADMIN: all bookings (filterable)
//   - CLIENT: own bookings only
//   - PERFORMER: own confirmed bookings only
func (h *BookingHandler) GetAll(c *gin.Context) {
	callerID := middleware.GetCallerID(c)
	callerRole := middleware.GetCallerRole(c)

	filter := domain.BookingFilter{
		EventDate: c.Query("event_date"),
		Page:      parseIntQuery(c, "page", 1),
		PageSize:  parseIntQuery(c, "page_size", 20),
	}

	// Admins may additionally filter by status and performer.
	if callerRole == domain.RoleAdmin {
		filter.Status = domain.BookingStatus(c.Query("status"))
		if pid, err := parseIDParam(c, "performer_id"); err == nil && pid != 0 {
			filter.PerformerID = pid
		}
	}

	bookings, total, err := h.bookingSvc.GetAll(c.Request.Context(), callerID, callerRole, filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OKWithMeta(c, bookings, response.Meta{
		Page:     filter.Page,
		PageSize: filter.PageSize,
		Total:    total,
	})
}

// GetByID godoc
// GET /api/v1/bookings/:id
func (h *BookingHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid booking id")
		return
	}

	callerID := middleware.GetCallerID(c)
	callerRole := middleware.GetCallerRole(c)

	booking, err := h.bookingSvc.GetByID(c.Request.Context(), id, callerID, callerRole)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, booking)
}

// Create godoc
// POST /api/v1/bookings
// Accessible by: CLIENT only — creates a pending booking request.
func (h *BookingHandler) Create(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	clientID := middleware.GetCallerID(c)

	booking, err := h.bookingSvc.CreateRequest(c.Request.Context(), clientID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, booking)
}

// UpdateStatus godoc
// PUT /api/v1/admin/bookings/:id/status
// Accessible by: ADMIN only — approves, rejects, or marks a booking as completed.
func (h *BookingHandler) UpdateStatus(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid booking id")
		return
	}

	var req domain.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := middleware.GetCallerID(c)

	booking, err := h.bookingSvc.UpdateStatus(c.Request.Context(), id, adminID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, booking)
}

// Delete godoc
// DELETE /api/v1/admin/bookings/:id
// Accessible by: ADMIN only.
func (h *BookingHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid booking id")
		return
	}

	if err := h.bookingSvc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, gin.H{"message": "booking deleted successfully"})
}
