package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiResponse struct {
	Error    string       `json:"error"`
	Response *interface{} `json:"response,omitempty"`
}

func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		srv.handlerProfile(w, r)
	default:
		makeOutput(w, ApiResponse{
			Error: "unknown method",
		}, http.StatusNotFound)
	}
}

func (srv *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	//requiredMethod := http.MethodGet
	//if r.Method != requiredMethod {
	//	http.Error(w, "bad method", http.StatusNotAcceptable)
	//	return
	//}
	// заполнение структуры params
	params := ProfileParams{
		Login: r.FormValue("login"),
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

//func (params *ProfileParams) validate()
