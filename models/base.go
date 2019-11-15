package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"path/filepath"
	"html/template"
	"github.com/jinzhu/gorm"
	"os"
	"fmt"
	"github.com/joho/godotenv"
)

var (
	db *gorm.DB
	templates *template.Template
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load the env file!")
		panic(err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	startDB(username, password, dbName, dbHost)

	pathTemplate := os.Getenv("path_template")
	startTemplates(pathTemplate)
}

func startDB(user, password, dbname, host string) {
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Println("Failed to connect to database!")
		panic(err)
	}

	db = conn

	// Will add columns to tables if the struct types have changed.
	db.Debug().AutoMigrate(&Picture{})
}

func startTemplates(path string) {
	pathTemplate := path
	if pathTemplate == "" {
		pathTemplate, _ = os.Getwd()
	}

	templates = template.Must(template.ParseFiles(
		filepath.Join(pathTemplate, "templates", "header.html"),
		filepath.Join(pathTemplate, "templates", "menu.html"),
		filepath.Join(pathTemplate, "templates", "body_upper.html"),
		filepath.Join(pathTemplate, "templates", "body_lower.html"),
		filepath.Join(pathTemplate, "templates", "piclist.html"),
	))
}

// GetDB returns a handle to database object.
func GetDB() (*gorm.DB) {
	return db
}

// GetHTMLTemplates returns a handle to the HTML Templates.
func GetHTMLTemplates() (*template.Template) {
	return templates
}
