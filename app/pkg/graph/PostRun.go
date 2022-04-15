package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"pluralith/pkg/auxiliary"
)

func PostRun(formFile string) (map[string]interface{}, error) {
	functionName := "PostRun"

	// Open form file
	diagramExport, openErr := os.Open(formFile)
	if openErr != nil {
		return nil, fmt.Errorf("%v: %w", functionName, openErr)
	}

	// Initialize multipart writer
	var uploadBody = &bytes.Buffer{}
	uploadWriter := multipart.NewWriter(uploadBody)

	// Populate form fields
	// file: "pdf"
	formWriter, formErr := uploadWriter.CreateFormFile("pdf", formFile)
	if formErr != nil {
		return nil, fmt.Errorf("%v: %w", functionName, formErr)
	}
	readAll, _ := io.ReadAll(diagramExport)
	formWriter.Write(readAll)

	// Close multipart writer
	uploadWriter.Close()

	// Construct request
	request, _ := http.NewRequest("POST", "https://api.pluralith.com/v1/run/post", uploadBody)
	// request, _ := http.NewRequest("POST", "http://localhost:8080/v1/run/post", uploadBody)
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)
	request.Header.Add("Content-Type", uploadWriter.FormDataContentType())

	// Add project id query string
	queryString := request.URL.Query()
	queryString.Add("projectId", auxiliary.StateInstance.PluralithConfig.ProjectId)
	request.URL.RawQuery = queryString.Encode()

	// Instantiate client and execute request
	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Parse response
	responseBody, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, fmt.Errorf("%v: %w", functionName, readErr)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("%v: %w", functionName, readErr)
	}

	var bodyObject map[string]interface{}
	parseErr := json.Unmarshal(responseBody, &bodyObject)
	if parseErr != nil {
		return nil, fmt.Errorf("parsing response failed -> %v: %w", functionName, responseErr)
	}

	dataObject := bodyObject["data"].(map[string]interface{})

	return dataObject, nil
}
