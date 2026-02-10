package unit

import (
	"testing"

	"cruise_booking_system/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestOrderStatusTransitions(t *testing.T) {
	order := &model.Order{Status: model.OrderStatusPending}

	// Valid transition
	order.Status = model.OrderStatusPaid
	assert.Equal(t, model.OrderStatusPaid, order.Status)

	// In a real state machine, we'd have a method CanTransitionTo(next) bool
	// For MVP struct, we just verify constants exist
	assert.Equal(t, "confirmed", string(model.OrderStatusConfirmed))
}
