package http

import (
	"fmt"
	"github.com/abbot/go-http-auth"
	"io.pburakov/homehub/agent/config"
	"io.pburakov/homehub/agent/util"
	"log"
	"net/http"
)

var (
	salt  = []byte("salt")
	magic = []byte("$1$")
)

func ProvideSecret(w *config.WebServer) func(string, string) string {
	return func(u string, r string) string {
		if u == w.Username {
			hash := string(auth.MD5Crypt([]byte(w.Password), []byte(salt), []byte(magic)))
			return hash
		}
		return ""
	}
}

func ServeWithAuthHandler(a *auth.BasicAuth, c *config.Configuration) http.HandlerFunc {
	return a.Wrap(func(res http.ResponseWriter, req *auth.AuthenticatedRequest) {
		http.FileServer(http.Dir(c.Motion.Dir)).ServeHTTP(res, &req.Request)
	})
}

func ServeFolder(c *config.Configuration) {
	log.Printf("Starting http file server on port %d", c.WebServer.Port)
	a := auth.NewBasicAuthenticator("homehub.io", ProvideSecret(&c.WebServer))
	http.HandleFunc("/", ServeWithAuthHandler(a, c))
	e := http.ListenAndServe(fmt.Sprintf(":%d", c.WebServer.Port), nil)
	if e != nil {
		util.Fatal(fmt.Errorf("error starting http file server: %s", e))
	}
}
