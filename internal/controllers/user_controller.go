package controllers

import (
	"net/http"

	"smithsolutions/go-api/internal/core"
	"smithsolutions/go-api/internal/services"
)

type QueryManyUsers struct {
	where *services.WhereUser
}

type UserController struct {
	mux         *http.ServeMux
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	mux := http.NewServeMux()

	controller := &UserController{
		mux:         mux,
		userService: userService,
	}

	controller.setupEndpoints()

	return controller
}

func (c *UserController) setupEndpoints() {
	c.mux.HandleFunc("/", c.GetMany)
}

func (c *UserController) GetMux() *http.ServeMux {
	return c.mux
}

func (c *UserController) GetMany(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetMany(services.WhereUser{}, nil)

	if err != nil {
		core.WriteJSON(w, http.StatusInternalServerError, &core.Response{
			Error: err.Error(),
		})
		return
	}

	response := core.Response{
		Data: users,
	}
	core.WriteJSON(w, http.StatusOK, response)
}
