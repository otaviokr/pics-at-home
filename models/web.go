package models

import (
	"encoding/json"
	"strconv"
	"fmt"
	"image/jpeg"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/otaviokr/pics-at-home/utils"
)

// CreatePicWeb inserts a new picture into database.
func (a *App) CreatePicWeb(w http.ResponseWriter, r *http.Request) {
	pic := &Picture{}
	err := json.NewDecoder(r.Body).Decode(pic)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	if validation, ok := pic.Validate(a.GetDB()); !ok {
		fmt.Println("Error validating new picture data.")
		fmt.Println(validation)
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	response := pic.Create(a.GetDB())
	utils.Respond(w, response)
}

// GetRandomPicWeb fetches a random picture from the server.
func (a *App) GetRandomPicWeb(w http.ResponseWriter, r *http.Request) {
	pic := GetRandomPicture(a.GetDB())

	if pic == nil {
		fmt.Println("No random Pic returned!")
	}

	jpeg.Encode(w, pic, nil)
}

// GetRecentPicsWeb will fetch a number of recent pictures (recently added) from the server.
func (a *App) GetRecentPicsWeb(w http.ResponseWriter, r *http.Request) {
	var pics []Picture

	pics = GetRecentPics(20, a.GetDB())
	renderTemplateList(w, "piclist", pics, *a)
}

// GetDetailPicWeb get the picture with ID specified in URL.
func (a *App) GetDetailPicWeb(w http.ResponseWriter, r *http.Request) {
	pic := &Picture{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["picID"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Could not parse ID from request"))
		return
	}

	pic = GetPictureByID(uint(id), a.GetDB())
	renderTemplate(w, "detail", pic, *a)
}

func renderTemplate(w http.ResponseWriter, tmpl string, picture *Picture, a App) {
	err := a.GetHTMLTemplates().ExecuteTemplate(w, tmpl+".html", picture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateList(w http.ResponseWriter, tmpl string, pictures []Picture, a App) {
	err := a.GetHTMLTemplates().ExecuteTemplate(w, tmpl+".html", pictures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
