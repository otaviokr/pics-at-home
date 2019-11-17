package models

import (
	"fmt"
	"net/http"
	"html/template"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	// Line below is to invoke the Postgres driver.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	// StaticDir is the path used to serve static files.
	StaticDir = "/static/"
)

// Config stores the configuration data.
type Config struct {
	DBUser string
	DBPassword string
	DBName string
	DBHost string
	DBPort string
	TemplatePath string
}

// App contains the key components of the application.
type App struct {
	db        *gorm.DB
	router    *mux.Router
	templates *template.Template
	config    *Config
}

// Initialize kickstarts the components
func (a *App) Initialize() {
	
	a.StartDB()
	a.StartTemplates()
	a.StartRouter()
}

// StartDB initializes database connection.
func (a *App) StartDB() {
	

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", 
		a.config.DBHost, a.config.DBUser, a.config.DBName, a.config.DBPassword)
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Println("Failed to connect to database!")
		panic(err)
	}

	a.db = conn

	// Will add columns to tables if the struct types have changed.
	a.db.Debug().AutoMigrate(&Picture{})
}

// StartTemplates will load the HTML template files.
func (a *App) StartTemplates() {
	pathTemplate := a.config.TemplatePath
	if pathTemplate == "" {
		pathTemplate, _ = os.Getwd()
	}

	a.templates = template.Must(template.ParseFiles(
		filepath.Join(pathTemplate, "templates", "header.html"),
		filepath.Join(pathTemplate, "templates", "menu.html"),
		filepath.Join(pathTemplate, "templates", "body_upper.html"),
		filepath.Join(pathTemplate, "templates", "body_lower.html"),
		filepath.Join(pathTemplate, "templates", "piclist.html"),
	))
}

// StartRouter will setup the HTTP request router.
func (a *App) StartRouter() {
	a.router = mux.NewRouter()
	// TODO Define JWT Authentication.
	// router.Use(app.JwtAuthentication)

	// Handle to give a random picture.
	a.router.HandleFunc("/api/pic/random", a.GetRandomPicAPI).Methods("GET")
	a.router.HandleFunc("/api/pic/create", a.CreatePicAPI).Methods("POST")

	a.router.HandleFunc("/pic/random", a.GetRandomPicWeb).Methods("GET")
	a.router.HandleFunc("/pic/recent", a.GetRecentPicsWeb).Methods("GET")
	a.router.HandleFunc("/pic/detail/{picID:[0-9]+}", a.GetDetailPicWeb).Methods("GET")

	a.router.HandleFunc("/pic", a.GetRecentPicsWeb).Methods("GET")

	a.router.PathPrefix(StaticDir).Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(os.Getenv("IMG_PATH")))))
	a.router.HandleFunc("/", a.GetRecentPicsWeb).Methods("GET")
}

// ListenAndServe is just a wrapper.
func (a *App) ListenAndServe() (error) {
	port := a.config.DBPort
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Port used: %s\n", port)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), a.router)
}

// SetDB sets the database handle.
func (a *App) SetDB(db *gorm.DB) {
	a.db = db
}

// GetDB returns a handle to database object.
func (a *App) GetDB() *gorm.DB {
	return a.db
}

// GetHTMLTemplates returns a handle to the HTML Templates.
func (a *App) GetHTMLTemplates() *template.Template {
	return a.templates
}

// GetRouter returns a handle to the router.
func (a *App) GetRouter() *mux.Router {
	return a.router
}

// SetConfig sets the configuration handle.
func (a *App) SetConfig(config *Config) {
	a.config = config
}

// GetConfig returns a handle to the config struct.
func (a *App) GetConfig() *Config {
	return a.config
}