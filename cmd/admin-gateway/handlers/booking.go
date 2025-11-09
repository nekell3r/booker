package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	bookingpb "booker/pkg/proto/booking"
	commonpb "booker/pkg/proto/common"
	venuepb "booker/pkg/proto/venue"
)

// Auth handlers
func (h *Handler) Login(c echo.Context) error {
	// TODO: Implement JWT login
	return c.JSON(http.StatusOK, map[string]string{"token": "dummy-token"})
}

func (h *Handler) RefreshToken(c echo.Context) error {
	// TODO: Implement token refresh
	return c.JSON(http.StatusOK, map[string]string{"token": "dummy-token"})
}

// Venue handlers
func (h *Handler) ListVenues(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.venueClient.ListVenues(c.Request().Context(), &venuepb.ListVenuesRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetVenue(c echo.Context) error {
	resp, err := h.venueClient.GetVenue(c.Request().Context(), &venuepb.GetVenueRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateVenue(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Timezone string `json:"timezone"`
		Address  string `json:"address"`
	}
	if err := c.Bind(&req); err != nil {
		log.Warn().Err(err).Msg("Failed to bind CreateVenue request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	log.Info().
		Str("name", req.Name).
		Str("timezone", req.Timezone).
		Str("address", req.Address).
		Msg("Creating venue")

	resp, err := h.venueClient.CreateVenue(c.Request().Context(), &venuepb.CreateVenueRequest{
		Name:     req.Name,
		Timezone: req.Timezone,
		Address:  req.Address,
	})
	if err != nil {
		log.Error().Err(err).
			Str("name", req.Name).
			Msg("Failed to create venue")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Info().
		Str("venue_id", resp.Id).
		Str("name", resp.Name).
		Msg("Venue created successfully")

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateVenue(c echo.Context) error {
	var req struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueClient.UpdateVenue(c.Request().Context(), &venuepb.UpdateVenueRequest{
		Id:      c.Param("id"),
		Name:    req.Name,
		Address: req.Address,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteVenue(c echo.Context) error {
	venueID := c.Param("id")
	
	log.Info().Str("venue_id", venueID).Msg("Deleting venue")
	
	_, err := h.venueClient.DeleteVenue(c.Request().Context(), &venuepb.DeleteVenueRequest{
		Id: venueID,
	})
	if err != nil {
		log.Error().Err(err).
			Str("venue_id", venueID).
			Msg("Failed to delete venue")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Info().Str("venue_id", venueID).Msg("Venue deleted successfully")
	return c.NoContent(http.StatusNoContent)
}

// Room handlers
func (h *Handler) ListRooms(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.venueClient.ListRooms(c.Request().Context(), &venuepb.ListRoomsRequest{
		VenueId: c.Param("venueId"),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetRoom(c echo.Context) error {
	resp, err := h.venueClient.GetRoom(c.Request().Context(), &venuepb.GetRoomRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateRoom(c echo.Context) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueClient.CreateRoom(c.Request().Context(), &venuepb.CreateRoomRequest{
		VenueId: c.Param("venueId"),
		Name:    req.Name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateRoom(c echo.Context) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueClient.UpdateRoom(c.Request().Context(), &venuepb.UpdateRoomRequest{
		Id:   c.Param("id"),
		Name: req.Name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteRoom(c echo.Context) error {
	_, err := h.venueClient.DeleteRoom(c.Request().Context(), &venuepb.DeleteRoomRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Table handlers
func (h *Handler) ListTables(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.venueClient.ListTables(c.Request().Context(), &venuepb.ListTablesRequest{
		RoomId: c.Param("roomId"),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetTable(c echo.Context) error {
	resp, err := h.venueClient.GetTable(c.Request().Context(), &venuepb.GetTableRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateTable(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Capacity int32  `json:"capacity"`
		CanMerge bool   `json:"can_merge"`
		Zone     string `json:"zone"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueClient.CreateTable(c.Request().Context(), &venuepb.CreateTableRequest{
		RoomId:   c.Param("roomId"),
		Name:     req.Name,
		Capacity: req.Capacity,
		CanMerge: req.CanMerge,
		Zone:     req.Zone,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateTable(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Capacity int32  `json:"capacity"`
		Zone     string `json:"zone"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueClient.UpdateTable(c.Request().Context(), &venuepb.UpdateTableRequest{
		Id:       c.Param("id"),
		Name:     req.Name,
		Capacity: req.Capacity,
		Zone:     req.Zone,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteTable(c echo.Context) error {
	_, err := h.venueClient.DeleteTable(c.Request().Context(), &venuepb.DeleteTableRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Schedule handlers
func (h *Handler) GetOpeningHours(c echo.Context) error {
	resp, err := h.venueClient.GetOpeningHours(c.Request().Context(), &venuepb.GetOpeningHoursRequest{
		VenueId: c.Param("venueId"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) SetOpeningHours(c echo.Context) error {
	var req struct {
		Days []struct {
			Weekday   int32  `json:"weekday"`
			OpenTime  string `json:"open_time"`
			CloseTime string `json:"close_time"`
		} `json:"days"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	days := make([]*venuepb.DayHours, len(req.Days))
	for i, d := range req.Days {
		days[i] = &venuepb.DayHours{
			Weekday:   d.Weekday,
			OpenTime:  d.OpenTime,
			CloseTime: d.CloseTime,
		}
	}

	resp, err := h.venueClient.SetOpeningHours(c.Request().Context(), &venuepb.SetOpeningHoursRequest{
		VenueId: c.Param("venueId"),
		Days:    days,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) SetSpecialHours(c echo.Context) error {
	var req struct {
		Date      string `json:"date"`
		OpenTime  string `json:"open_time"`
		CloseTime string `json:"close_time"`
		IsClosed  bool   `json:"is_closed"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueClient.SetSpecialHours(c.Request().Context(), &venuepb.SetSpecialHoursRequest{
		VenueId:   c.Param("venueId"),
		Date:      req.Date,
		OpenTime:  req.OpenTime,
		CloseTime: req.CloseTime,
		IsClosed:  req.IsClosed,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// Booking handlers
func (h *Handler) ListBookings(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.bookingClient.ListBookings(c.Request().Context(), &bookingpb.ListBookingsRequest{
		VenueId: c.QueryParam("venue_id"),
		Date:    c.QueryParam("date"),
		Status:  c.QueryParam("status"),
		TableId: c.QueryParam("table_id"),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetBooking(c echo.Context) error {
	resp, err := h.bookingClient.GetBooking(c.Request().Context(), &bookingpb.GetBookingRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateBooking(c echo.Context) error {
	var req struct {
		VenueID string `json:"venue_id"`
		Table   struct {
			VenueID string `json:"venue_id"`
			RoomID  string `json:"room_id"`
			TableID string `json:"table_id"`
		} `json:"table"`
		Slot struct {
			Date            string `json:"date"`
			StartTime       string `json:"start_time"`
			DurationMinutes int32  `json:"duration_minutes"`
		} `json:"slot"`
		PartySize      int32  `json:"party_size"`
		CustomerName   string `json:"customer_name"`
		CustomerPhone  string `json:"customer_phone"`
		Comment        string `json:"comment"`
		IdempotencyKey string `json:"idempotency_key"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingClient.CreateBooking(c.Request().Context(), &bookingpb.CreateBookingRequest{
		VenueId: req.VenueID,
		Table: &commonpb.TableRef{
			VenueId: req.Table.VenueID,
			RoomId:  req.Table.RoomID,
			TableId: req.Table.TableID,
		},
		Slot: &commonpb.Slot{
			Date:            req.Slot.Date,
			StartTime:       req.Slot.StartTime,
			DurationMinutes: req.Slot.DurationMinutes,
		},
		PartySize:      req.PartySize,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		Comment:        req.Comment,
		AdminId:        adminID,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) ConfirmBooking(c echo.Context) error {
	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingClient.ConfirmBooking(c.Request().Context(), &bookingpb.ConfirmBookingRequest{
		Id:      c.Param("id"),
		AdminId: adminID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CancelBooking(c echo.Context) error {
	var req struct {
		Reason string `json:"reason"`
	}
	c.Bind(&req)

	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingClient.CancelBooking(c.Request().Context(), &bookingpb.CancelBookingRequest{
		Id:      c.Param("id"),
		AdminId: adminID,
		Reason:  req.Reason,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) MarkSeated(c echo.Context) error {
	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingClient.MarkSeated(c.Request().Context(), &bookingpb.MarkSeatedRequest{
		Id:      c.Param("id"),
		AdminId: adminID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) MarkFinished(c echo.Context) error {
	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingClient.MarkFinished(c.Request().Context(), &bookingpb.MarkFinishedRequest{
		Id:      c.Param("id"),
		AdminId: adminID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) MarkNoShow(c echo.Context) error {
	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingClient.MarkNoShow(c.Request().Context(), &bookingpb.MarkNoShowRequest{
		Id:      c.Param("id"),
		AdminId: adminID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CheckAvailability(c echo.Context) error {
	var req struct {
		VenueID string `json:"venue_id"`
		Slot    struct {
			Date            string `json:"date"`
			StartTime       string `json:"start_time"`
			DurationMinutes int32  `json:"duration_minutes"`
		} `json:"slot"`
		PartySize int32 `json:"party_size"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueClient.CheckAvailability(c.Request().Context(), &venuepb.CheckAvailabilityRequest{
		VenueId: req.VenueID,
		Slot: &commonpb.Slot{
			Date:            req.Slot.Date,
			StartTime:       req.Slot.StartTime,
			DurationMinutes: req.Slot.DurationMinutes,
		},
		PartySize: req.PartySize,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) WebSocket(c echo.Context) error {
	// TODO: Implement WebSocket for live updates
	return c.String(http.StatusNotImplemented, "WebSocket not implemented yet")
}

// Metrics endpoint
func (h *Handler) Metrics(c echo.Context) error {
	promhttp.Handler().ServeHTTP(c.Response(), c.Request())
	return nil
}
