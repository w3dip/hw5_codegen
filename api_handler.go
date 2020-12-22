package main

import "net/http"
import "encoding/json"
import "io"
import "fmt"
import "strconv"
import "strings"

func (srv *MyApi) handleProfile(w http.ResponseWriter, r *http.Request) {

	// заполнение структуры params

	login := r.FormValue("login")

	if login == "" {
		makeOutput(w, ApiResponse{
			Error: "login must me not empty",
		}, http.StatusBadRequest)
		return
	}

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

	if login == "" {
		makeOutput(w, ApiResponse{
			Error: "login must me not empty",
		}, http.StatusBadRequest)
		return
	}

	if len(login) < 10 {
		makeOutput(w, ApiResponse{
			Error: "login len must be >= 10",
		}, http.StatusBadRequest)
		return
	}

	full_name := r.FormValue("full_name")

	status := r.FormValue("status")

	if status != "" {
		elem_values := make([]string, 0)

		elem_values = append(elem_values, "user")

		elem_values = append(elem_values, "moderator")

		elem_values = append(elem_values, "admin")

		found_elem := false
		for _, elem := range elem_values {
			if elem == status {
				found_elem = true
			}
		}
		if !found_elem {
			separated_str := strings.Join(elem_values, ", ")
			makeOutput(w, ApiResponse{
				Error: "status must be one of [" + separated_str + "]",
			}, http.StatusBadRequest)
			return
		}
	}

	if status == "" {
		status = "user"
	}

	age, convert_err := strconv.Atoi(r.FormValue("age"))
	if convert_err != nil {
		makeOutput(w, ApiResponse{
			Error: "age must be int",
		}, http.StatusBadRequest)
		return
	}

	if age < 0 {
		makeOutput(w, ApiResponse{
			Error: "age must be >= 0",
		}, http.StatusBadRequest)
		return
	}

	if age > 128 {
		makeOutput(w, ApiResponse{
			Error: "age must be <= 128",
		}, http.StatusBadRequest)
		return
	}

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

	if username == "" {
		makeOutput(w, ApiResponse{
			Error: "username must me not empty",
		}, http.StatusBadRequest)
		return
	}

	if len(username) < 3 {
		makeOutput(w, ApiResponse{
			Error: "username len must be >= 3",
		}, http.StatusBadRequest)
		return
	}

	account_name := r.FormValue("account_name")

	class := r.FormValue("class")

	if class != "" {
		elem_values := make([]string, 0)

		elem_values = append(elem_values, "warrior")

		elem_values = append(elem_values, "sorcerer")

		elem_values = append(elem_values, "rouge")

		found_elem := false
		for _, elem := range elem_values {
			if elem == class {
				found_elem = true
			}
		}
		if !found_elem {
			separated_str := strings.Join(elem_values, ", ")
			makeOutput(w, ApiResponse{
				Error: "class must be one of [" + separated_str + "]",
			}, http.StatusBadRequest)
			return
		}
	}

	if class == "" {
		class = "warrior"
	}

	level, convert_err := strconv.Atoi(r.FormValue("level"))
	if convert_err != nil {
		makeOutput(w, ApiResponse{
			Error: "level must be int",
		}, http.StatusBadRequest)
		return
	}

	if level < 1 {
		makeOutput(w, ApiResponse{
			Error: "level must be >= 1",
		}, http.StatusBadRequest)
		return
	}

	if level > 50 {
		makeOutput(w, ApiResponse{
			Error: "level must be <= 50",
		}, http.StatusBadRequest)
		return
	}

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
