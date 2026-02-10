package core

import (
	"context"
	"log"
)

type MarketingService struct {
	// dependencies
}

func NewMarketingService() *MarketingService {
	return &MarketingService{}
}

func (s *MarketingService) TriggerCartAbandonment(ctx context.Context, userID string) error {
	log.Printf("Sending cart abandonment email to user %s", userID)
	return nil
}

func (s *MarketingService) TriggerPostTripReview(ctx context.Context, userID string, orderID string) error {
	log.Printf("Inviting user %s to review order %s", userID, orderID)
	return nil
}
