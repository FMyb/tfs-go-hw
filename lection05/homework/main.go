package main

import (
	"log"
	"net/http"

	"github.com/FMyb/tfs-go-hw/lection05/homework/controllers"
	"github.com/FMyb/tfs-go-hw/lection05/homework/data"
	"github.com/FMyb/tfs-go-hw/lection05/homework/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	messages := controllers.NewMessages(data.NewMessages())
	users := controllers.NewUsers(data.NewUsers())
	root.Post("/login", users.Login)
	root.Post("/users", users.UserRegister)

	r := chi.NewRouter()
	r.Use(utils.JwtAuthentication)
	r.Get("/messages", messages.GetMessagesFromPublic)
	r.Post("/messages", messages.SendToPublic)

	r.Get("/users/me/messages", messages.GetUserMessages)
	r.Post("/users/me/messages", messages.SendToUser)

	root.Mount("/", r)

	log.Fatal(http.ListenAndServe(":5000", root))
}
