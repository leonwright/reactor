package server

import (
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/leonwright/reactor/routes"
	"github.com/leonwright/reactor/utils"
)

func registerRoutes() {
	jwtMiddleware := utils.CreateMiddleWare()
	r := mux.NewRouter()
	r.Handle("/api/private", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			message := "Hello from a private endpoint! You need to be authenticated to see this."
			utils.SendJSONResponse(message, w)
		}))))
	// This route is only accessible if the user has a valid Access Token with the read:messages scope
	// We are chaining the jwtmiddleware middleware into the negroni handler function which will check
	// for a valid token and scope.
	r.Handle("/api/private-scoped", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
			token := authHeaderParts[1]

			hasScope, _ := utils.CheckScope("read:self", token)

			if !hasScope {
				message := "Insufficient scope."
				utils.SendJSONResponse(message, w)
				return
			}
			message := "Hello from a private endpoint! You need to be authenticated to see this."
			utils.SendJSONResponse(message, w)
		}))))
	r.Handle("/repos", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(routes.GetUserRepositories))))
	r.Handle("/register", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(routes.RegisterWithBuilder))))
	http.Handle("/", r)
}

// StartServer starts the Reactor REST API.
func StartServer(cfg utils.Config) {
	registerRoutes()
	var listenOn = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.RestPort)
	deb.Infof("Starting REST API at port %s", listenOn)
	if err := http.ListenAndServe(listenOn, nil); err != nil {
		log.Fatal(err)
	}
}
