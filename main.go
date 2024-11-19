package main

import (
	"api-gateway/config"
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("example")

type APIGateway struct {
	Services []string
}

// Forward request đến các backend server
func (g *APIGateway) ForwardRequests(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
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
			proxyReq, err := http.NewRequest(r.Method, r.URL.ResolveReference(serviceURL).String(), bytes.NewReader(bodyBytes))
			if err != nil {
				log.Printf("Failed to create request for backend %s: %v\n", serviceURL, err)
				return
			}

			// Sao chép headers từ request gốc
			proxyReq.Header = r.Header

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

	http.HandleFunc("/", gateway.ForwardRequests)

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
