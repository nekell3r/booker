package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"booker/cmd/venue-svc/config"
	"booker/cmd/venue-svc/repository"
	"booker/pkg/kafka"
	bookingpb "booker/pkg/proto/booking"
	commonpb "booker/pkg/proto/common"
	venuepb "booker/pkg/proto/venue"
	"booker/pkg/tracing"
)

type Service struct {
	venuepb.UnimplementedVenueServiceServer
	repo          *repository.Repository
	producer      *kafka.Producer
	bookingClient bookingpb.BookingServiceClient
	cfg           *config.Config
}

func New(repo *repository.Repository, producer *kafka.Producer, bookingClient bookingpb.BookingServiceClient, cfg *config.Config) *Service {
	return &Service{
		repo:          repo,
		producer:      producer,
		bookingClient: bookingClient,
		cfg:           cfg,
	}
}

func (s *Service) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.Venue, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateVenue")
	defer span.End()

	log.Info().
		Str("name", req.Name).
		Str("timezone", req.Timezone).
		Str("address", req.Address).
		Msg("Creating venue")

	id, err := s.repo.CreateVenue(ctx, req.Name, req.Timezone, req.Address)
	if err != nil {
		log.Error().Err(err).
			Str("name", req.Name).
			Msg("Failed to create venue in database")
		return nil, err
	}

	venue, err := s.repo.GetVenue(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Str("venue_id", id).
			Msg("Failed to retrieve created venue")
		return nil, err
	}

	log.Info().
		Str("venue_id", id).
		Str("name", venue.Name).
		Msg("Venue created successfully")

	return toVenueProto(venue), nil
}

