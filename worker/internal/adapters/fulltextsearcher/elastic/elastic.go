package elastic

import (
	"context"

	"github.com/Unlites/ml-analysis-provider/worker/internal/domain"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const analysisIndex = "analysis"

// ElasticFullTextSearcher is a full text search adapter
type ElasticFullTextSearcher struct {
	client *elasticsearch.TypedClient
}

// NewElasticFullTextSearcher creates new ElasticFullTextSearcher
func NewElasticFullTextSearcher(client *elasticsearch.TypedClient) *ElasticFullTextSearcher {
	return &ElasticFullTextSearcher{
		client: client,
	}
}

// SearchAnalyzes searches analyzes by given filter and returnes slice of founded analyzes ids
func (s *ElasticFullTextSearcher) SearchAnalyzes(
	ctx context.Context,
	filter domain.AnalyzesFilter,
) ([]string, error) {
	query := &types.Query{
		Match: map[string]types.MatchQuery{
			"query":  {Query: filter.Query},
			"answer": {Query: filter.Answer},
		},
	}

	if filter.IsUserSatisfied != nil {
		query.Term = map[string]types.TermQuery{
			"is_user_satisfied": {Value: *filter.IsUserSatisfied},
		}
	}

	s.client.Search().Index(analysisIndex).Request(
		&search.Request{
			Query: query,
			From:  &filter.Offset,
			Size:  &filter.Limit,
		},
	)

	return nil, nil
}
