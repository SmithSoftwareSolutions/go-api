package main

import (
	"database/sql"

	"smithsolutions/go-api/internal/services"
)

type ServiceMap struct {
	UserService  *services.UserService
	EventService *services.EventService
}

func BuildServiceMap(db *sql.DB) *ServiceMap {
	// services
	userService := services.NewUserService(db, nil)
	eventService := services.NewEventService(db, userService)
	userService.SetEventService(eventService)

	return &ServiceMap{
		UserService:  userService,
		EventService: eventService,
	}
}
