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

	"globe-and-citizen/layer8/proxy/resource_server/controller"
)

func TestPostgresConnection(t *testing.T) {
	config.InitDB()
	t.Log(*config.DB)
}

func TestIndex(t *testing.T) {
	r, err := http.NewRequest("GET", "http://localhost:5001/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rwr := httptest.NewRecorder()
	handler := http.HandlerFunc(Router)
	handler.ServeHTTP(rwr, r)

	status := rwr.Code
	if status != http.StatusOK {
		t.Errorf("Index was not available. Status code got <%v> wanted <%v>", status, http.StatusOK)
	}
}

func TestLoginUserHandler(t *testing.T) {
	precheckPayload := []byte(`{"username": "tester"}`)
	r, err := http.NewRequest("POST", "http://localhost:50001/api/v1/login-precheck", bytes.NewBuffer(precheckPayload))
	if err != nil {
		t.Fatal(err)
	}
	rwr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.LoginPrecheckHandler)
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
	handler2 := http.HandlerFunc(controller.LoginUserHandler)
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

func dep_TestDatabase(t *testing.T) {
	r, err := http.NewRequest("POST", "http://localhost:5001/login", nil)
	if err != nil {
		t.Fatal(err) // t.Log + T.FailNow()
	}

	rwr := httptest.NewRecorder()
	handler := http.HandlerFunc(Router)
	handler.ServeHTTP(rwr, r)

	t.Logf("%v", rwr)
}

// inline struct declaration
// pcResp := struct {
// 	username string
// 	salt string
// }{
// 	username: "value",
// 	salt: "Shark",
// }
