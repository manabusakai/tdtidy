package tdtidy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_selectTaskDefinitions(t *testing.T) {
	now := time.Now()
	app := &App{}

	type args struct {
		threshold time.Time
		tds       []taskdef
	}
	tests := []struct {
		name    string
		args    args
		want    families
		wantErr bool
	}{
		{
			name: "If you do not specify a retention period",
			args: args{
				threshold: now,
				tds: []taskdef{
					{
						family:         "tdtidy",
						revision:       2,
						registeredAt:   getPastDate(now, 5),
						deregisteredAt: nil,
					},
					{
						family:         "tdtidy",
						revision:       1,
						registeredAt:   getPastDate(now, 10),
						deregisteredAt: nil,
					},
				},
			},
			want: families{
				"tdtidy": {
					"tdtidy:2",
					"tdtidy:1",
				},
			},
			wantErr: false,
		},
		{
			name: "If you specify a retention period.",
			args: args{
				threshold: now.AddDate(0, 0, -7).UTC(),
				tds: []taskdef{
					{
						family:         "tdtidy",
						revision:       2,
						registeredAt:   getPastDate(now, 5),
						deregisteredAt: nil,
					},
					{
						family:         "tdtidy",
						revision:       1,
						registeredAt:   getPastDate(now, 10),
						deregisteredAt: nil,
					},
				},
			},
			want: families{
				"tdtidy": {
					"tdtidy:1",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := app.selectTaskDefinitions(tt.args.threshold, tt.args.tds)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.selectTaskDefinitions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("App.selectTaskDefinitions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getPastDate(now time.Time, daysAgo int) *time.Time {
	t := now.AddDate(0, 0, -daysAgo).UTC()
	return &t
}
