package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kevindoubleu/warnetop/device"
	"github.com/kevindoubleu/warnetop/lib/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pog")
	})

	deviceRepo := device.NewDeviceRepositoryPSQL()
	defer deviceRepo.Close()
	deviceService := device.NewDeviceService(deviceRepo)
	deviceHandler := device.NewDeviceHandler(deviceService)

	completeDeviceHandler := chainMiddleware(
		deviceHandler,
		middleware.RequestResponseLogger,
		middleware.WithReqID,
	)
	mux.Handle("/device", completeDeviceHandler)

	server := &http.Server{
		Handler: mux,

		Addr:         "localhost:8888",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	fmt.Println(server.ListenAndServe())
}

func chainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
