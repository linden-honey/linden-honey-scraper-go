package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func writeJSON(data interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		err := errors.Wrap(err, "Error happend during data marshalling")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(([]byte)(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}
