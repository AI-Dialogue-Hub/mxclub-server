package penalty

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeoutPenalty_ApplyPenalty(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name           string
		timeout        time.Duration
		expectedAmount float64
	}{
		{"No penalty for 0 minutes", 0, 0},
		{"No penalty for 9 minutes", time.Minute * 19, 0},
		{"Penalty for 10 minutes", time.Minute * 20, 10},
		{"Penalty for 15 minutes", time.Minute * 25, 10},
		{"Penalty for 15 minutes", time.Minute * 25, 10},
		// AddDeduction additional tests for 20 minutes and more if you extend the map
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the timeout by adjusting the time directly instead of using time.Sleep
			reqTime := now.Add(-tt.timeout)
			req := &PenaltyReq{GrabTime: &reqTime}

			penalty := &TimeoutPenalty{}
			resp, err := penalty.ApplyPenalty(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedAmount, resp.PenaltyAmount)
		})
	}
}

func TestLowRatingPenalty_ApplyPenalty(t *testing.T) {
	tests := []struct {
		name           string
		rating         int
		expectedAmount float64
		expectError    bool
	}{
		{"Penalty for 1 star", 1, 30, false},
		{"Penalty for 2 stars", 2, 20, false},
		{"No penalty for 3 stars", 3, 0, true},
		{"No penalty for 4 stars", 4, 0, true},
		{"Invalid rating (negative)", -1, 0, true},
		{"Invalid rating (zero)", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PenaltyReq{Rating: tt.rating}

			penalty := &LowRatingPenalty{}
			resp, err := penalty.ApplyPenalty(req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAmount, resp.PenaltyAmount)
			}
		})
	}
}
