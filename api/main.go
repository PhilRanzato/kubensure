package main

import (
	"fmt"
	"net/http"

	"github.com/PhilRanzato/kubensure/api/api"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/pods", api.PodHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/services", api.ServiceHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/pod-exec", api.PodExecHandler).Methods("POST")
	r.HandleFunc("/api/test-connection", api.TestConnection).Methods("POST")
	r.HandleFunc("/api/pvc", api.PvcHandler).Methods("GET")
	return r
}

func main() {
	r := newRouter()
	fmt.Println("Server listening on port 80")
	// enable CORS
	router := handlers.CORS(handlers.AllowedHeaders([]string{"Accept", "X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)
	http.ListenAndServe(":80", router)
}
