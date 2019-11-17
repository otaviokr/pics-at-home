package main

import (
	"fmt"
	"os"

	"github.com/otaviokr/pics-at-home/models"
	"github.com/joho/godotenv"
)

func main() {

	a := models.App{}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load the env file!")
	panic(err)
	}
	
	config := models.Config{}
	config.DBUser = os.Getenv("db_user")
	config.DBPassword = os.Getenv("db_pass")
	config.DBName = os.Getenv("db_name")
	config.DBHost = os.Getenv("db_host")
	config.DBPort = os.Getenv(("db_port"))

	config.TemplatePath = os.Getenv("template_path")

	a.SetConfig(&config)

	a.Initialize()
	
	err = a.ListenAndServe()
	if err != nil {
		fmt.Print(err)
	}
}
