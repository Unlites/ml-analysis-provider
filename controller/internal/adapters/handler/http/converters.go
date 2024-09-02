package httphandler

import (
	"github.com/Unlites/ml-analysis-provider/controller/internal/domain"
)

// toDomainAnalyzesFilter returnes AnalyzesFilter from GetAnalyzesParams
func toDomainAnalyzesFilter(params GetAnalyzesParams) domain.AnalyzesFilter {
	var filter domain.AnalyzesFilter

	if params.Limit == nil {
		filter.Limit = 10
	}

	if params.Offset == nil {
		filter.Offset = 0
	}

	if params.Query != nil {
		filter.Query = *params.Query
	}

	if params.Answer != nil {
		filter.Answer = *params.Answer
	}

	filter.IsUserSatisfied = params.IsUserSatisfied

	return filter
}

// toDomainAnalysis returnes Analysis from AnalysisRequest
func toDomainAnalysis(req AnalysisRequest) domain.Analysis {
	return domain.Analysis{
		Query:           req.Query,
		Answer:          req.Answer,
		IsUserSatisfied: req.IsUserSatisfied,
	}
}

// toAnalysisResponse returnes AnalysisResponse from Analysis
func toAnalysisResponse(analysis domain.Analysis) AnalysisResponse {
	return AnalysisResponse{
		Id:              analysis.Id,
		Query:           analysis.Query,
		Answer:          analysis.Answer,
		IsUserSatisfied: analysis.IsUserSatisfied,
	}
}

// toAnalysisResponses returnes slice of AnalysisResponse from given slice of Analysis
func toAnalysisResponses(analysis []domain.Analysis) []AnalysisResponse {
	responses := make([]AnalysisResponse, len(analysis))
	for i, a := range analysis {
		responses[i] = toAnalysisResponse(a)
	}
	return responses
}
