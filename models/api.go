package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/otaviokr/pics-at-home/utils"
)

// CreatePicAPI inserts a new picture into database.
func (a *App) CreatePicAPI (w http.ResponseWriter, r *http.Request) {
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

// GetRandomPicAPI fetches a random picture from the server.
func (a *App) GetRandomPicAPI(w http.ResponseWriter, r *http.Request) {
	pic := Picture{}

	pic = GetRandomPictureInfo(a.GetDB())

	response := utils.Message(true, "success")
	response["picture"] = pic
	utils.Respond(w, response)
}
