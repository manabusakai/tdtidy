package tdtidy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_selectTaskDefinitions(t *testing.T) {
	now := time.Now()
	app := &App{
		ctx: context.TODO(),
		opt: &option{},
	}

	type args struct {
		threshold time.Time
		tds       []taskdef
	}
	tests := []struct {
		name    string
		args    args
		want    []taskdef
		wantErr bool
	}{
		{
			name: "If you do not specify a retention period",
			args: args{
				threshold: now,
				tds: []taskdef{
					{
						arn:            "arn:aws:ecs:ap-northeast-1:123456789012:task-definition/tdtidy:2",
						family:         "tdtidy",
						revision:       2,
						registeredAt:   getPastDate(now, 5),
						deregisteredAt: nil,
					},
					{
						arn:            "arn:aws:ecs:ap-northeast-1:123456789012:task-definition/tdtidy:1",
						family:         "tdtidy",
						revision:       1,
						registeredAt:   getPastDate(now, 10),
						deregisteredAt: nil,
					},
				},
			},
			want: []taskdef{
				{
					arn:            "arn:aws:ecs:ap-northeast-1:123456789012:task-definition/tdtidy:2",
					family:         "tdtidy",
					revision:       2,
					registeredAt:   getPastDate(now, 5),
					deregisteredAt: nil,
				},
				{
					arn:            "arn:aws:ecs:ap-northeast-1:123456789012:task-definition/tdtidy:1",
					family:         "tdtidy",
					revision:       1,
					registeredAt:   getPastDate(now, 10),
					deregisteredAt: nil,
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
						arn:            "arn:aws:ecs:ap-northeast-1:123456789012:task-definition/tdtidy:2",
						family:         "tdtidy",
						revision:       2,
						registeredAt:   getPastDate(now, 5),
						deregisteredAt: nil,
					},
					{
						arn:            "arn:aws:ecs:ap-northeast-1:123456789012:task-definition/tdtidy:1",
						family:         "tdtidy",
						revision:       1,
						registeredAt:   getPastDate(now, 10),
						deregisteredAt: nil,
					},
				},
			},
			want: []taskdef{
				{
					arn:            "arn:aws:ecs:ap-northeast-1:123456789012:task-definition/tdtidy:1",
					family:         "tdtidy",
					revision:       1,
					registeredAt:   getPastDate(now, 10),
					deregisteredAt: nil,
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
