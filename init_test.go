package tdtidy

import (
	"flag"
	"os"
	"testing"

	. "github.com/samber/lo"
)

func Test_initOption(t *testing.T) {
	type arg struct {
		name  string
		value string
	}
	tests := []struct {
		name string
		cmd  string
		args []arg
		want *option
	}{
		{
			name: "default option",
			args: []arg{},
			want: &option{
				dryRun:          ToPtr(false),
				retentionPeriod: ToPtr(0),
				familyPrefix:    nil,
			},
		},
		{
			name: "all option",
			cmd:  "delete",
			args: []arg{
				{
					name:  "dry-run",
					value: "true",
				},
				{
					name:  "retention-period",
					value: "7",
				},
				{
					name:  "family-prefix",
					value: "dummy",
				},
			},
			want: &option{
				dryRun:          ToPtr(true),
				retentionPeriod: ToPtr(7),
				familyPrefix:    ToPtr("dummy"),
			},
		},
		{
			name: "dry-run option only",
			cmd:  "delete",
			args: []arg{
				{
					name:  "dry-run",
					value: "true",
				},
			},
			want: &option{
				dryRun:          ToPtr(true),
				retentionPeriod: ToPtr(0),
				familyPrefix:    nil,
			},
		},
		{
			name: "retention-period option only",
			cmd:  "delete",
			args: []arg{
				{
					name:  "retention-period",
					value: "7",
				},
			},
			want: &option{
				dryRun:          ToPtr(false),
				retentionPeriod: ToPtr(7),
				familyPrefix:    nil,
			},
		},
		{
			name: "family-prefix option only",
			cmd:  "delete",
			args: []arg{
				{
					name:  "family-prefix",
					value: "dummy",
				},
			},
			want: &option{
				dryRun:          ToPtr(false),
				retentionPeriod: ToPtr(0),
				familyPrefix:    ToPtr("dummy"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add a subcommand into os.Args
			args := os.Args
			os.Args = append([]string{tt.cmd}, args...)

			defer func() {
				os.Args = args
			}()

			fs := flag.NewFlagSet(tt.name, flag.ContinueOnError)
			dryRun = fs.Bool("dry-run", false, "")
			retentionPeriod = fs.Int("retention-period", 0, "")
			familyPrefix = fs.String("family-prefix", "", "")

			for _, arg := range tt.args {
				fs.Set(arg.name, arg.value)
			}

			got, err := initOption()
			if err != nil {
				t.Fatalf("initOption() error: %v,", err)
			}

			if *got.dryRun != *tt.want.dryRun {
				t.Errorf("dryRun = %v, want %v", *got.dryRun, *tt.want.dryRun)
			}
			if *got.retentionPeriod != *tt.want.retentionPeriod {
				t.Errorf("retentionPeriod = %v, want %v", *got.retentionPeriod, *tt.want.retentionPeriod)
			}
			if (got.familyPrefix != nil && tt.want.familyPrefix == nil) || (got.familyPrefix == nil && tt.want.familyPrefix != nil) {
				t.Errorf("familyPrefix = %v, want %v", got.familyPrefix, tt.want.familyPrefix)
			}
			if (got.familyPrefix != nil && tt.want.familyPrefix != nil) && (*got.familyPrefix != *tt.want.familyPrefix) {
				t.Errorf("familyPrefix = %v, want %v", *got.familyPrefix, *tt.want.familyPrefix)
			}
		})
	}
}
