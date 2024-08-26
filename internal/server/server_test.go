package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"simple-service/tools/storage"
	"testing"
)

type JSONData map[string]interface{}

func makeRequest(resource string, method string, params ...JSONData) (JSONData, int) {
	result := make(JSONData)
	URL, err := url.Parse(resource)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest(method, URL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")

	if len(params) > 0 {
		jsonData, _ := json.Marshal(params[0])
		reader := bytes.NewReader(jsonData)
		request.Body = io.NopCloser(reader)
	}

	response := httptest.NewRecorder()
	hendler := GetAPIHendler()
	hendler.ServeHTTP(response, request)

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result, response.Code
}

func TestList(t *testing.T) {
	resource := "/data/list"
	data, code := makeRequest(resource, http.MethodGet)

	if code != http.StatusOK {
		t.Error("invalid response code")
	}

	_, ok := data["result"]

	if !ok {
		t.Error("invalid response structure")
	}
}

func TestCreateSelect(t *testing.T) {
	resource := "/data/create"
	createRecord := storage.Telnum{
		Msisdn:     "79201112233",
		Region:     "msk",
		Abc:        "74951002030",
		Enabled:    true,
		ServiceKey: 3507,
	}

	bytes, _ := json.Marshal(&createRecord)
	params := make(JSONData)
	json.Unmarshal(bytes, &params)

	data, code := makeRequest(resource, http.MethodPut, params)
	key := data["key"].(string)

	if code != http.StatusCreated {
		t.Error("invalid response code")
	}

	//Проверяем содержимое ответа
	if key != storage.MakeHash(createRecord) {
		t.Error("invalid create key")
	}

	resource = fmt.Sprintf("/data/item/%s", key)
	data, code = makeRequest(resource, http.MethodGet)

	selectRecord := storage.Telnum{}
	itemBytes, _ := json.Marshal(data["item"])
	_ = json.Unmarshal(itemBytes, &selectRecord)

	if code != http.StatusOK {
		t.Error("invalid response code")
	}

	//Проверяем содержимое ответа
	if key != storage.MakeHash(selectRecord) {
		t.Error("invalid select key")
	}
}

func TestCreateUpdate(t *testing.T) {
	resource := "/data/create"
	createRecord := storage.Telnum{
		Msisdn:     "79201112233",
		Region:     "msk",
		Abc:        "74951002030",
		Enabled:    true,
		ServiceKey: 3507,
	}

	bytes, _ := json.Marshal(&createRecord)
	params := make(JSONData)
	json.Unmarshal(bytes, &params)

	data, code := makeRequest(resource, http.MethodPut, params)
	key := data["key"].(string)

	if code != http.StatusCreated {
		t.Error("invalid response code")
	}

	//Проверяем содержимое ответа
	if key != storage.MakeHash(createRecord) {
		t.Error("invalid create key")
	}

	updateRecord := storage.Telnum{
		Msisdn:     "79209201001",
		Region:     "msk",
		Abc:        "79209201002",
		Enabled:    true,
		ServiceKey: 3507,
	}

	bytes, _ = json.Marshal(&updateRecord)
	params = make(JSONData)
	json.Unmarshal(bytes, &params)

	resource = fmt.Sprintf("/data/update/%s", key)
	data, code = makeRequest(resource, http.MethodPost, params)
	key = data["key"].(string)

	if code != http.StatusOK {
		t.Error("invalid response code")
	}

	//Проверяем содержимое ответа
	if key != storage.MakeHash(updateRecord) {
		t.Error("invalid update key")
	}
}

func TestCreateDelete(t *testing.T) {
	resource := "/data/create"
	createRecord := storage.Telnum{
		Msisdn:     "79201112233",
		Region:     "msk",
		Abc:        "74951002030",
		Enabled:    true,
		ServiceKey: 3507,
	}

	bytes, _ := json.Marshal(&createRecord)
	params := make(JSONData)
	json.Unmarshal(bytes, &params)

	data, code := makeRequest(resource, http.MethodPut, params)
	key := data["key"].(string)

	if code != http.StatusCreated {
		t.Error("invalid response code")
	}

	//Проверяем содержимое ответа
	if key != storage.MakeHash(createRecord) {
		t.Error("invalid create key")
	}

	resource = fmt.Sprintf("/data/delete/%s", key)
	data, code = makeRequest(resource, http.MethodDelete)
	result := data["result"].(bool)

	if code != http.StatusOK {
		t.Error("invalid response code")
	}

	//Проверяем содержимое ответа
	if result != true {
		t.Error("invalid result delete")
	}
}
