package logs

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func Log() {
	err := os.Mkdir("backend_logs", os.ModePerm)
	if err != nil {
		log.Println("[WARN]not dir log",err)

	}
	logFile, err := os.OpenFile("backend_logs/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("[WARN]didn't open file ",err)
	}
	multiwriter := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiwriter)

}
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {return ip 	} 
	ip = r.Header.Get("X-Real-IP")
	if ip != "" { 
		return ip 	 
	} 
	host,_,err := net.SplitHostPort(r.RemoteAddr) 
	if err !=nil{return  r.RemoteAddr}  
	return host
} 
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r)
		log.Printf("[REQUEST] method = %s path = %s ip = %s browser = %s", 
		r.Method,
		r.URL.Path,
		ip, 
		r.UserAgent(),
		)
	})
}
