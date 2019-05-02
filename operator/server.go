package operator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kris-nova/logger"
)

type ServerConfiguration struct {
	Port        int
	BindAddress string
}

func ListenAndWait(cfg *ServerConfiguration) error {
	r := mux.NewRouter()
	r.HandleFunc("/requestResources", requestResources)
	logger.Always("Listening on %s:%d...", cfg.BindAddress, cfg.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.Port), r)
	return nil
}

func requestResources(w http.ResponseWriter, r *http.Request) {

	// Just for testing:
	UpdateCRDNumberInstances(10)
	w.WriteHeader(200)
	w.Write([]byte("Status: OK\n"))
	return

	decoder := json.NewDecoder(r.Body)
	mutation := &SparkClusterApiOperatorRequest{}
	err := decoder.Decode(&mutation)
	if err != nil {
		//logger.Warning("JSON Error: %v", err)
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("JSON Error: %v\n", err)))
		return
	}
	logger.Always("Request: %+v", mutation)

	// Todo Nova - Set this to the size of whatever instance we are using
	serverSize := &SparkClusterApiOperatorRequest{
		CPUCount:       1,
		MemoryMBs:      1000,
		ContainerCount: 1,
	}
	numInstances := ComputeNumberOfExpectedInstances(serverSize, mutation)
	if numInstances < 0 {
		//logger.Warning("Invalid request\n")
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("Error updated resources: %v", err)))
		return
	}
	err = UpdateCRDNumberInstances(numInstances)
	if err != nil {
		//logger.Warning("Error updated resources: %v\n", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error updated resources: %v", err)))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Status: OK\n"))
	return
}
