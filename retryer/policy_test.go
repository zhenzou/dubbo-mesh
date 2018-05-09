// 重试策略

package retryer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimesPolicy_Continue(t *testing.T) {
	p := NewTimesPolicy(3)
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "time_1",
			want: true,
		},
		{
			name: "time_2",
			want: true,
		},
		{
			name: "time_3",
			want: true,
		},
		{
			name: "time_4",
			want: false,
		},
		{
			name: "time_5",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.Continue(); got != tt.want {
				t.Errorf("TimesPolicy.Continue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimes_Interval(t *testing.T) {

	policy := &TimesPolicy{maxAttempts: 5, period: 100, maxPeriod: int64(1 * time.Millisecond)}

	assert.Equal(t, 1, policy.attempt)
	assert.Equal(t, 0, policy.sleptForMillis)
	policy.Continue()
	assert.Equal(t, 2, policy.attempt)
	assert.Equal(t, 150, policy.sleptForMillis)
	policy.Continue()
	assert.Equal(t, 3, policy.attempt)
	assert.Equal(t, 375, policy.sleptForMillis)
	policy.Continue()
	assert.Equal(t, 4, policy.attempt)
	assert.Equal(t, 712, policy.sleptForMillis)
	policy.Continue()
	assert.Equal(t, 5, policy.attempt)
	assert.Equal(t, 1218, policy.sleptForMillis)
	policy.Continue()
}

func TestDeadlinePolicy_Continue(t *testing.T) {
	type fields struct {
		deadline time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "1_second",
			fields: fields{1 * time.Second},
			want:   false,
		},
		{
			name:   "2_second",
			fields: fields{2 * time.Second},
			want:   false,
		},
		{
			name:   "3_second",
			fields: fields{3 * time.Second},
			want:   true,
		},
		{
			name:   "4_second",
			fields: fields{4 * time.Second},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeadlinePolicy(tt.fields.deadline)
			time.Sleep(2 * time.Second)
			if got := d.Continue(); got != tt.want {
				t.Errorf("DeadlinePolicy.Continue() = %v, want %v", got, tt.want)
			}
		})
	}
}
