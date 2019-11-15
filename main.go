package main

import (
	"github.com/gorilla/mux"
	"os"
	"fmt"
	"net/http"
	"github.com/otaviokr/pics/controllers"
)

const (
	// STATIC_DIR is the path used to serve static files.
	STATIC_DIR = "/static/"
)

func main() {
	router := mux.NewRouter()
	// TODO Define JWT Authentication.
	// router.Use(app.JwtAuthentication)

	// Handle to give a random picture.
	router.HandleFunc("/api/pic/random", controllers.GetRandomPicAPI).Methods("GET")
	router.HandleFunc("/api/pic/create", controllers.CreatePicAPI).Methods("POST")

	router.HandleFunc("/pic/random", controllers.GetRandomPicWeb).Methods("GET")
	router.HandleFunc("/pic/recent", controllers.GetRecentPicsWeb).Methods("GET")
	router.HandleFunc("/pic/detail/{picID:[0-9]+}", controllers.GetDetailPicWeb).Methods("GET")

	router.HandleFunc("/pic", controllers.GetRecentPicsWeb).Methods("GET")

	router.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir(os.Getenv("IMG_PATH")))))
	router.HandleFunc("/", controllers.GetRecentPicsWeb).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Port used: %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		fmt.Print(err)
	}
}