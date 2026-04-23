package handler

import (
	"epbms/internal/domain"
	"epbms/pkg/response"
	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin-specific HTTP requests.
type AdminHandler struct {
	userRepo    domain.UserRepository
	bookingRepo domain.BookingRepository
	perfRepo    domain.PerformerRepository
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(userRepo domain.UserRepository, bookingRepo domain.BookingRepository, perfRepo domain.PerformerRepository) *AdminHandler {
	return &AdminHandler{
		userRepo:    userRepo,
		bookingRepo: bookingRepo,
		perfRepo:    perfRepo,
	}
}

// DashboardStats godoc
// GET /api/v1/admin/stats
type dashboardStats struct {
	TotalBookings     int64 `json:"total_bookings"`
	PendingBookings   int64 `json:"pending_bookings"`
	ConfirmedBookings int64 `json:"confirmed_bookings"`
	RejectedBookings  int64 `json:"rejected_bookings"`
	CompletedBookings int64 `json:"completed_bookings"`
	TotalPerformers   int64 `json:"total_performers"`
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	countByStatus := func(status domain.BookingStatus) (int64, error) {
		_, total, err := h.bookingRepo.FindAll(ctx, domain.BookingFilter{Status: status, Page: 1, PageSize: 1})
		return total, err
	}

	total, _, err := h.bookingRepo.FindAll(ctx, domain.BookingFilter{Page: 1, PageSize: 1})
	if err != nil {
		response.Error(c, err)
		return
	}

	pending, err := countByStatus(domain.StatusPending)
	if err != nil {
		response.Error(c, err)
		return
	}
	confirmed, err := countByStatus(domain.StatusConfirmed)
	if err != nil {
		response.Error(c, err)
		return
	}
	rejected, err := countByStatus(domain.StatusRejected)
	if err != nil {
		response.Error(c, err)
		return
	}
	completed, err := countByStatus(domain.StatusCompleted)
	if err != nil {
		response.Error(c, err)
		return
	}

	_, perfTotal, err := h.perfRepo.FindAll(ctx, domain.PerformerFilter{Page: 1, PageSize: 1})
	if err != nil {
		response.Error(c, err)
		return
	}

	_ = total
	response.OK(c, dashboardStats{
		TotalBookings:     pending + confirmed + rejected + completed,
		PendingBookings:   pending,
		ConfirmedBookings: confirmed,
		RejectedBookings:  rejected,
		CompletedBookings: completed,
		TotalPerformers:   perfTotal,
	})
}
