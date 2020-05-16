package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func WriteRawJSON(data []byte, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		log.Println(errors.Wrap(err, "Error happened during response writing"))
	}
}

func WriteJSON(data interface{}, status int, w http.ResponseWriter) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		err := errors.Wrap(err, "Error happened during data marshalling")
		log.Println(err)
		WriteError(err, http.StatusInternalServerError, w)
	} else {
		WriteRawJSON(dataBytes, status, w)
	}
}

func WriteError(err error, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	_, err = w.Write(([]byte)(err.Error()))
	if err != nil {
		log.Println(errors.Wrap(err, "Error happened during response writing"))
	}
}
