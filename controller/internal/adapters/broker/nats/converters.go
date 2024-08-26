package broker

import "github.com/Unlites/ml-analysis-provider/controller/internal/domain"

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

// toAnalyzesFilterNats returnes AnalyzesFilterNats from AnalyzesFilter
func toAnalyzesFilterNats(ad domain.AnalyzesFilter) AnalyzesFilterNats {
	return AnalyzesFilterNats{
		Query:           ad.Query,
		Answer:          ad.Answer,
		IsUserSatisfied: ad.IsUserSatisfied,
	}
}
