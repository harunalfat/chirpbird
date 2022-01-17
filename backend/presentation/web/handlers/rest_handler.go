package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/web"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
)

func jsonError(rw http.ResponseWriter, code int, errs ...error) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(code)

	var errorStrings []string
	for _, e := range errs {
		errorStrings = append(errorStrings, e.Error())
	}
	log.Printf("Giving error response\n%s", errorStrings)

	resp := web.Response{
		Data:   code,
		Errors: errorStrings,
	}
	json.NewEncoder(rw).Encode(resp)
}

func jsonResponse(rw http.ResponseWriter, code int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(code)
	resp := web.Response{
		Data: data,
	}
	json.NewEncoder(rw).Encode(resp)
}

type RestHandler struct {
	channelUseCase *usecases.ChannelUseCase
	messageUseCase *usecases.MessageUseCase
	userUseCase    *usecases.UserUseCase
}

func NewRestHandler(channelUseCase *usecases.ChannelUseCase, messageUseCase *usecases.MessageUseCase, userUseCase *usecases.UserUseCase) *RestHandler {
	return &RestHandler{
		channelUseCase,
		messageUseCase,
		userUseCase,
	}
}

func (handler *RestHandler) RegisterUser(rw http.ResponseWriter, r *http.Request) {
	var user entities.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jsonError(rw, http.StatusBadRequest, err)
		return
	}

	user, err := handler.userUseCase.CreateIfUsernameNotExist(r.Context(), user)
	if err != nil {
		jsonError(rw, http.StatusBadRequest, err)
		return
	}

	jsonResponse(rw, http.StatusCreated, user)
}
