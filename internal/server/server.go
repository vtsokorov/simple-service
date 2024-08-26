package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"simple-service/tools/storage"
	"strings"
	"sync"
)

var (
	template   = regexp.MustCompile(`^\/data(\/.*)?$`)
	apiHendler *EndpointHendler
	once       sync.Once
	messages   map[int]string = map[int]string{
		http.StatusNotFound:            "record not found",
		http.StatusInternalServerError: "internal server error",
		http.StatusBadRequest:          "bad request",
	}
)

func splitPath(path string) []string {
	matches := template.FindStringSubmatch(path)

	action := "undefine"
	if len(matches) > 1 {
		action = matches[1]
	}
	return strings.Split(action, "/")
}

type EndpointHendler struct {
	cache storage.Cache
}

func (e *EndpointHendler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	items := splitPath(r.URL.Path)

	switch {
	case r.Method == http.MethodGet && items[1] == "list":
		e.list(w, r)
		return
	case r.Method == http.MethodPut && items[1] == "create":
		e.create(w, r)
		return
	case r.Method == http.MethodGet && items[1] == "item":
		e.item(w, r)
		return
	case r.Method == http.MethodDelete && items[1] == "delete":
		e.delete(w, r)
		return
	case r.Method == http.MethodPost && items[1] == "update":
		e.update(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(messages[http.StatusNotFound]))
		return
	}
}

func (e *EndpointHendler) list(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{}
	data := e.cache.Show()

	w.WriteHeader(http.StatusOK)
	result["result"] = data

	json.NewEncoder(w).Encode(result)
}

func (e *EndpointHendler) item(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{}
	items := splitPath(r.URL.Path)

	if len(items) >= 3 {
		data, ok := e.cache.Get(items[2])
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			result["error"] = messages[http.StatusNotFound]
		} else {
			w.WriteHeader(http.StatusOK)
			result["item"] = data
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		result["error"] = messages[http.StatusInternalServerError]
	}

	json.NewEncoder(w).Encode(result)
}

func (v *EndpointHendler) create(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{}

	var request storage.Telnum
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result["error"] = messages[http.StatusBadRequest]
	} else {
		key := v.cache.Create(request)

		w.WriteHeader(http.StatusCreated)
		result["result"] = true
		result["key"] = key
	}

	json.NewEncoder(w).Encode(result)
}

func (e *EndpointHendler) delete(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{}
	items := splitPath(r.URL.Path)

	if len(items) >= 3 {
		ok := e.cache.Delete(items[2])
		result["result"] = ok

		if ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
			result["error"] = messages[http.StatusNotFound]
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		result["error"] = messages[http.StatusInternalServerError]
	}

	json.NewEncoder(w).Encode(result)
}

func (e *EndpointHendler) update(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{}
	items := splitPath(r.URL.Path)

	if len(items) >= 3 {
		var request storage.Telnum
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			result["error"] = messages[http.StatusBadRequest]
		} else {
			key, ok := e.cache.Update(items[2], request)
			result["result"] = ok

			if ok {
				w.WriteHeader(http.StatusOK)
				result["key"] = key
			} else {
				w.WriteHeader(http.StatusNotFound)
				result["error"] = messages[http.StatusNotFound]
			}
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		result["error"] = messages[http.StatusInternalServerError]
	}

	json.NewEncoder(w).Encode(result)
}

func createObject() *EndpointHendler {
	return &EndpointHendler{
		cache: storage.CreateCacheObject(),
	}
}

func GetAPIHendler() *EndpointHendler {
	once.Do(func() {
		apiHendler = createObject()
	})
	return apiHendler
}

func RunServer(port int) {
	mux := http.NewServeMux()

	hendler := GetAPIHendler()
	mux.Handle("/data/", hendler)

	log.Printf("server is listening port :%d...", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
