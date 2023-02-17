package e2etest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/kumareswaramoorthi/companies/api/models"
	"github.com/stretchr/testify/require"
)

var token string
var testID string

func init() {
	client := &http.Client{}

	requestBody := []byte(`{"email": "admin@company.com", "password": "password"}`)
	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal(err)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	type tokenResponse struct {
		Token string `json:"token"`
	}
	var tokenResp tokenResponse
	err = json.Unmarshal(responseBody, &tokenResp)
	if err != nil {
		log.Fatal(err)
	}
	token = tokenResp.Token

	createTestData()
}

func TestCreateCompany(t *testing.T) {
	reqJson := `{
		"id": "041d2027-e6fa-4d6d-836d-eedb235c82bc",
		"name": "xyz6",
		"description": "new company",
		"amount_of_employees": 100,
		"type" : "Corporations",
		"registered": true
	}`

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/company", strings.NewReader(reqJson))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	var actualResponse models.Company
	err = json.Unmarshal(respBody, &actualResponse)
	require.Nil(t, err)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Equal(t, actualResponse.Name, "xyz6")
}

func TestGetCompany(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/company/"+testID, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	var actualResponse models.Company
	err = json.Unmarshal(respBody, &actualResponse)
	require.Nil(t, err)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, actualResponse.Name, "test")
}

func TestPatchCompany(t *testing.T) {
	reqJson := `{
		"name": "updated company",
		"description": "updated company description"
	}`

	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", "http://localhost:8080/api/v1/company/"+testID, strings.NewReader(reqJson))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	var actualResponse models.Company
	err = json.Unmarshal(respBody, &actualResponse)
	require.Nil(t, err)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, actualResponse.Name, "updated company")
	require.Equal(t, actualResponse.Description, "updated company description")
}

func TestDeleteCompany(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/api/v1/company/"+testID, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func createTestData() {
	reqJson := `{
		"id": "041d2027-e6fa-4d6d-836d-eedb235c82be",
		"name": "test",
		"description": "test company",
		"amount_of_employees": 100,
		"type" : "Corporations",
		"registered": true
		}`

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/company", strings.NewReader(reqJson))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var resp models.Company
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		log.Fatal(err)
	}
	testID = resp.ID
}
