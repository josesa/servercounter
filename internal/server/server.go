package server

import (
	"fmt"
	"net/http"

	"github.com/josesa/servercounter/internal/service"
)

type webserver struct {
	service *service.HitCounter
}

func New(hc *service.HitCounter) *webserver {
	ws := &webserver{
		service: hc,
	}

	return ws
}

// Request handles the requests and presents the expected value for presentation
func (ws *webserver) Request(w http.ResponseWriter, r *http.Request) {
	count := ws.service.IncAndGetCount()

	// Formats output
	output := fmt.Sprintf("%d", count)
	w.Write([]byte(output))
}
