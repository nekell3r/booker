package testutil

import (
	"context"
	"time"

	"google.golang.org/grpc"

	bookingrepo "booker/cmd/booking-svc/repository"
	venuerepo "booker/cmd/venue-svc/repository"
	bookingpb "booker/pkg/proto/booking"
	commonpb "booker/pkg/proto/common"
	venuepb "booker/pkg/proto/venue"
)

// MockBookingRepository is a mock implementation of booking repository
type MockBookingRepository struct {
	CreateBookingFunc          func(ctx context.Context, booking *bookingrepo.Booking) error
	GetBookingFunc             func(ctx context.Context, id string) (*bookingrepo.Booking, error)
	ListBookingsFunc           func(ctx context.Context, filters *bookingrepo.BookingFilters) ([]*bookingrepo.Booking, int32, error)
	UpdateBookingStatusFunc    func(ctx context.Context, id, status string) error
	GetExpiredHoldsFunc        func(ctx context.Context) ([]*bookingrepo.Booking, error)
	CheckTableAvailabilityFunc func(ctx context.Context, venueID string, tableIDs []string, date, startTime, endTime string) (map[string]bool, error)
	AddToOutboxFunc            func(ctx context.Context, topic, key string, payload []byte) error
	GetPendingOutboxFunc       func(ctx context.Context, limit int32) ([]*bookingrepo.OutboxMessage, error)
	UpdateOutboxStatusFunc     func(ctx context.Context, id, status string, retryCount int32) error
}

func (m *MockBookingRepository) CreateBooking(ctx context.Context, booking *bookingrepo.Booking) error {
	if m.CreateBookingFunc != nil {
		return m.CreateBookingFunc(ctx, booking)
	}
	return nil
}

