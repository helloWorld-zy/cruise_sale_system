package image

import (
	"context"
	"fmt"
)

type PosterService struct{}

func NewPosterService() *PosterService {
	return &PosterService{}
}

func (s *PosterService) GenerateItineraryPoster(ctx context.Context, orderID string) (string, error) {
	// Mock implementation: Return a placeholder URL
	// Real implementation would use 'github.com/fogleman/gg' or similar
	return fmt.Sprintf("https://mock-storage.com/posters/%s.png", orderID), nil
}
