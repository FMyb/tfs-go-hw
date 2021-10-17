package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/FMyb/tfs-go-hw/lection05/homework/controllers"
	"github.com/FMyb/tfs-go-hw/lection05/homework/utils"
)

func main() {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Post("/login", controllers.Login)
	root.Post("/users", controllers.UserRegister)

	r := chi.NewRouter()
	r.Use(utils.JwtAuthentication)
	r.Get("/messages", controllers.GetMessagesFromPublic)
	r.Post("/messages", controllers.SendToPublic)

	r.Get("/users/me/messages", controllers.GetUserMessages)
	r.Post("/users/me/messages", controllers.SendToUser)

	root.Mount("/", r)

	log.Fatal(http.ListenAndServe(":5000", root))
}
