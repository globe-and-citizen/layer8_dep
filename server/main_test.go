package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/config"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	ctl "globe-and-citizen/layer8/server/resource_server/controller"
	"globe-and-citizen/layer8/server/resource_server/utils"
)

func TestPostgresConnection(t *testing.T) {
	config.InitDB()
	if config.DB.Error != nil {
		t.Fatalf("error connecting the db: %s", config.DB.Error.Error())
	}
	t.Log("Database successfully connected")
}

func TestLoginUserHandler(t *testing.T) {
	precheckPayload := []byte(`{"username": "tester"}`)
	r, err := http.NewRequest("POST", "http://localhost:50001/api/v1/login-precheck", bytes.NewBuffer(precheckPayload))
	if err != nil {
		t.Fatal(err)
	}
	rwr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctl.LoginPrecheckHandler)
	handler.ServeHTTP(rwr, r)

	t.Logf("%v", rwr.Body)
	pcResp := struct {
		Username string `json:username`
		Salt     string `json:salt`
	}{}
	resopnseBytes, err := io.ReadAll(rwr.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(resopnseBytes, &pcResp)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", pcResp)
	// Test 2 //
	payloadString := fmt.Sprintf(`{"username": "tester", "password": "12341234", "salt": "%s"}`, pcResp.Salt)
	loginPayload := []byte(payloadString)
	r2, err := http.NewRequest("POST", "http://localhost:50001/api/v1/login-user", bytes.NewBuffer(loginPayload))
	if err != nil {
		t.Fatal(err)
	}
	rwr2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(ctl.LoginUserHandler)
	handler2.ServeHTTP(rwr2, r2)

	status := rwr2.Code
	body2, err := io.ReadAll(rwr2.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body2))
	if status != http.StatusOK {
		t.Log("Bro ur nearly there...", status)
	}
}

func TestRegisterClientHandler(t *testing.T) {
	// Init the DB Connection
	config.InitDB()
	// Create a mock request with Name & RedirectURI as the body.
	uuid_obj := uuid.New()
	uuid_str := uuid_obj.String()
	jsonBodyAsString := fmt.Sprintf(`{"name":"new-client-%s","redirect_uri":"anc.com"}`, uuid_str[0:4])
	reqBody := []byte(jsonBodyAsString)
	reqBodyReader := bytes.NewBuffer(reqBody)
	req, err := http.NewRequest("POST", "http://localhost:50001/api/v1/login-user", reqBodyReader)
	if err != nil {
		t.Fatalf("error creating mock http request: %s", err.Error())
	}

	// Use an adapter to turn the function under test into a http.Handler
	rwr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctl.RegisterClientHandler)
	handler.ServeHTTP(rwr, req)

	if rwr.Code != 200 {
		t.Fatalf("req failed. rwr.Code received was: %d", rwr.Code)
	}

	rawResponse, err := io.ReadAll(rwr.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %s", err.Error())
	}

	response := &utils.Response{}
	if err := json.Unmarshal(rawResponse, response); err != nil {
		t.Fatalf("unable to unmarshal rawData: %s", err.Error())
	}
	t.Log(response.Status)
	t.Log(response.Error)
	t.Log(response.Data)
	t.Log(response.Message)
}

func TestIndex(t *testing.T) { // Router no longer exists
	r, err := http.NewRequest("GET", "http://localhost:5001/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rwr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctl.IndexHandler)
	handler.ServeHTTP(rwr, r)

	status := rwr.Code
	if status != http.StatusOK {
		t.Errorf("Index was not available. Status code got <%v> wanted <%v>", status, http.StatusOK)
	}
}
