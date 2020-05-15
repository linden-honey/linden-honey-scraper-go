package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func WriteJSON(data interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		err := errors.Wrap(err, "Error happened during data marshalling")
		log.Println(err)
		WriteError(err, status, w)
	} else {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonBytes)
		if err != nil {
			log.Println(errors.Wrap(err, "Error happened during response writing"))
		}
	}
}

func WriteError(err error, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	_, err = w.Write(([]byte)(err.Error()))
	if err != nil {
		log.Println(errors.Wrap(err, "Error happened during response writing"))
	}
}
