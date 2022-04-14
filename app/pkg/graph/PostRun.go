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

func PostRun(formFile string) (map[string]string, error) {
	functionName := "PostRun"

	var urls = make(map[string]string)

	// Open form file
	diagramExport, openErr := os.Open(formFile)
	if openErr != nil {
		return urls, fmt.Errorf("%v: %w", functionName, openErr)
	}

	// Initialize multipart writer
	var uploadBody = &bytes.Buffer{}
	uploadWriter := multipart.NewWriter(uploadBody)

	// Populate form fields
	// file: "pdf"
	formWriter, formErr := uploadWriter.CreateFormFile("pdf", formFile)
	if formErr != nil {
		return urls, fmt.Errorf("%v: %w", functionName, formErr)
	}
	readAll, _ := io.ReadAll(diagramExport)
	formWriter.Write(readAll)

	// field: "source"
	formWriter, formErr = uploadWriter.CreateFormField("source")
	if formErr != nil {
		return urls, fmt.Errorf("%v: %w", functionName, formErr)
	}
	formWriter.Write([]byte("CI"))

	// Close multipart writer
	uploadWriter.Close()

	// Construct request
	request, _ := http.NewRequest("POST", "https://api.pluralith.com/v1/run/post", uploadBody)
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
		return urls, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Parse response
	responseBody, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return urls, fmt.Errorf("%v: %w", functionName, readErr)
	}

	if response.StatusCode != 200 {
		return urls, fmt.Errorf("%v: %w", functionName, readErr)
	}

	var bodyObject map[string]interface{}
	parseErr := json.Unmarshal(responseBody, &bodyObject)
	if parseErr != nil {
		return urls, fmt.Errorf("parsing response failed -> %v: %w", functionName, responseErr)
	}

	dataObject := bodyObject["data"].(map[string]interface{})

	urls["pluralithURL"] = dataObject["pluralithURL"].(string) // = bodyObject["data"].(structs.ExportURLs)
	urls["thumbnailURL"] = dataObject["pngURL"].(string)
	return urls, nil
}
