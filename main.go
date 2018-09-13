package main

import (
	"fmt"
	"go_rest_pg_starter/config"
	"go_rest_pg_starter/controllers"
	"go_rest_pg_starter/email"
	"go_rest_pg_starter/middlewares"
	"go_rest_pg_starter/models"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	/*
		Config, Setup services
	*/

	config := config.GetConfig()
	dbConfig := config.Database
	services, err := models.NewServices(
		models.WithGorm(dbConfig.Dialect(), dbConfig.ConnectionInfo()),
		models.WithLogMode(!config.IsProd()),
		models.WithUser(config.Pepper, config.HMACKey),
		models.WithPost(),
	)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	services.AutoMigrate()
	// services.DestructiveReset() // Comment this out for not resetting DB everytime it restarts

	/*
		Mailgun setup
	*/
	mailgunConfig := config.Mailgun
	emailer := email.NewClient(
		email.WithSender("Support", "support@"+mailgunConfig.Domain),
		email.WithMailgun(mailgunConfig.Domain, mailgunConfig.APIKey, mailgunConfig.PublicAPIKey),
	)

	r := mux.NewRouter()
	r = r.PathPrefix("/api").Subrouter()

	/*
		Defines controllers
	*/
	usersCtrl := controllers.NewUsers(services.User, emailer)
	postsCtrl := controllers.NewPosts(services.Post)

	userMW := middlewares.User{
		UserService: services.User,
	}

	/*
		Users routes
	*/
	r.HandleFunc("/signup", usersCtrl.Create).Methods("POST")
	r.HandleFunc("/login", usersCtrl.Login).Methods("POST")
	r.HandleFunc("/forgot_password", usersCtrl.InitiateReset).Methods("POST")
	r.HandleFunc("/update_password", usersCtrl.CompleteReset).Methods("POST")
	r.HandleFunc("/me", userMW.RequireUser(usersCtrl.Me)).Methods("GET")

	/*
		Posts routes
	*/
	r.HandleFunc("/posts", userMW.RequireUser(postsCtrl.Create)).Methods("POST")
	r.HandleFunc("/posts/{id:[0-9]+}", postsCtrl.GetOne).Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}/update", userMW.RequireUser(postsCtrl.Update)).Methods("PUT")
	r.HandleFunc("/posts/{id:[0-9]+}/delete", userMW.RequireUser(postsCtrl.Delete)).Methods("DELETE")

	fmt.Printf("Starting the server on localhost:%d...\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), middlewares.PassSignKey(r))
}
