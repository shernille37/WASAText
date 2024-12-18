package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/constants"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

type ImageURL struct {
	Image string `json:"image"`
}

var validImages = map[string]bool{
	".gif":  true,
	".jpeg": true,
	".png":  true,
}

func uploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var res ImageURL

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusBadRequest)
		return
	}

	// Retrieve the file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, constants.INVALID_IMAGE, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the file has a valid image extension
	if !validImages[filepath.Ext(handler.Filename)] {
		http.Error(w, constants.INVALID_IMAGE, http.StatusBadRequest)
		return
	}

	// Create and Save file to /tmp/images
	dstPath := filepath.Join("/tmp/images", handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content to destination
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	res.Image = dstPath

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}
