package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shernille37/WASAText/service/api"
)

func main() {

	if err := run(); err != nil {
		os.Exit(1)
	}
	
}

func run() error {

	fmt.Println("Init API Server")

	apirouter, err := api.New()

	if err != nil {
		return fmt.Errorf("error starting api: %w", err)
	}
	
	router := apirouter.Handler()

	apiServer := http.Server {
		
		Handler: router,
	}

	go func() {
		fmt.Println("API Started")
		apiServer.ListenAndServe()
	}()

	
	return nil

}