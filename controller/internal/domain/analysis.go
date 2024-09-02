package domain

// Analysis represents a result of working with a ML model
type Analysis struct {
	Id              int
	Query           string
	Answer          string
	IsUserSatisfied bool
}

// AnalyzesFilter is a filter for getting analyzes
type AnalyzesFilter struct {
	Query           string
	Answer          string
	IsUserSatisfied *bool
	Limit           int
	Offset          int
}
