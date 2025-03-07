package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func printRequestBody(r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}
	log.Printf("%s\n", string(body))
}

func noValidationHandler(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r)
	fmt.Fprintf(w, "No validation needed.\n")
}

func tokenValidationHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: No Authorization header provided", http.StatusUnauthorized)
		return
	}

	// 检查Bearer令牌的格式
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		http.Error(w, "Unauthorized: Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	// 提取令牌
	token := authHeader[len(prefix):]
	if token != "token" {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		log.Printf("Unauthorized: Invalid token\n")
		return
	}

	printRequestBody(r)
	fmt.Fprintf(w, "Token validated.\n")
}

func basicAuthHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != "admin" || password != "changeme" {
		http.Error(w, "Unauthorized: Invalid username or password", http.StatusUnauthorized)
		log.Printf("Unauthorized: Invalid username or password\n")
		return
	}
	printRequestBody(r)
	fmt.Fprintf(w, "Basic Auth validated.\n")
}

func headerValidationHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-SECRET")
	if secret == "" {
		http.Error(w, "Unauthorized: Missing X-SECRET", http.StatusUnauthorized)
		log.Printf("Unauthorized: Missing X-SECRET\n")
		return
	}
	printRequestBody(r)
	fmt.Fprintf(w, "Header validated.\n")
}

func main() {
	http.HandleFunc("/webhook/no-validation", noValidationHandler)
	http.HandleFunc("/webhook/token-validation", tokenValidationHandler)
	http.HandleFunc("/webhook/basic-auth", basicAuthHandler)
	http.HandleFunc("/webhook/header-validation", headerValidationHandler)

	log.Println("Server starting on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
