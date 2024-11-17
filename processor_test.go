package tdtidy

import (
	"context"
	"testing"
	"time"

	"github.com/manabusakai/tdtidy/internal/ecs"
	. "github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func Test_filterTaskDefinitions(t *testing.T) {
	// Set up
	client := ecs.NewClient(context.TODO())
	time.Local = time.UTC
	now := time.Now()

	type fields struct {
		client ecs.Client
		opt    *option
	}
	type args struct {
		tds []ecs.TaskDefinition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []ecs.TaskDefinition
	}{
		{
			name: "If you do not specify a retention period",
			fields: fields{
				client: client,
				opt: &option{
					retentionPeriod: ToPtr(0),
				},
			},
			args: args{
				tds: []ecs.TaskDefinition{
					{
						Revision:       2,
						RegisteredAt:   ToPtr(now.AddDate(0, 0, -5)),
						DeregisteredAt: nil,
					},
					{
						Revision:       1,
						RegisteredAt:   ToPtr(now.AddDate(0, 0, -10)),
						DeregisteredAt: nil,
					},
				},
			},
			want: []ecs.TaskDefinition{
				{
					Revision:       2,
					RegisteredAt:   ToPtr(now.AddDate(0, 0, -5)),
					DeregisteredAt: nil,
				},
				{
					Revision:       1,
					RegisteredAt:   ToPtr(now.AddDate(0, 0, -10)),
					DeregisteredAt: nil,
				},
			},
		},
		{
			name: "If you specify a retention period",
			fields: fields{
				client: client,
				opt: &option{
					retentionPeriod: ToPtr(7),
				},
			},
			args: args{
				tds: []ecs.TaskDefinition{
					{
						Revision:       2,
						RegisteredAt:   ToPtr(now.AddDate(0, 0, -5)),
						DeregisteredAt: nil,
					},
					{
						Revision:       1,
						RegisteredAt:   ToPtr(now.AddDate(0, 0, -10)),
						DeregisteredAt: nil,
					},
				},
			},
			want: []ecs.TaskDefinition{
				{
					Revision:       1,
					RegisteredAt:   ToPtr(now.AddDate(0, 0, -10)),
					DeregisteredAt: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Processor{
				client: tt.fields.client,
				opt:    tt.fields.opt,
			}
			if got := p.filterTaskDefinitions(tt.args.tds); !assert.Equal(t, got, tt.want) {
				t.Errorf("App.filterTaskDefinitions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidTaskDefinition(t *testing.T) {
	// Set up
	client := ecs.NewClient(context.TODO())
	time.Local = time.UTC
	now := time.Now()

	type fields struct {
		client ecs.Client
		opt    *option
	}
	type args struct {
		td ecs.TaskDefinition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "RegisteredAt is missing",
			fields: fields{
				client: client,
				opt: &option{
					retentionPeriod: ToPtr(7),
				},
			},
			args: args{
				td: ecs.TaskDefinition{
					RegisteredAt:   nil,
					DeregisteredAt: nil,
				},
			},
			want: false,
		},
		{
			name: "RegisteredAt is after the threshold",
			fields: fields{
				client: client,
				opt: &option{
					retentionPeriod: ToPtr(7),
				},
			},
			args: args{
				td: ecs.TaskDefinition{
					RegisteredAt:   ToPtr(now.AddDate(0, 0, -1)),
					DeregisteredAt: nil,
				},
			},
			want: false,
		},
		{
			name: "RegisteredAt is before the threshold",
			fields: fields{
				client: client,
				opt: &option{
					retentionPeriod: ToPtr(7),
				},
			},
			args: args{
				td: ecs.TaskDefinition{
					RegisteredAt:   ToPtr(now.AddDate(0, 0, -10)),
					DeregisteredAt: nil,
				},
			},
			want: true,
		},
		{
			name: "DeregisteredAt is after the threshold",
			fields: fields{
				client: client,
				opt: &option{
					retentionPeriod: ToPtr(7),
				},
			},
			args: args{
				td: ecs.TaskDefinition{
					RegisteredAt:   ToPtr(now.AddDate(0, 0, -1)),
					DeregisteredAt: ToPtr(now.AddDate(0, 0, -1)),
				},
			},
			want: false,
		},
		{
			name: "DeregisteredAt is before the threshold",
			fields: fields{
				client: client,
				opt: &option{
					retentionPeriod: ToPtr(7),
				},
			},
			args: args{
				td: ecs.TaskDefinition{
					RegisteredAt:   ToPtr(now.AddDate(0, 0, -10)),
					DeregisteredAt: ToPtr(now.AddDate(0, 0, -10)),
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Processor{
				client: tt.fields.client,
				opt:    tt.fields.opt,
			}
			if got := p.isValidTaskDefinition(&tt.args.td, tt.fields.opt.threshold()); got != tt.want {
				t.Errorf("App.isValidTaskDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}
