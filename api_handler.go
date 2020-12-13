package main

import "net/http"
import "encoding/json"
import "io"
import "fmt"
import "strconv"

func (srv *MyApi) handleProfile(w http.ResponseWriter, r *http.Request) {

	// заполнение структуры params

	login := r.FormValue("login")

	params := ProfileParams{

		Login: login,
	}
	// валидирование параметров
	ctx := r.Context()
	var res interface{}
	res, err := srv.Profile(ctx, params)
	// прочие обработки
	if err != nil {
		fmt.Printf("error happend: %+v\n", err)
		switch err.(type) {
		case ApiError:
			err := err.(ApiError)
			makeOutput(w, ApiResponse{
				Error: err.Err.Error(),
			}, err.HTTPStatus)
		default:
			makeOutput(w, ApiResponse{
				Error: err.Error(),
			}, http.StatusInternalServerError)
		}
		return
	}
	makeOutput(w, ApiResponse{
		Response: &res,
	}, http.StatusOK)
}

func (srv *MyApi) handleCreate(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("X-Auth") != "100500" {
		makeOutput(w, ApiResponse{
			Error: "unauthorized",
		}, http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		makeOutput(w, ApiResponse{
			Error: "bad method",
		}, http.StatusNotAcceptable)
		return
	}

	// заполнение структуры params

	login := r.FormValue("login")

	full_name := r.FormValue("full_name")

	status := r.FormValue("status")

	age, _ := strconv.Atoi(r.FormValue("age"))

	params := CreateParams{

		Login: login,

		Name: full_name,

		Status: status,

		Age: age,
	}
	// валидирование параметров
	ctx := r.Context()
	var res interface{}
	res, err := srv.Create(ctx, params)
	// прочие обработки
	if err != nil {
		fmt.Printf("error happend: %+v\n", err)
		switch err.(type) {
		case ApiError:
			err := err.(ApiError)
			makeOutput(w, ApiResponse{
				Error: err.Err.Error(),
			}, err.HTTPStatus)
		default:
			makeOutput(w, ApiResponse{
				Error: err.Error(),
			}, http.StatusInternalServerError)
		}
		return
	}
	makeOutput(w, ApiResponse{
		Response: &res,
	}, http.StatusOK)
}

func (srv *OtherApi) handleCreate(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("X-Auth") != "100500" {
		makeOutput(w, ApiResponse{
			Error: "unauthorized",
		}, http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		makeOutput(w, ApiResponse{
			Error: "bad method",
		}, http.StatusNotAcceptable)
		return
	}

	// заполнение структуры params

	username := r.FormValue("username")

	account_name := r.FormValue("account_name")

	class := r.FormValue("class")

	level, _ := strconv.Atoi(r.FormValue("level"))

	params := OtherCreateParams{

		Username: username,

		Name: account_name,

		Class: class,

		Level: level,
	}
	// валидирование параметров
	ctx := r.Context()
	var res interface{}
	res, err := srv.Create(ctx, params)
	// прочие обработки
	if err != nil {
		fmt.Printf("error happend: %+v\n", err)
		switch err.(type) {
		case ApiError:
			err := err.(ApiError)
			makeOutput(w, ApiResponse{
				Error: err.Err.Error(),
			}, err.HTTPStatus)
		default:
			makeOutput(w, ApiResponse{
				Error: err.Error(),
			}, http.StatusInternalServerError)
		}
		return
	}
	makeOutput(w, ApiResponse{
		Response: &res,
	}, http.StatusOK)
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
