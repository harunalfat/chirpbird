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
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(code)

	var errorStrings []string
	for _, e := range errs {
		errorStrings = append(errorStrings, e.Error())
	}

	resp := web.Response{
		Data:   code,
		Errors: errorStrings,
	}
	json.NewEncoder(rw).Encode(resp)
}

func jsonResponse(rw http.ResponseWriter, code, data interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp := web.Response{
		Data: data,
	}
	json.NewEncoder(rw).Encode(resp)
}

type RestHandler struct {
	channelUseCase usecases.ChannelUseCase
	userUseCase    usecases.UserUseCase
}

func NewRestHandlerImpl(channelUseCase usecases.ChannelUseCase, userUseCase usecases.UserUseCase) *RestHandler {
	return &RestHandler{
		channelUseCase,
		userUseCase,
	}
}

func (handler *RestHandler) CreateChannel(rw http.ResponseWriter, r *http.Request) {
	var channel entities.Channel
	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		jsonError(rw, http.StatusBadRequest, err)
		return
	}

	userID := r.URL.Query().Get("userId")
	creator, err := handler.userUseCase.Fetch(r.Context(), userID)
	if err != nil {
		jsonError(rw, http.StatusBadRequest, err)
		return
	}

	inserted, err := handler.channelUseCase.Create(r.Context(), channel, creator)
	if err != nil {
		jsonError(rw, http.StatusBadRequest, err)
		return
	}

	if _, err = handler.userUseCase.EmbedChannelIfNotExist(r.Context(), creator, inserted); err != nil {
		jsonError(rw, http.StatusBadRequest, err)
		return
	}

	jsonResponse(rw, http.StatusCreated, map[string]interface{}{"channel": inserted})
}

func (handler *RestHandler) InviteToChannel(rw http.ResponseWriter, r *http.Request) {
	var payload entities.InvitePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		jsonError(rw, http.StatusBadRequest, err)
		return
	}

	if err = handler.userUseCase.EmbedChannelToMultipleUsersIfNotExist(r.Context(), payload.UserIDs, payload.ChannelID); err != nil {
		jsonError(rw, http.StatusInternalServerError, err)
		return
	}

	for _, username := range payload.UserIDs {
		if err = handler.userUseCase.SubsribeUserConnectionToChannel(r.Context(), username, payload.ChannelID); err != nil {
			log.Printf("Failed to subsribe client [%s] to channel [%s]", username, payload.ChannelID)
		}
	}
	jsonResponse(rw, http.StatusOK, "ok")
}
