package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AdminUploadImage uploads file to Supabase Storage using POST
func AdminUploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// limit to 5MB
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		log.Println("[ERROR] Failed to parse multipart form:", err)
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("[ERROR] File not found in form:", err)
		http.Error(w, "File not found in form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Println("[LOG] Received file:", header.Filename, "size", header.Size)

	// validate extension
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowed[ext] {
		log.Println("[WARN] Disallowed file extension:", ext)
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, file); err != nil {
		log.Println("[ERROR] Error reading file:", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	fname := time.Now().Format("20060102150405") + "_" + filepath.Base(header.Filename)

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")
	supabaseBucket := os.Getenv("SUPABASE_BUCKET")
	if supabaseBucket == "" {
		supabaseBucket = "public"
	}

	if supabaseURL == "" || supabaseKey == "" {
		log.Println("[ERROR] SUPABASE_URL or SUPABASE_SERVICE_KEY not set")
		http.Error(w, "Server upload not configured", http.StatusInternalServerError)
		return
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", strings.TrimRight(supabaseURL, "/"), supabaseBucket, fname)

	req, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Println("[ERROR] Failed to create upload request:", err)
		http.Error(w, "Failed to create upload request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	ct := header.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	req.Header.Set("Content-Type", ct)

	client := &http.Client{Timeout: 30 * time.Second}
	log.Println("[LOG] Uploading to Supabase:", uploadURL)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] Upload to Supabase failed:", err)
		http.Error(w, "Upload to Supabase failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Println("[ERROR] Supabase upload error:", string(body))
		http.Error(w, fmt.Sprintf("Supabase upload error: %s", string(body)), http.StatusInternalServerError)
		return
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", strings.TrimRight(supabaseURL, "/"), supabaseBucket, fname)
	log.Println("[LOG] Upload succeeded:", publicURL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"url": publicURL})
}
