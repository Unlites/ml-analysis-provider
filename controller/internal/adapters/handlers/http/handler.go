package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Unlites/ml-analysis-provider/controller/internal/application"
	"github.com/Unlites/ml-analysis-provider/controller/internal/domain"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

// HTTPHandler is the http handler adapter for application
type HTTPHandler struct {
	uc application.Usecase
}

// NewHTTPHandler creates a new http handler adapter for application
func NewHTTPHandler(uc application.Usecase) (http.Handler, error) {
	swagger, err := GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load swagger spec: %w", err)
	}

	swagger.Servers = nil
	impl := &HTTPHandler{uc: uc}

	handler := Handler(impl)
	handler = middleware.OapiRequestValidator(swagger)(handler)

	return http.StripPrefix("/api/v1", handler), nil
}

// GetAnalyzes implements ServerInterface
func (h *HTTPHandler) GetAnalyzes(w http.ResponseWriter, r *http.Request, params GetAnalyzesParams) {
	analyzes, err := h.uc.GetAnalyzes(r.Context(), h.toDomainAnalyzesFilter(params))
	if err != nil {
		h.encodeResponse(
			w, http.StatusInternalServerError,
			ErrorResponse{"failed to get analyzes: " + err.Error()},
		)
		return
	}

	h.encodeResponse(w, http.StatusOK, h.toAnalysisResponses(analyzes))
}

// AddAnalysis implements ServerInterface
func (h *HTTPHandler) AddAnalysis(w http.ResponseWriter, r *http.Request) {
	var req AnalysisRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.encodeResponse(
			w, http.StatusBadRequest,
			ErrorResponse{"invalid request: " + err.Error()},
		)
		return
	}

	if err := h.uc.AddAnalysis(r.Context(), h.toDomainAnalysis(req)); err != nil {
		h.encodeResponse(
			w, http.StatusInternalServerError,
			ErrorResponse{"failed to add analysis: " + err.Error()},
		)
		return
	}

	h.encodeResponse(w, http.StatusCreated, nil)
}

// GetAnalyzesId implements ServerInterface
func (h *HTTPHandler) GetAnalysisById(w http.ResponseWriter, r *http.Request, id string) {
	analysis, err := h.uc.GetAnalysisById(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if err == domain.ErrNotFound {
			status = http.StatusNotFound
		}

		h.encodeResponse(
			w, status,
			ErrorResponse{"failed to get analysis: " + err.Error()},
		)
		return
	}

	h.encodeResponse(w, http.StatusOK, h.toAnalysisResponse(analysis))
}
