package controllers

import (
	"net/http"
	"image/jpeg"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/otaviokr/pics/models"
	"github.com/gorilla/mux"
	"github.com/otaviokr/pics/utils"
)

// CreatePicWeb inserts a new picture into database.
var CreatePicWeb = func(w http.ResponseWriter, r *http.Request) {
	pic := &models.Picture{}
	err := json.NewDecoder(r.Body).Decode(pic)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	if validation, ok := pic.Validate(); !ok {
		fmt.Println("Error validating new picture data.")
		fmt.Println(validation)
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	response := pic.Create()
	utils.Respond(w, response)
}

// GetRandomPicWeb fetches a random picture from the server.
var GetRandomPicWeb = func(w http.ResponseWriter, r *http.Request) {
	pic := models.GetRandomPicture()

	if pic == nil {
		fmt.Println("No random Pic returned!")
	}
	
	jpeg.Encode(w, pic, nil)
}

// GetRecentPicsWeb will fetch a number of recent pictures (recently added) from the server.
var GetRecentPicsWeb = func(w http.ResponseWriter, r *http.Request) {
	var pics []models.Picture

	pics = models.GetRecentPics(20)
	renderTemplateList(w, "piclist", pics)
}

// GetDetailPicWeb get the picture with ID specified in URL.
var GetDetailPicWeb = func(w http.ResponseWriter, r *http.Request) {
	pic := &models.Picture{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["picID"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Could not parse ID from request"))
		return
	}

	pic = models.GetPictureByID(uint(id))
	renderTemplate(w, "detail", pic)
}

func renderTemplate(w http.ResponseWriter, tmpl string, picture *models.Picture) {
	err := models.GetHTMLTemplates().ExecuteTemplate(w, tmpl+".html", picture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateList(w http.ResponseWriter, tmpl string, pictures []models.Picture) {
	err := models.GetHTMLTemplates().ExecuteTemplate(w, tmpl+".html", pictures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}