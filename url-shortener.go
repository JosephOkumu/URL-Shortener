package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// ShortURLs is a map that stores shortened URLs with their corresponding original URLs.
var shortURLs = make(map[string]string)

func main() {
	// Set up HTTP request handlers
	http.HandleFunc("/", handleForm)           //Handle root endpoint
	http.HandleFunc("/shorten", handleShorten) // Handle URL shortening endpoint
	http.HandleFunc("/short/", handleRedirect) // Handle redirecting to original URL

	// Start HTTP server
	fmt.Println("URL Shortener running on: 3030")
	http.ListenAndServe(":3030", nil)
}

// handleForm handles GET requests to the root endpoint and displays the URL shortening form
func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title> URL-Shortener</title>
		</head>
		<body>
			<h2>URL-Shortener</h2>
			<form method= "post" action="/shorten">
				<input type="url" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form> 
		</body>
		</html>
	`)
}

// handleShorten handles POST requests to the /shorten endpoint and generates a shortened URL
func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the original URL from the form
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Generate a short key and store the original URL
	shortKey := generateShortKey()
	shortURLs[shortKey] = originalURL

	// Construct the shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:3030/short/%s", shortKey)

	// Render HTML response with original and shortened URLs
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title> URL-Shortener</title>
		</head>
		<body>
			<h2>URL-Shortener</h2>
			<p>Original URL: `, originalURL, `</p>
			<p>Shortened URL: <a href="`, shortenedURL, `">`, shortenedURL, `</a></p>
		</body>
		</html>
		`)
}

// handleRedirect handles requests to shortened URLs and redirects to the original URL.
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Extract the short key from the request path
	shortKey := strings.TrimPrefix(r.URL.Path, "/short/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the map
	originalURL, found := shortURLs[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	// Redirect to the orginal URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// generateShortKey generates a short random key for URL shortening
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
