package handle

import (
"bytes"
"encoding/json"
"fmt"
"io"
"mime/multipart"
"net/http"
"os"
"path/filepath"
"strings"
"time"
)

// AdminUploadImage uploads file to Supabase storage if configured, otherwise saves locally.
func AdminUploadImage(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
w.WriteHeader(http.StatusMethodNotAllowed)
return
}

err := r.ParseMultipartForm(20 << 20) // 20 MB
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

buf := &bytes.Buffer{}
if _, err := io.Copy(buf, file); err != nil {
http.Error(w, "Error reading file", http.StatusInternalServerError)
return
}

fname := time.Now().Format("20060102150405") + "_" + filepath.Base(header.Filename)

supabaseURL := os.Getenv("SUPABASE_URL")
if supabaseURL == "" {
// fallback to DATABASE_URL if user configured it that way
supabaseURL = os.Getenv("DATABASE_URL")
}
supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")
supabaseBucket := os.Getenv("SUPABASE_BUCKET")
if supabaseBucket == "" {
supabaseBucket = "public"
}

if supabaseURL != "" && supabaseKey != "" {
uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", strings.TrimRight(supabaseURL, "/"), supabaseBucket, fname)

req, err := http.NewRequest(http.MethodPut, uploadURL, bytes.NewReader(buf.Bytes()))
if err != nil {
http.Error(w, "Failed to create upload request", http.StatusInternalServerError)
return
}
req.Header.Set("Authorization", "Bearer "+supabaseKey)
req.Header.Set("apikey", supabaseKey)
ct := header.Header.Get("Content-Type")
if ct == "" {
ct = "application/octet-stream"
}
req.Header.Set("Content-Type", ct)

client := &http.Client{Timeout: 30 * time.Second}
resp, err := client.Do(req)
if err != nil {
http.Error(w, "Upload to Supabase failed", http.StatusInternalServerError)
return
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
body, _ := io.ReadAll(resp.Body)
http.Error(w, fmt.Sprintf("Supabase upload error: %s", string(body)), http.StatusInternalServerError)
return
}

publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", strings.TrimRight(supabaseURL, "/"), supabaseBucket, fname)

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(map[string]string{"url": publicURL})
return
}

// fallback: save locally
if err := os.MkdirAll("./uploads", 0755); err != nil {
http.Error(w, "Unable to create uploads dir", http.StatusInternalServerError)
return
}

outPath := filepath.Join("./uploads", fname)
out, err := os.Create(outPath)
if err != nil {
http.Error(w, "Unable to create file", http.StatusInternalServerError)
return
}
defer out.Close()

if _, err := out.Write(buf.Bytes()); err != nil {
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
json.NewEncoder(w).Encode(map[string]string{"url": url})
}
