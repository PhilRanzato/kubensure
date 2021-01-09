package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	backend "github.com/PhilRanzato/kubensure/backend"
)

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	cs := backend.GetClientSet()

	fmt.Println("Get services")

	svcs, _ := json.Marshal(backend.GetServices(cs))

	json.NewEncoder(w).Encode(string(svcs))

}
