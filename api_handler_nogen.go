package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiResponse struct {
	Error    string `json:"error"`
	Response User   `json:"response"`
}

func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		srv.handlerProfile(w, r)
	default:
		// 404
	}
}

func (srv *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	// заполнение структуры params
	params := ProfileParams{
		Login: r.FormValue("login"),
	}
	// валидирование параметров
	ctx := r.Context()
	res, err := srv.Profile(ctx, params)
	// прочие обработки
	if err != nil {
		fmt.Printf("error happend: %+v\n", err)
		switch err.(type) {
		case *ApiError:
			err := err.(*ApiError)
			http.Error(w, err.Err.Error(), err.HTTPStatus)
		default:
			http.Error(w, "internal error", 500)
		}
		return
	}

	result, err := json.Marshal(ApiResponse{
		Response: *res,
	})
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	_, err_write := io.WriteString(w, string(result))
	if err_write != nil {
		http.Error(w, "internal error", 500)
		return
	}
}
