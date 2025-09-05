package services

import (
	"database/sql"

	"smithsolutions/go-api/internal/filters"
	"smithsolutions/go-api/internal/models"
	"smithsolutions/go-api/internal/util"
)

type CreateEvent struct {
	OwnerUserId int

	Label          string
	CoverPhotoPath *string
}

// func (c CreateEvent) SQL() ([]string, []any, error) {
// 	return util.GetCreateSQL(c)
// }

type UpdateEvent struct {
	Label          *string
	CoverPhotoPath *string
}

func (u UpdateEvent) SQL() (string, []any, error) {
	return util.GetUpdateSQL(u)
}

type WhereEvent struct {
	OwnerUserId *filters.IntFilter
}

func (w WhereEvent) SQL() (string, []any, error) {
	return util.GetWhereSQL(w)
}

type IncludeWithEvent struct {
	User bool
}

type EventService struct {
	ResourceService[models.Event, CreateEvent, UpdateEvent, WhereEvent, IncludeWithEvent]

	userService *UserService
}

func NewEventService(db *sql.DB, userService *UserService) *EventService {
	serviceBase := SetupResourceService[models.Event, CreateEvent, UpdateEvent, WhereEvent, IncludeWithEvent](db, "events", &models.Event{})

	eventService := &EventService{
		ResourceService: serviceBase,
		userService:     userService,
	}

	eventService.ResourceService.attachRelationsOverrider = eventService

	return eventService
}

func (s *EventService) AttachRelations(model *models.Event, include IncludeWithEvent) error {
	if include.User {
		user, err := s.userService.GetOneById(model.OwnerUserId, nil)
		if err != nil {
			return err
		}

		model.Owner = user
	}

	return nil
}
