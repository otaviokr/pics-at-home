package controllers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/otaviokr/pics/models"
	"github.com/otaviokr/pics/utils"
)

// CreatePicAPI inserts a new picture into database.
var CreatePicAPI = func(w http.ResponseWriter, r *http.Request) {
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


// GetRandomPicAPI fetches a random picture from the server.
var GetRandomPicAPI = func(w http.ResponseWriter, r *http.Request) {
	pic := &models.Picture{}

	pic = models.GetRandomPictureInfo()
	
	response := utils.Message(true, "success")
	response["picture"] = pic
	utils.Respond(w, response)
}
