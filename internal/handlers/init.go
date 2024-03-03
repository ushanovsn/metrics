package handlers

import (
	"net/http"
)

func ServerMux() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", startPage)

	mux.Handle("/update/", http.StripPrefix("/update/", http.HandlerFunc(updatePage)))



	return mux
}