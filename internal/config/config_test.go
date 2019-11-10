package config_test

import (
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/jacoelho/codewars/internal/config"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		expect config.Config
	}{
		{
			desc:  "example",
			input: "-port 9090 -secret abc -slack-webhook hook",
			expect: config.Config{
				Port:      "9090",
				SlackHook: "hook",
				Secret:    "abc",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			cfg := config.New()
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			cfg.Flags(fs)

			if err := fs.Parse(strings.Split(tt.input, " ")); err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			if !reflect.DeepEqual(&tt.expect, cfg) {
				t.Fatalf("expected %v, got %v", &tt.expect, cfg)
			}
		})
	}
}
