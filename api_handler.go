package main

import "net/http"
import "encoding/json"
import "io"

func (srv *MyApi) handleProfile(w http.ResponseWriter, r *http.Request) {
}

func (srv *MyApi) handleCreate(w http.ResponseWriter, r *http.Request) {
}

func (srv *OtherApi) handleCreate(w http.ResponseWriter, r *http.Request) {
}

func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/user/profile":
		srv.handleProfile(w, r)

	case "/user/create":
		srv.handleCreate(w, r)

	default:
		makeOutput(w, ApiResponse{
			Error: "unknown method",
		}, http.StatusNotFound)
	}
}

func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/user/create":
		srv.handleCreate(w, r)

	default:
		makeOutput(w, ApiResponse{
			Error: "unknown method",
		}, http.StatusNotFound)
	}
}

type ApiResponse struct {
	Error    string       `json:"error"`
	Response *interface{} `json:"response,omitempty"`
}

func makeOutput(w http.ResponseWriter, body interface{}, status int) {
	w.WriteHeader(status)
	result, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	_, err_write := io.WriteString(w, string(result))
	if err_write != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
