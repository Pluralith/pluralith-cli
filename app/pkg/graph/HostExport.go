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

type ExportURLs struct {
	pdf string
	png string
}

func HostExport(formFile string) (ExportURLs, error) {
	functionName := "LogExport"

	var urls = ExportURLs{}

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

	// field: "id"
	formWriter, formErr = uploadWriter.CreateFormField("id")
	if formErr != nil {
		return urls, fmt.Errorf("%v: %w", functionName, formErr)
	}
	formWriter.Write([]byte("hant"))

	// field: "source"
	formWriter, formErr = uploadWriter.CreateFormField("source")
	if formErr != nil {
		return urls, fmt.Errorf("%v: %w", functionName, formErr)
	}
	formWriter.Write([]byte("CLI"))

	// Close multipart writer
	uploadWriter.Close()

	// Construct request
	request, _ := http.NewRequest("POST", "https://pluralith-api-dev-r3vxinfd2a-lm.a.run.app/v1/user/export/publish", uploadBody)
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)
	request.Header.Add("Content-Type", uploadWriter.FormDataContentType())

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

	urls = bodyObject["data"].(ExportURLs)
	return urls, nil
}
