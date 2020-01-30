package controller

import (
	"net/http"
)

type DocsController struct {
	Spec string
}

func (c *DocsController) GetSpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)(c.Spec))
}
