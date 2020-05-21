package controller

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// docsController represents the documents controller implementation
type docsController struct {
	logger *log.Logger
	Spec   string
}

// NewDocsController returns a pointer to the new instance of docsController
func NewDocsController(logger *log.Logger, spec string) *docsController {
	return &docsController{
		logger: logger,
		Spec:   spec,
	}
}

// GetSpec handles getting swagger specification via http
func (c *docsController) GetSpec(w http.ResponseWriter, _ *http.Request) {
	err := WriteRawJSON([]byte(c.Spec), http.StatusOK, w)
	if err != nil {
		c.logger.Error(err)
		WriteError(err, http.StatusInternalServerError, w)
	}
}
