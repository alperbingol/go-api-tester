package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <METHOD> <URL> [BODY] [HEADERS]")
		fmt.Println("Example GET: go run main.go GET https://jsonplaceholder.typicode.com/todos/1")
		fmt.Println("Example POST: go run main.go POST https://jsonplaceholder.typicode.com/posts '{\"title\": \"foo\", \"body\": \"bar\", \"userId\": 1}' \"Content-Type:application/json\"")
		return
	}

	method := os.Args[1] // HTTP method (GET, POST, PUT, DELETE)
	url := os.Args[2]    // API URL
	var jsonBody []byte  // Request body (for POST, PUT)

	// Handle request body (for POST and PUT)
	if (method == "POST" || method == "PUT") && len(os.Args) > 3 {
		jsonBody = []byte(os.Args[3])
	}

	// Create request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add headers (if provided)
	if len(os.Args) > 4 {
		for _, header := range os.Args[4:] {
			parts := bytes.SplitN([]byte(header), []byte(":"), 2) // Only split at first `:`
			if len(parts) == 2 {
				req.Header.Add(string(parts[0]), string(parts[1]))
			}
		}
	}

	// Debug: Print request headers to verify they are being set
	fmt.Println("\nRequest Headers Before Sending:")
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading resposnse:", err)
		return
	}

	// Pretty-print JSON response
	fmt.Println("Status Code:", resp.Status)

	// Print response headers
	fmt.Println("\nRequest Headers:")
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// Pretty-print JSON response
	prettyPrintJSON(body)

}

// Pretty-print JSON function
func prettyPrintJSON(body []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		fmt.Println("Response Body:\n", string(body)) // Print as-is if not JSON
	} else {
		fmt.Println("Response Body:\n", prettyJSON.String())
	}
}