func (m *MockBookingRepository) GetBooking(ctx context.Context, id string) (*bookingrepo.Booking, error) {
	if m.GetBookingFunc != nil {
		return m.GetBookingFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBookingRepository) ListBookings(ctx context.Context, filters *bookingrepo.BookingFilters) ([]*bookingrepo.Booking, int32, error) {
	if m.ListBookingsFunc != nil {
		return m.ListBookingsFunc(ctx, filters)
	}
	return []*bookingrepo.Booking{}, 0, nil
}

func (m *MockBookingRepository) UpdateBookingStatus(ctx context.Context, id, status string) error {
	if m.UpdateBookingStatusFunc != nil {
		return m.UpdateBookingStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *MockBookingRepository) GetExpiredHolds(ctx context.Context) ([]*bookingrepo.Booking, error) {
	if m.GetExpiredHoldsFunc != nil {
		return m.GetExpiredHoldsFunc(ctx)
	}
	return []*bookingrepo.Booking{}, nil
}

func (m *MockBookingRepository) CheckTableAvailability(ctx context.Context, venueID string, tableIDs []string, date, startTime, endTime string) (map[string]bool, error) {
	if m.CheckTableAvailabilityFunc != nil {
		return m.CheckTableAvailabilityFunc(ctx, venueID, tableIDs, date, startTime, endTime)
	}
	return make(map[string]bool), nil
}

func (m *MockBookingRepository) AddToOutbox(ctx context.Context, topic, key string, payload []byte) error {
	if m.AddToOutboxFunc != nil {
		return m.AddToOutboxFunc(ctx, topic, key, payload)
	}
	return nil
}

func (m *MockBookingRepository) GetPendingOutbox(ctx context.Context, limit int32) ([]*bookingrepo.OutboxMessage, error) {
	if m.GetPendingOutboxFunc != nil {
		return m.GetPendingOutboxFunc(ctx, limit)
	}
	return []*bookingrepo.OutboxMessage{}, nil
}

func (m *MockBookingRepository) UpdateOutboxStatus(ctx context.Context, id, status string, retryCount int32) error {
	if m.UpdateOutboxStatusFunc != nil {
		return m.UpdateOutboxStatusFunc(ctx, id, status, retryCount)
	}
	return nil
}

// MockVenueRepository is a mock implementation of venue repository
type MockVenueRepository struct {
	CreateVenueFunc  func(ctx context.Context, name, timezone, address string) (string, error)
	GetVenueFunc     func(ctx context.Context, id string) (*venuerepo.Venue, error)
	ListVenuesFunc   func(ctx context.Context, limit, offset int32) ([]*venuerepo.Venue, int32, error)
	UpdateVenueFunc  func(ctx context.Context, id, name, address string) error
	DeleteVenueFunc  func(ctx context.Context, id string) error
	CreateRoomFunc   func(ctx context.Context, venueID, name string) (string, error)
	GetRoomFunc      func(ctx context.Context, id string) (*venuerepo.Room, error)
	ListRoomsFunc    func(ctx context.Context, venueID string, limit, offset int32) ([]*venuerepo.Room, int32, error)
	UpdateRoomFunc   func(ctx context.Context, id, name string) error
	DeleteRoomFunc   func(ctx context.Context, id string) error
	CreateTableFunc  func(ctx context.Context, roomID, name string, capacity int32, canMerge bool, zone string) (string, error)
	GetTableFunc     func(ctx context.Context, id string) (*venuerepo.Table, error)
	ListTablesFunc   func(ctx context.Context, roomID, venueID string, limit, offset int32) ([]*venuerepo.Table, int32, error)
	UpdateTableFunc  func(ctx context.Context, id, name string, capacity int32, zone string) error
	DeleteTableFunc  func(ctx context.Context, id string) error
	GetAllTablesFunc func(ctx context.Context, venueID string) ([]*venuerepo.Table, error)
}

func (m *MockVenueRepository) CreateVenue(ctx context.Context, name, timezone, address string) (string, error) {
	if m.CreateVenueFunc != nil {
		return m.CreateVenueFunc(ctx, name, timezone, address)
	}
	return "venue-1", nil
}

func (m *MockVenueRepository) GetVenue(ctx context.Context, id string) (*venuerepo.Venue, error) {
	if m.GetVenueFunc != nil {
		return m.GetVenueFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockVenueRepository) ListVenues(ctx context.Context, limit, offset int32) ([]*venuerepo.Venue, int32, error) {
	if m.ListVenuesFunc != nil {
		return m.ListVenuesFunc(ctx, limit, offset)
	}
	return []*venuerepo.Venue{}, 0, nil
}

func (m *MockVenueRepository) UpdateVenue(ctx context.Context, id, name, address string) error {
	if m.UpdateVenueFunc != nil {
		return m.UpdateVenueFunc(ctx, id, name, address)
	}
	return nil
}

func (m *MockVenueRepository) DeleteVenue(ctx context.Context, id string) error {
	if m.DeleteVenueFunc != nil {
		return m.DeleteVenueFunc(ctx, id)
	}
	return nil
}

func (m *MockVenueRepository) CreateRoom(ctx context.Context, venueID, name string) (string, error) {
	if m.CreateRoomFunc != nil {
		return m.CreateRoomFunc(ctx, venueID, name)
	}
	return "room-1", nil
}

func (m *MockVenueRepository) GetRoom(ctx context.Context, id string) (*venuerepo.Room, error) {
	if m.GetRoomFunc != nil {
		return m.GetRoomFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockVenueRepository) ListRooms(ctx context.Context, venueID string, limit, offset int32) ([]*venuerepo.Room, int32, error) {
	if m.ListRoomsFunc != nil {
		return m.ListRoomsFunc(ctx, venueID, limit, offset)
	}
	return []*venuerepo.Room{}, 0, nil
}

func (m *MockVenueRepository) UpdateRoom(ctx context.Context, id, name string) error {
	if m.UpdateRoomFunc != nil {
		return m.UpdateRoomFunc(ctx, id, name)
	}
	return nil
}

func (m *MockVenueRepository) DeleteRoom(ctx context.Context, id string) error {
	if m.DeleteRoomFunc != nil {
		return m.DeleteRoomFunc(ctx, id)
	}
	return nil
}

func (m *MockVenueRepository) CreateTable(ctx context.Context, roomID, name string, capacity int32, canMerge bool, zone string) (string, error) {
	if m.CreateTableFunc != nil {
		return m.CreateTableFunc(ctx, roomID, name, capacity, canMerge, zone)
	}
	return "table-1", nil
}

func (m *MockVenueRepository) GetTable(ctx context.Context, id string) (*venuerepo.Table, error) {
	if m.GetTableFunc != nil {
		return m.GetTableFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockVenueRepository) ListTables(ctx context.Context, roomID, venueID string, limit, offset int32) ([]*venuerepo.Table, int32, error) {
	if m.ListTablesFunc != nil {
		return m.ListTablesFunc(ctx, roomID, venueID, limit, offset)
	}
	return []*venuerepo.Table{}, 0, nil
}

func (m *MockVenueRepository) UpdateTable(ctx context.Context, id, name string, capacity int32, zone string) error {
	if m.UpdateTableFunc != nil {
		return m.UpdateTableFunc(ctx, id, name, capacity, zone)
	}
	return nil
}

func (m *MockVenueRepository) DeleteTable(ctx context.Context, id string) error {
	if m.DeleteTableFunc != nil {
		return m.DeleteTableFunc(ctx, id)
	}
	return nil
}

// MockKafkaProducer is a mock implementation of kafka producer
type MockKafkaProducer struct {
	PublishBookingEventFunc func(ctx context.Context, topic string, event *commonpb.BookingEvent) error
	PublishVenueEventFunc   func(ctx context.Context, topic string, event *commonpb.VenueEvent) error
	CloseFunc               func() error
}

func (m *MockKafkaProducer) PublishBookingEvent(ctx context.Context, topic string, event *commonpb.BookingEvent) error {
	if m.PublishBookingEventFunc != nil {
		return m.PublishBookingEventFunc(ctx, topic, event)
	}
	return nil
}

func (m *MockKafkaProducer) PublishVenueEvent(ctx context.Context, topic string, event *commonpb.VenueEvent) error {
	if m.PublishVenueEventFunc != nil {
		return m.PublishVenueEventFunc(ctx, topic, event)
	}
	return nil
}

func (m *MockKafkaProducer) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// MockRedisClient is a mock implementation of redis client
type MockRedisClient struct {
	SetHoldFunc    func(ctx context.Context, key string, bookingID string, ttl time.Duration) (bool, error)
	GetHoldFunc    func(ctx context.Context, key string) (string, error)
	DeleteHoldFunc func(ctx context.Context, key string) error
	IncrFunc       func(ctx context.Context, key string) (int64, error)
	ExpireFunc     func(ctx context.Context, key string, expiration time.Duration) error
	DelFunc        func(ctx context.Context, key string) error
}

func (m *MockRedisClient) SetHold(ctx context.Context, key string, bookingID string, ttl time.Duration) (bool, error) {
	if m.SetHoldFunc != nil {
		return m.SetHoldFunc(ctx, key, bookingID, ttl)
	}
	return true, nil
}

func (m *MockRedisClient) GetHold(ctx context.Context, key string) (string, error) {
	if m.GetHoldFunc != nil {
		return m.GetHoldFunc(ctx, key)
	}
	return "", nil
}

func (m *MockRedisClient) DeleteHold(ctx context.Context, key string) error {
	if m.DeleteHoldFunc != nil {
		return m.DeleteHoldFunc(ctx, key)
	}
	return nil
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	if m.IncrFunc != nil {
		return m.IncrFunc(ctx, key)
	}
	return 0, nil
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if m.ExpireFunc != nil {
		return m.ExpireFunc(ctx, key, expiration)
	}
	return nil
}

func (m *MockRedisClient) Del(ctx context.Context, key string) error {
	if m.DelFunc != nil {
		return m.DelFunc(ctx, key)
	}
	return nil
}

// MockVenueServiceClient is a mock implementation of venue gRPC client
// Note: This is a simplified mock. In production, you'd use a proper mocking library
// or generate mocks from proto interfaces
type MockVenueServiceClient struct {
	CheckAvailabilityFunc func(ctx context.Context, req *venuepb.CheckAvailabilityRequest, opts ...grpc.CallOption) (*venuepb.CheckAvailabilityResponse, error)
	ListVenuesFunc        func(ctx context.Context, req *venuepb.ListVenuesRequest, opts ...grpc.CallOption) (*venuepb.ListVenuesResponse, error)
	GetVenueFunc          func(ctx context.Context, req *venuepb.GetVenueRequest, opts ...grpc.CallOption) (*venuepb.Venue, error)
	CreateVenueFunc       func(ctx context.Context, req *venuepb.CreateVenueRequest, opts ...grpc.CallOption) (*venuepb.Venue, error)
	UpdateVenueFunc       func(ctx context.Context, req *venuepb.UpdateVenueRequest, opts ...grpc.CallOption) (*venuepb.Venue, error)
	DeleteVenueFunc       func(ctx context.Context, req *venuepb.DeleteVenueRequest, opts ...grpc.CallOption) (*venuepb.DeleteVenueResponse, error)
}

func (m *MockVenueServiceClient) CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest, opts ...grpc.CallOption) (*venuepb.CheckAvailabilityResponse, error) {
	if m.CheckAvailabilityFunc != nil {
		return m.CheckAvailabilityFunc(ctx, req, opts...)
	}
	return &venuepb.CheckAvailabilityResponse{}, nil
}

func (m *MockVenueServiceClient) ListVenues(ctx context.Context, req *venuepb.ListVenuesRequest, opts ...grpc.CallOption) (*venuepb.ListVenuesResponse, error) {
	if m.ListVenuesFunc != nil {
		return m.ListVenuesFunc(ctx, req, opts...)
	}
	return &venuepb.ListVenuesResponse{}, nil
}

func (m *MockVenueServiceClient) GetVenue(ctx context.Context, req *venuepb.GetVenueRequest, opts ...grpc.CallOption) (*venuepb.Venue, error) {
	if m.GetVenueFunc != nil {
		return m.GetVenueFunc(ctx, req, opts...)
	}
	return &venuepb.Venue{}, nil
}

func (m *MockVenueServiceClient) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest, opts ...grpc.CallOption) (*venuepb.Venue, error) {
	if m.CreateVenueFunc != nil {
		return m.CreateVenueFunc(ctx, req, opts...)
	}
	return &venuepb.Venue{}, nil
}

func (m *MockVenueServiceClient) UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest, opts ...grpc.CallOption) (*venuepb.Venue, error) {
	if m.UpdateVenueFunc != nil {
		return m.UpdateVenueFunc(ctx, req, opts...)
	}
	return &venuepb.Venue{}, nil
}

func (m *MockVenueServiceClient) DeleteVenue(ctx context.Context, req *venuepb.DeleteVenueRequest, opts ...grpc.CallOption) (*venuepb.DeleteVenueResponse, error) {
	if m.DeleteVenueFunc != nil {
		return m.DeleteVenueFunc(ctx, req, opts...)
	}
	return &venuepb.DeleteVenueResponse{}, nil
}

// MockBookingServiceClient is a mock implementation of booking gRPC client
type MockBookingServiceClient struct {
	CheckTableAvailabilityFunc func(ctx context.Context, req *bookingpb.CheckTableAvailabilityRequest, opts ...grpc.CallOption) (*bookingpb.CheckTableAvailabilityResponse, error)
	CreateBookingFunc          func(ctx context.Context, req *bookingpb.CreateBookingRequest, opts ...grpc.CallOption) (*bookingpb.Booking, error)
	GetBookingFunc             func(ctx context.Context, req *bookingpb.GetBookingRequest, opts ...grpc.CallOption) (*bookingpb.Booking, error)
	ListBookingsFunc           func(ctx context.Context, req *bookingpb.ListBookingsRequest, opts ...grpc.CallOption) (*bookingpb.ListBookingsResponse, error)
	ConfirmBookingFunc         func(ctx context.Context, req *bookingpb.ConfirmBookingRequest, opts ...grpc.CallOption) (*bookingpb.Booking, error)
	CancelBookingFunc          func(ctx context.Context, req *bookingpb.CancelBookingRequest, opts ...grpc.CallOption) (*bookingpb.Booking, error)
}

func (m *MockBookingServiceClient) CheckTableAvailability(ctx context.Context, req *bookingpb.CheckTableAvailabilityRequest, opts ...grpc.CallOption) (*bookingpb.CheckTableAvailabilityResponse, error) {
	if m.CheckTableAvailabilityFunc != nil {
		return m.CheckTableAvailabilityFunc(ctx, req, opts...)
	}
	return &bookingpb.CheckTableAvailabilityResponse{}, nil
}

func (m *MockBookingServiceClient) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest, opts ...grpc.CallOption) (*bookingpb.Booking, error) {
	if m.CreateBookingFunc != nil {
		return m.CreateBookingFunc(ctx, req, opts...)
	}
	return &bookingpb.Booking{}, nil
}

func (m *MockBookingServiceClient) GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest, opts ...grpc.CallOption) (*bookingpb.Booking, error) {
	if m.GetBookingFunc != nil {
		return m.GetBookingFunc(ctx, req, opts...)
	}
	return &bookingpb.Booking{}, nil
}

func (m *MockBookingServiceClient) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest, opts ...grpc.CallOption) (*bookingpb.ListBookingsResponse, error) {
	if m.ListBookingsFunc != nil {
		return m.ListBookingsFunc(ctx, req, opts...)
	}
	return &bookingpb.ListBookingsResponse{}, nil
}

func (m *MockBookingServiceClient) ConfirmBooking(ctx context.Context, req *bookingpb.ConfirmBookingRequest, opts ...grpc.CallOption) (*bookingpb.Booking, error) {
	if m.ConfirmBookingFunc != nil {
		return m.ConfirmBookingFunc(ctx, req, opts...)
	}
	return &bookingpb.Booking{}, nil
}

func (m *MockBookingServiceClient) CancelBooking(ctx context.Context, req *bookingpb.CancelBookingRequest, opts ...grpc.CallOption) (*bookingpb.Booking, error) {
	if m.CancelBookingFunc != nil {
		return m.CancelBookingFunc(ctx, req, opts...)
	}
	return &bookingpb.Booking{}, nil
}
