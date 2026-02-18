package service

import "context"

type BatchItem struct {
	CorrelationID string
	OriginalURL   string
}

type BatchResult struct {
	CorrelationID string
	ShortURL      string
}

type ShortenBatchService struct {
	shorten *ShortenService
}

func NewShortenBatchService(shorten *ShortenService) *ShortenBatchService {
	return &ShortenBatchService{shorten: shorten}
}

func (s *ShortenBatchService) Shorten(ctx context.Context, userID string, items []BatchItem) ([]BatchResult, error) {
	results := make([]BatchResult, 0, len(items))

	for _, item := range items {
		result, err := s.shorten.Shorten(ctx, userID, item.OriginalURL)
		if err != nil {
			s.shorten.logger.Errorf("Failed to shorten batch item %s: %v", item.CorrelationID, err)
			return nil, err
		}

		results = append(results, BatchResult{CorrelationID: item.CorrelationID, ShortURL: result.ShortURL})
	}

	return results, nil
}
