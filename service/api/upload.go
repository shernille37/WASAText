package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shernille37/WASAText/service/api/constants"
	"github.com/shernille37/WASAText/service/api/reqcontext"
)

type ImageURL struct {
	Image string `json:"image"`
}

// var validImages = map[string]bool{
// 	".gif":  true,
// 	".jpeg": true,
// 	".png":  true,
// 	".jpg":  true,
// }

var imageDirectory = "/tmp/images"

func uploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Create the image directory if it doesn't exist with permission (0755)
	if err := os.MkdirAll(imageDirectory, 0755); err != nil {
		http.Error(w, "Failed to create directory", http.StatusBadRequest)
		return
	}

	var res ImageURL

	// Parse the multipart form
	// Limit to 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, constants.PARSE_ERROR, http.StatusBadRequest)
		return
	}

	// Retrieve the file
	file, handler, err := r.FormFile("image")
	if err != nil {
		// http.Error(w, constants.INVALID_IMAGE, http.StatusBadRequest)
		http.Error(w, "Cannot retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// if !validImages[filepath.Ext(handler.Filename)] {
	// 	http.Error(w, constants.INVALID_IMAGE, http.StatusBadRequest)
	// 	return
	// }

	// Check if the file has a valid image extension by checking the MIME Type
	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	file.Seek(0, 0)

	// Detect MIME type
	mimeType := http.DetectContentType(buffer)
	if !strings.HasPrefix(mimeType, "image/") {
		http.Error(w, constants.INVALID_IMAGE, http.StatusBadRequest)
		return
	}

	// Create a unique filename
	fileID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, constants.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	uniqueFilename := fmt.Sprintf("%s-%s", fileID.String(), handler.Filename)
	// Create and Save file to /tmp/images
	dstPath := filepath.Join("/tmp/images", uniqueFilename)
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

	res.Image = fmt.Sprintf("/images/%s", uniqueFilename)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}
