package handle

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AdminUploadImage handles multipart file upload and stores files under ./uploads/
func AdminUploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not found in form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		http.Error(w, "Unable to create uploads dir", http.StatusInternalServerError)
		return
	}

	fname := time.Now().Format("20060102150405") + "_" + filepath.Base(header.Filename)
	outPath := filepath.Join("./uploads", fname)
	out, err := os.Create(outPath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	apiBase := os.Getenv("API_BASE_URL")
	var url string
	if apiBase != "" {
		url = strings.TrimRight(apiBase, "/") + "/uploads/" + fname
	} else {
		url = "/uploads/" + fname
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"url":"%s"}` , url)
}
