package main

import (
	"api-gateway/config"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("example")

type APIGateway struct {
	Services []string
}

// Forward request đến các backend server
func (g *APIGateway) ForwardRequests(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)

	var prettyJSON bytes.Buffer

	if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err != nil {
		logger.Debug("Invalid JSON body, logging raw body", string(bodyBytes))
	} else {
		logger.Debug("Parsed JSON Body", prettyJSON.String())
	}

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// dynamic path
	vars := mux.Vars(r) // Lấy các biến từ route
	dynamicPath := vars["any"]

	logger.Debug("dynamicPath", dynamicPath)
	// Đóng Body để tránh rò rỉ tài nguyên
	defer r.Body.Close()

	for _, service := range g.Services {
		go func(service string) {
			// Tạo URL mới cho backend
			serviceURL, err := url.Parse(service)
			if err != nil {
				log.Printf("Invalid backend URL: %s\n", service)
				return
			}
			logger.Debug("Method", r.Method)

			// Gắn path vào backend URL
			//r.URL.ResolveReference(serviceURL).String()
			backendURL := serviceURL.ResolveReference(&url.URL{
				Path:     dynamicPath,
				RawQuery: r.URL.RawQuery,
			})
			logger.Debug("url : ", backendURL)

			proxyReq, err := http.NewRequest(r.Method, backendURL.String(), bytes.NewReader(bodyBytes))
			if err != nil {
				log.Printf("Failed to create request for backend %s: %v\n", serviceURL, err)
				return
			}

			// Sao chép headers từ request gốc
			proxyReq.Header = r.Header

			logger.Debug("Header request :", proxyReq.Header)

			// Gửi request đến backend
			client := &http.Client{}
			resp, err := client.Do(proxyReq)
			if err != nil {
				log.Printf("Error forwarding to backend %s: %v\n", serviceURL, err)
				return
			}

			log.Printf("Response from %s: %d\n", serviceURL, resp.StatusCode)
		}(service)
	}

	// Trả về response thành công
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request forwarded to services"))
}

func main() {
	// load config
	config, err := config.LoadConfigDockerfile()

	if err != nil {
		logger.Error("error when get config, %v", err)
		return
	}

	// Tạo API Gateway
	gateway := &APIGateway{Services: config.Servers}

	// dung mux handle dynamic router
	router := mux.NewRouter()
	router.HandleFunc("/{any:.*}", gateway.ForwardRequests).Methods("GET", "POST", "PUT", "DELETE")

	log.Println("API Gateway running on :8080")
	http.ListenAndServe(":8080", router)
}
