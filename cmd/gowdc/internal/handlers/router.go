package handlers

import (
	"gowdc/internal/web/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(logger *log.Logger) http.Handler {

	// WDC
	wdc := GoWDC{
		Log: logger,
	}

	mux := mux.NewRouter()
	mux.Handle("/api", middleware.Chain(wdc.Get(), middleware.Logging())).Methods("GET")

	// Hosting for Tableau Web Data Connector
	mux.PathPrefix("/").Handler(
		http.StripPrefix("/",
			http.FileServer(http.Dir("public"))))

	return mux
}
