package main

import (
	"net/http"

	"smithsolutions/go-api/internal/controllers"
	"smithsolutions/go-api/internal/core"
)

func RegisterHandlers(serviceMap ServiceMap) *http.ServeMux {

	rootMux := http.NewServeMux()

	rootMux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {

		res := struct{ Message string }{
			Message: "Stable",
		}

		core.WriteJSON(w, 200, res)
	})

	userController := controllers.NewUserController(serviceMap.UserService)

	registerHandler(rootMux, "/users/", userController.GetMux())

	return rootMux
}

func registerHandler(rootMux *http.ServeMux, pattern string, handler http.Handler) {
	rootMux.Handle(pattern, http.StripPrefix(pattern[:len(pattern)-1], handler))
}
