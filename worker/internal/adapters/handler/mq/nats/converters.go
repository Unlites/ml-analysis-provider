package nats

import "github.com/Unlites/ml-analysis-provider/worker/internal/domain"

// toDomainAnalysis returnes Analysis from AnalysisNats
func toDomainAnalysis(an AnalysisNats) domain.Analysis {
	return domain.Analysis{
		Id:              an.Id,
		Query:           an.Query,
		Answer:          an.Answer,
		IsUserSatisfied: an.IsUserSatisfied,
	}
}

// toAnalysisNats returnes AnalysisNats from Analysis
func toAnalysisNats(ad domain.Analysis) AnalysisNats {
	return AnalysisNats{
		Id:              ad.Id,
		Query:           ad.Query,
		Answer:          ad.Answer,
		IsUserSatisfied: ad.IsUserSatisfied,
	}
}

// toDomainAnalyzesFilter returnes AnalyzesFilter from AnalyzesFilterNats
func toDomainAnalyzesFilter(afn AnalyzesFilterNats) domain.AnalyzesFilter {
	return domain.AnalyzesFilter{
		Query:           afn.Query,
		Answer:          afn.Answer,
		IsUserSatisfied: afn.IsUserSatisfied,
		Limit:           afn.Limit,
		Offset:          afn.Offset,
	}
}

// toAnalyzesDomain returnes slice of AnalysisNats from given slice of Analysis
func toAnalyzesNats(analyzes []domain.Analysis) []AnalysisNats {
	responses := make([]AnalysisNats, len(analyzes))
	for i, a := range analyzes {
		responses[i] = toAnalysisNats(a)
	}
	return responses
}
