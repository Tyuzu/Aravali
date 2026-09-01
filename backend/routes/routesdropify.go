package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddStaticRoutes(router *httprouter.Router) {
	// Serve static uploaded files directly using standard file server
	router.ServeFiles("/static/uploads/*filepath", http.Dir("static/uploads"))

}
