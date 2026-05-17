package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func Log() {
	err := os.Mkdir("backend_logs", os.ModePerm)

	if err != nil {
		log.Println("[WARN]not dir log", err)
	}

	logFile, err := os.OpenFile("backend_logs/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("[WARN]didn't open file ", err)
	}

	// Инициализируем наш клиент Loki
	loki := NewLokiWriter()

	var multiwriter io.Writer
	if loki != nil {
		// Если URL есть, пишем в консоль, файл И отправляем в Loki
		multiwriter = io.MultiWriter(os.Stdout, logFile, loki)
	} else {
		// Если URL нет (локально), пишем только в консоль и файл
		multiwriter = io.MultiWriter(os.Stdout, logFile)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiwriter)

}
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
func LogMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r) 
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[REQUEST] method = %s path = %s ip = %s user_agent = %s latency = %s",
			r.Method,
			r.URL.Path,
			ip,
			r.UserAgent(),
			time.Since(start),
		)
	})
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LokiPushRequest struct {
	Streams []LokiStream `json:"streams"`
}

// LokiWriter реализует io.Writer для прямой отправки в Loki
type LokiWriter struct {
	URL        string
	HTTPClient *http.Client
	Labels     map[string]string
}

func NewLokiWriter() *LokiWriter {
	url := os.Getenv("LOKI_URL")
	if url == "" {
	log.Println("err set link Loki")
		return nil
	}
	return &LokiWriter{
		URL:url,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
		Labels: map[string]string{
			"app": "go-backend",
			"env": "huggingface",
		},
	}
}

func (l *LokiWriter) Write(p []byte) (n int, err error) {
	// Возвращаем длину, чтобы стандартный логгер думал, что всё успешно записано
	n = len(p)

	// Очищаем строку от лишних переносов для корректного отображения в Grafana
	message := string(bytes.TrimSpace(p))

	// Отправляем в фоне (goroutine), чтобы сетевые задержки Loki не тормозили бэкенд
	go func(msg string) {
		nowNano := fmt.Sprintf("%d", time.Now().UnixNano())

		payload := LokiPushRequest{
			Streams: []LokiStream{
				{
					Stream: l.Labels,
					Values: [][]string{
						{nowNano, msg},
					},
				},
			},
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return
		}

		req, err := http.NewRequest("POST", l.URL, bytes.NewBuffer(body))
		if err != nil { 
				fmt.Printf("[LOKI-ERROR] Failed to send log: %v\n", err)
			return
		} 
	
		req.Header.Set("Content-Type", "application/json")

		resp, err := l.HTTPClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close() 
		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			var buf bytes.Buffer
			buf.ReadFrom(resp.Body)
			fmt.Printf("[LOKI-ERROR] Bad status %d: %s\n", resp.StatusCode, buf.String())
		}
	}(message)

	return n, nil
}
