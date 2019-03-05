package http

import (
	"fmt"
	"io.pburakov/homehub/agent/config"
	"io.pburakov/homehub/agent/util"
	"log"
	"net/http"
)

func ServeFolder(c *config.Configuration) {
	log.Printf("Starting http file server on port %d", c.WebServer.Port)
	fs := http.FileServer(http.Dir(c.Motion.Dir))
	http.Handle("/", fs)
	e := http.ListenAndServe(fmt.Sprintf(":%d", c.WebServer.Port), nil)
	if e != nil {
		util.Fatal(fmt.Errorf("error starting http file server: %s", e))
	}
}
