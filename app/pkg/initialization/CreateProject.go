package initialization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func CreateProject(initData InitData) error {
	functionName := "CreateProject"

	creationSpinner := ux.NewSpinner("Creating Project", "Project Created Successfully", "Project Creation Failed, Try Again!", true)
	creationSpinner.Start()

	projectData := make(map[string]string)
	projectData["id"] = initData.ProjectId
	projectData["orgId"] = initData.OrgId
	projectData["name"] = initData.ProjectName

	// Stringify run cache
	projectDataBytes, marshalErr := json.MarshalIndent(projectData, "", " ")
	if marshalErr != nil {
		creationSpinner.Fail()
		return fmt.Errorf("creating project json string failed -> %v: %w", functionName, marshalErr)
	}

	// request, _ := http.NewRequest("POST", "https://api.pluralith.com/v1/project/create", bytes.NewBuffer(projectDataBytes))
	request, _ := http.NewRequest("POST", "http://localhost:8080/v1/project/create", bytes.NewBuffer(projectDataBytes))
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		creationSpinner.Fail()
		return fmt.Errorf("creating project failed -> %v: %w", functionName, responseErr)
	}

	// Parse response for file URLs
	responseBody, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		creationSpinner.Fail()
		return fmt.Errorf("%v: %w", functionName, readErr)
	}

	var bodyObject map[string]interface{}
	parseErr := json.Unmarshal(responseBody, &bodyObject)
	if parseErr != nil {
		creationSpinner.Fail()
		return fmt.Errorf("parsing response failed -> %v: %w", functionName, parseErr)
	}

	if response.StatusCode >= 300 {
		creationSpinner.Fail()
		return fmt.Errorf("request failed -> %v: %v", functionName, bodyObject)
	}

	creationSpinner.Success()
	return nil
}
