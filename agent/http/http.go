package http

import (
	"fmt"
	"github.com/abbot/go-http-auth"
	"io.pburakov/homehub/agent/config"
	"io.pburakov/homehub/agent/util"
	"log"
	"net/http"
)

func Secret(user, realm string) string {
	if user == "john" {
		// password is "hello"
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}

func ServeFolder(c *config.Configuration) {
	log.Printf("Starting http file server on port %d", c.WebServer.Port)
	a := auth.NewBasicAuthenticator("homehub.io", Secret)
	http.HandleFunc("/", a.Wrap(func(res http.ResponseWriter, req *auth.AuthenticatedRequest) {
		http.FileServer(http.Dir(c.Motion.Dir)).ServeHTTP(res, &req.Request)
	}))
	e := http.ListenAndServe(fmt.Sprintf(":%d", c.WebServer.Port), nil)
	if e != nil {
		util.Fatal(fmt.Errorf("error starting http file server: %s", e))
	}
}
