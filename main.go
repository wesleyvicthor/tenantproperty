package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wesleyvicthor/tenantproperty/internal"
)

func main() {
	handler := http.NewServeMux()
	srv := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	properties := internal.NewPropertyRepository()
	handler.Handle("/property/match", internal.NewPropertyMatchHandler(properties))
	handler.Handle("/properties", internal.NewPropertiesHandler(properties))
	handler.Handle("/properties/", http.StripPrefix("/properties/", internal.NewPropertyHandler(properties)))

	go func() {
		fmt.Println("Listening on ... :8080")
		srv.ListenAndServe()
	}()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-term
	fmt.Println("Shutting down...")
	srv.Shutdown(context.Background())
	fmt.Println("Done!")
}
