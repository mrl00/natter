// Binário healthcheck verifica se o servidor está respondendo corretamente.
// Utilizado pelo Docker healthcheck para determinar se o container está saudável.
//
// Variáveis de ambiente:
//   - APP_HOST: host do servidor (padrão: localhost).
//   - APP_PORT: porta do servidor (padrão: 8080).
package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	targetURL := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/health",
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(targetURL.String())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Erro ao conectar: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("OK")
		os.Exit(0)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Status inesperado: %d\n", resp.StatusCode)
	os.Exit(1)
}
