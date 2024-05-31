package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	//	Setup flags and parse them
	filePath := flag.String("file", "", "the file to upload")
	targetURL := flag.String("url", "https://packassist.cagedtornado.com/v1/package", "the package-assist endpoint")
	authToken := flag.String("token", "", "the auth token")
	flag.Parse()

	//	If filePath is empty, just get out:
	if *filePath == "" {
		fmt.Println("Missing required argument: -file\n")
		flag.Usage()
		return
	}

	// Open the file
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to write our multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field and write the file to it
	part, err := writer.CreateFormFile("file", *filePath)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file to form file:", err)
		return
	}

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", *targetURL, &requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the Content-Type header to multipart/form-data with the boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set the X-PackAuth header with the authorization token.  This must match
	// what the targetURL is expecting
	req.Header.Set("X-PackAuth", *authToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
