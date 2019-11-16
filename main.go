package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/otaviokr/pics-at-home/controllers"
	"github.com/otaviokr/pics-at-home/models"
)

const (
	// StaticDir is the path used to serve static files.
	StaticDir = "/static/"
)

func main() {

	models.StartDB()

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

	router.PathPrefix(StaticDir).Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(os.Getenv("IMG_PATH")))))
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