func (s *Service) GetVenue(ctx context.Context, req *venuepb.GetVenueRequest) (*venuepb.Venue, error) {
	venue, err := s.repo.GetVenue(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toVenueProto(venue), nil
}

func (s *Service) ListVenues(ctx context.Context, req *venuepb.ListVenuesRequest) (*venuepb.ListVenuesResponse, error) {
	venues, total, err := s.repo.ListVenues(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	protoVenues := make([]*venuepb.Venue, len(venues))
	for i, v := range venues {
		protoVenues[i] = toVenueProto(v)
	}

	return &venuepb.ListVenuesResponse{
		Venues: protoVenues,
		Total:  total,
	}, nil
}

func (s *Service) UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest) (*venuepb.Venue, error) {
	err := s.repo.UpdateVenue(ctx, req.Id, req.Name, req.Address)
	if err != nil {
		return nil, err
	}

	venue, err := s.repo.GetVenue(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return toVenueProto(venue), nil
}

func (s *Service) DeleteVenue(ctx context.Context, req *venuepb.DeleteVenueRequest) (*venuepb.DeleteVenueResponse, error) {
	log.Info().Str("venue_id", req.Id).Msg("Deleting venue")

	err := s.repo.DeleteVenue(ctx, req.Id)
	if err != nil {
		log.Error().Err(err).
			Str("venue_id", req.Id).
			Msg("Failed to delete venue")
		return nil, err
	}

	log.Info().Str("venue_id", req.Id).Msg("Venue deleted successfully")
	return &venuepb.DeleteVenueResponse{Success: true}, nil
}

func (s *Service) CreateRoom(ctx context.Context, req *venuepb.CreateRoomRequest) (*venuepb.Room, error) {
	id, err := s.repo.CreateRoom(ctx, req.VenueId, req.Name)
	if err != nil {
		return nil, err
	}

	room, err := s.repo.GetRoom(ctx, id)
	if err != nil {
		return nil, err
	}

	return toRoomProto(room), nil
}

func (s *Service) GetRoom(ctx context.Context, req *venuepb.GetRoomRequest) (*venuepb.Room, error) {
	room, err := s.repo.GetRoom(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toRoomProto(room), nil
}

func (s *Service) ListRooms(ctx context.Context, req *venuepb.ListRoomsRequest) (*venuepb.ListRoomsResponse, error) {
	rooms, total, err := s.repo.ListRooms(ctx, req.VenueId, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	protoRooms := make([]*venuepb.Room, len(rooms))
	for i, r := range rooms {
		protoRooms[i] = toRoomProto(r)
	}

	return &venuepb.ListRoomsResponse{
		Rooms: protoRooms,
		Total: total,
	}, nil
}

func (s *Service) UpdateRoom(ctx context.Context, req *venuepb.UpdateRoomRequest) (*venuepb.Room, error) {
	err := s.repo.UpdateRoom(ctx, req.Id, req.Name)
	if err != nil {
		return nil, err
	}

	room, err := s.repo.GetRoom(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return toRoomProto(room), nil
}

func (s *Service) DeleteRoom(ctx context.Context, req *venuepb.DeleteRoomRequest) (*venuepb.DeleteRoomResponse, error) {
	err := s.repo.DeleteRoom(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &venuepb.DeleteRoomResponse{Success: true}, nil
}

func (s *Service) CreateTable(ctx context.Context, req *venuepb.CreateTableRequest) (*venuepb.Table, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateTable")
	defer span.End()

	id, err := s.repo.CreateTable(ctx, req.RoomId, req.Name, req.Capacity, req.CanMerge, req.Zone)
	if err != nil {
		return nil, err
	}

	table, err := s.repo.GetTable(ctx, id)
	if err != nil {
		return nil, err
	}

	// Publish event
	room, _ := s.repo.GetRoom(ctx, req.RoomId)
	if room != nil {
		event := &commonpb.VenueEvent{
			VenueId: room.VenueID,
			Payload: &commonpb.VenueEvent_LayoutUpdated{
				LayoutUpdated: &commonpb.TableLayoutUpdated{
					RoomId:   req.RoomId,
					TableIds: []string{id},
				},
			},
		}
		if err := s.producer.PublishVenueEvent(ctx, "table.layout.updated", event); err != nil {
			log.Error().Err(err).Msg("Failed to publish layout updated event")
		}
	}

	return toTableProto(table), nil
}

func (s *Service) GetTable(ctx context.Context, req *venuepb.GetTableRequest) (*venuepb.Table, error) {
	table, err := s.repo.GetTable(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toTableProto(table), nil
}

func (s *Service) ListTables(ctx context.Context, req *venuepb.ListTablesRequest) (*venuepb.ListTablesResponse, error) {
	tables, total, err := s.repo.ListTables(ctx, req.RoomId, req.VenueId, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	protoTables := make([]*venuepb.Table, len(tables))
	for i, t := range tables {
		protoTables[i] = toTableProto(t)
	}

	return &venuepb.ListTablesResponse{
		Tables: protoTables,
		Total:  total,
	}, nil
}

func (s *Service) UpdateTable(ctx context.Context, req *venuepb.UpdateTableRequest) (*venuepb.Table, error) {
	ctx, span := tracing.StartSpan(ctx, "UpdateTable")
	defer span.End()

	err := s.repo.UpdateTable(ctx, req.Id, req.Name, req.Capacity, req.Zone)
	if err != nil {
		return nil, err
	}

	table, err := s.repo.GetTable(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Publish event
	room, _ := s.repo.GetRoom(ctx, table.RoomID)
	if room != nil {
		event := &commonpb.VenueEvent{
			VenueId: room.VenueID,
			Payload: &commonpb.VenueEvent_LayoutUpdated{
				LayoutUpdated: &commonpb.TableLayoutUpdated{
					RoomId:   table.RoomID,
					TableIds: []string{req.Id},
				},
			},
		}
		if err := s.producer.PublishVenueEvent(ctx, "table.layout.updated", event); err != nil {
			log.Error().Err(err).Msg("Failed to publish layout updated event")
		}
	}

	return toTableProto(table), nil
}

func (s *Service) DeleteTable(ctx context.Context, req *venuepb.DeleteTableRequest) (*venuepb.DeleteTableResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "DeleteTable")
	defer span.End()

	table, _ := s.repo.GetTable(ctx, req.Id)

	err := s.repo.DeleteTable(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Publish event
	if table != nil {
		room, _ := s.repo.GetRoom(ctx, table.RoomID)
		if room != nil {
			event := &commonpb.VenueEvent{
				VenueId: room.VenueID,
				Payload: &commonpb.VenueEvent_LayoutUpdated{
					LayoutUpdated: &commonpb.TableLayoutUpdated{
						RoomId:   table.RoomID,
						TableIds: []string{req.Id},
					},
				},
			}
			if err := s.producer.PublishVenueEvent(ctx, "table.layout.updated", event); err != nil {
				log.Error().Err(err).Msg("Failed to publish layout updated event")
			}
		}
	}

	return &venuepb.DeleteTableResponse{Success: true}, nil
}

func (s *Service) SetOpeningHours(ctx context.Context, req *venuepb.SetOpeningHoursRequest) (*venuepb.SetOpeningHoursResponse, error) {
	// TODO: Implement opening hours storage
	return &venuepb.SetOpeningHoursResponse{Success: true}, nil
}

func (s *Service) GetOpeningHours(ctx context.Context, req *venuepb.GetOpeningHoursRequest) (*venuepb.OpeningHours, error) {
	// TODO: Implement opening hours retrieval
	return &venuepb.OpeningHours{VenueId: req.VenueId}, nil
}

func (s *Service) SetSpecialHours(ctx context.Context, req *venuepb.SetSpecialHoursRequest) (*venuepb.SetSpecialHoursResponse, error) {
	// TODO: Implement special hours storage
	// Publish event
	event := &commonpb.VenueEvent{
		VenueId: req.VenueId,
		Payload: &commonpb.VenueEvent_ScheduleUpdated{
			ScheduleUpdated: &commonpb.VenueScheduleUpdated{
				Date: req.Date,
			},
		},
	}
	if err := s.producer.PublishVenueEvent(ctx, "venue.schedule.updated", event); err != nil {
		log.Error().Err(err).Msg("Failed to publish schedule updated event")
	}

	return &venuepb.SetSpecialHoursResponse{Success: true}, nil
}

func (s *Service) CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest) (*venuepb.CheckAvailabilityResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "CheckAvailability")
	defer span.End()

	// Get all tables in the venue
	allTables, _, err := s.repo.ListTables(ctx, "", req.VenueId, 1000, 0)
	if err != nil {
		return nil, err
	}

	if len(allTables) == 0 {
		return &venuepb.CheckAvailabilityResponse{
			Tables: []*venuepb.TableAvailability{},
		}, nil
	}

	// Collect table IDs
	tableIDs := make([]string, 0, len(allTables))
	tableMap := make(map[string]*repository.Table)
	for _, table := range allTables {
		tableIDs = append(tableIDs, table.ID)
		tableMap[table.ID] = table
	}

	// Check availability via booking service
	availabilityResp, err := s.bookingClient.CheckTableAvailability(ctx, &bookingpb.CheckTableAvailabilityRequest{
		VenueId:  req.VenueId,
		TableIds: tableIDs,
		Slot:     req.Slot,
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to check table availability, assuming all available")
		// Continue with all tables marked as available
	}

	// Build availability map
	availabilityMap := make(map[string]bool)
	for _, info := range availabilityResp.GetTables() {
		availabilityMap[info.TableId] = info.Available
	}

	// Find tables that can accommodate party size
	result := make([]*venuepb.TableAvailability, 0)

	// First, find single tables that can accommodate
	for _, table := range allTables {
		available := availabilityMap[table.ID]
		if !available {
			continue
		}

		if table.Capacity >= req.PartySize {
			// Get room info for table ref
			room, err := s.repo.GetRoom(ctx, table.RoomID)
			if err != nil {
				log.Warn().Err(err).Str("room_id", table.RoomID).Msg("Failed to get room")
				continue
			}

			result = append(result, &venuepb.TableAvailability{
				Table: &commonpb.TableRef{
					VenueId: req.VenueId,
					RoomId:  room.ID,
					TableId: table.ID,
				},
				Available: true,
				Reason:    "",
			})
		}
	}

	// If no single table can accommodate, look for mergeable tables
	if len(result) == 0 {
		// Find pairs of mergeable tables that together can accommodate
		// We need to add both directions (table1+table2 and table2+table1)
		// so that regardless of which table is selected, we can find the merge option
		for i, table1 := range allTables {
			if !table1.CanMerge || !availabilityMap[table1.ID] {
				continue
			}

			for j, table2 := range allTables {
				// Skip same table and ensure we only check each pair once
				if i >= j || !table2.CanMerge || !availabilityMap[table2.ID] {
					continue
				}

				// Check if they're in the same room (can be merged)
				if table1.RoomID != table2.RoomID {
					continue
				}

				totalCapacity := table1.Capacity + table2.Capacity
				if totalCapacity >= req.PartySize {
					room, err := s.repo.GetRoom(ctx, table1.RoomID)
					if err != nil {
						log.Warn().Err(err).Str("room_id", table1.RoomID).Msg("Failed to get room")
						continue
					}

					// Add both directions: (table1, table2) and (table2, table1)
					// This ensures that regardless of which table is selected,
					// we can find the merge option

					// Option 1: table1 as primary, table2 as merged
					result = append(result, &venuepb.TableAvailability{
						Table: &commonpb.TableRef{
							VenueId: req.VenueId,
							RoomId:  room.ID,
							TableId: table1.ID,
						},
						Available: true,
						Reason:    fmt.Sprintf("Can be merged with table %s (total capacity: %d)", table2.Name, totalCapacity),
						MergedWithTable: &commonpb.TableRef{
							VenueId: req.VenueId,
							RoomId:  room.ID,
							TableId: table2.ID,
						},
					})

					// Option 2: table2 as primary, table1 as merged
					result = append(result, &venuepb.TableAvailability{
						Table: &commonpb.TableRef{
							VenueId: req.VenueId,
							RoomId:  room.ID,
							TableId: table2.ID,
						},
						Available: true,
						Reason:    fmt.Sprintf("Can be merged with table %s (total capacity: %d)", table1.Name, totalCapacity),
						MergedWithTable: &commonpb.TableRef{
							VenueId: req.VenueId,
							RoomId:  room.ID,
							TableId: table1.ID,
						},
					})

					log.Info().
						Str("table1_id", table1.ID).
						Str("table2_id", table2.ID).
						Msg("Adding merge suggestions for both tables")
				}
			}
		}
	}

	return &venuepb.CheckAvailabilityResponse{
		Tables: result,
	}, nil
}

func (s *Service) GetTableLayout(ctx context.Context, req *venuepb.GetTableLayoutRequest) (*venuepb.GetTableLayoutResponse, error) {
	// TODO: Implement with Redis cache
	tables, _, err := s.repo.ListTables(ctx, req.RoomId, req.VenueId, 100, 0)
	if err != nil {
		return nil, err
	}

	protoTables := make([]*venuepb.Table, len(tables))
	for i, t := range tables {
		protoTables[i] = toTableProto(t)
	}

	return &venuepb.GetTableLayoutResponse{
		RoomId: req.RoomId,
		Tables: protoTables,
	}, nil
}

// Converters
func toVenueProto(v *repository.Venue) *venuepb.Venue {
	return &venuepb.Venue{
		Id:        v.ID,
		Name:      v.Name,
		Timezone:  v.Timezone,
		Address:   v.Address,
		CreatedAt: v.CreatedAt.Unix(),
		UpdatedAt: v.UpdatedAt.Unix(),
	}
}

func toRoomProto(r *repository.Room) *venuepb.Room {
	return &venuepb.Room{
		Id:        r.ID,
		VenueId:   r.VenueID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt.Unix(),
		UpdatedAt: r.UpdatedAt.Unix(),
	}
}

func toTableProto(t *repository.Table) *venuepb.Table {
	return &venuepb.Table{
		Id:        t.ID,
		RoomId:    t.RoomID,
		Name:      t.Name,
		Capacity:  t.Capacity,
		CanMerge:  t.CanMerge,
		Zone:      t.Zone,
		CreatedAt: t.CreatedAt.Unix(),
		UpdatedAt: t.UpdatedAt.Unix(),
	}
}
