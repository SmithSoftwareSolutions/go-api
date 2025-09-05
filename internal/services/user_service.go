package services

import (
	"database/sql"

	"smithsolutions/go-api/internal/filters"
	"smithsolutions/go-api/internal/models"
	"smithsolutions/go-api/internal/util"
)

//go:generate gen CreateUser CreateSQL
type CreateUser struct {
	Email        string
	PasswordHash string
}

// func (c CreateUser) SQL() ([]string, []any, error) {
// 	return util.GetCreateSQL(c)
// }

type UpdateUser struct {
}

func (u UpdateUser) SQL() (string, []any, error) {
	return util.GetUpdateSQL(u)
}

type WhereUser struct {
	Email *filters.StringFilter
}

func (w WhereUser) SQL() (string, []any, error) {
	return util.GetWhereSQL(w)
}

type IncludeWithUser struct {
	Events bool
}

type UserService struct {
	ResourceService[models.User, CreateUser, UpdateUser, WhereUser, IncludeWithUser]

	// services
	eventService *EventService
}

func NewUserService(db *sql.DB, eventService *EventService) *UserService {
	serviceBase := SetupResourceService[models.User, CreateUser, UpdateUser, WhereUser, IncludeWithUser](db, "users", &models.User{})

	userService := &UserService{
		ResourceService: serviceBase,
		eventService:    eventService,
	}

	userService.ResourceService.attachRelationsOverrider = userService

	return userService
}

// post initialization dependency injection
func (s *UserService) SetEventService(eventService *EventService) {
	s.eventService = eventService
}

func (s *UserService) AttachRelations(model *models.User, include IncludeWithUser) error {
	if include.Events {
		events, err := s.eventService.GetMany(WhereEvent{
			OwnerUserId: filters.IntEquals(model.Id),
		}, nil)
		if err != nil {
			return err
		}
		model.Events = events
	}

	return nil
}
