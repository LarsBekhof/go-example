package routes

import (
	"github.com/gorilla/mux"
	"go-example/internal/controllers"
)

func RegisterRoutes(router *mux.Router) {
	bookController := controllers.InitBookController()

	router.HandleFunc("/users", controllers.InitUserController).Methods(controllers.UserMethods)
	router.HandleFunc("/users/{user_id}/books", bookController.AbstractController.HandleRequest).Methods(bookController.ImplementedMethods...)
}
