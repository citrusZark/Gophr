package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleHome(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db := NewDBImageStore()
	defer db.Close()
	images, err := db.FindAll(0)
	if err != nil {
		panic(err)
	}
	// Display homepage
	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": images,
	})
}
