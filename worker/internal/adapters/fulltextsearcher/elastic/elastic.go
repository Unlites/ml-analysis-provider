package elastic

import (
	"context"
	"fmt"

	"github.com/Unlites/ml-analysis-provider/worker/internal/domain"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const analysisIndex = "ml_analysis"

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
		Bool: &types.BoolQuery{},
	}

	if filter.Query != "" {
		query.Bool.Must = append(query.Bool.Must, types.Query{
			Match: map[string]types.MatchQuery{
				"query": {
					Query: filter.Query,
				},
			},
		})
	}

	if filter.Answer != "" {
		query.Bool.Must = append(query.Bool.Must, types.Query{
			Match: map[string]types.MatchQuery{
				"answer": {
					Query: filter.Answer,
				},
			},
		})
	}

	if filter.IsUserSatisfied != nil {
		query.Bool.Must = append(query.Bool.Must, types.Query{
			Term: map[string]types.TermQuery{
				"is_user_satisfied": {Value: *filter.IsUserSatisfied},
			},
		})
	}

	res, err := s.client.Search().Index(analysisIndex).Request(
		&search.Request{
			Query: query,
			From:  &filter.Offset,
			Size:  &filter.Limit,
		},
	).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("fialed to do analyzes full text search: %w", err)
	}

	ids := make([]string, len(res.Hits.Hits))
	for i, hit := range res.Hits.Hits {
		ids[i] = *hit.Id_
	}

	return ids, nil
}
