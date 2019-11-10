package user_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/jacoelho/codewars/internal/user"
)

func TestUserDecode(t *testing.T) {
	testCases := []struct {
		desc    string
		input   string
		expect  user.Event
		wantErr bool
	}{
		{
			desc:  "happy path",
			input: `{"action":"honor_changed","user":{"id":"570a39399cd34a8a740012a9","honor":308,"honor_delta":2}}`,
			expect: user.Event{
				HonorUpgraded: &user.HonorUpgraded{
					ID:         "570a39399cd34a8a740012a9",
					Honor:      308,
					HonorDelta: 2,
				},
			},
		},
		{
			desc:  "empty event",
			input: `{"action":"honor_changed","user":{}}`,
			expect: user.Event{
				HonorUpgraded: &user.HonorUpgraded{},
			},
		},
		{
			desc:  "invalid honor ",
			input: `{"action":"honor_changed","user":{"id":"570a39399cd34a8a740012a9","honor":"foobar","honor_delta":2}}`,
			expect: user.Event{
				HonorUpgraded: &user.HonorUpgraded{
					ID:         "570a39399cd34a8a740012a9",
					Honor:      308,
					HonorDelta: 2,
				},
			},
			wantErr: true,
		},
		{
			desc:  "unknown action",
			input: `{"action":"unknown","user":{"id":"570a39399cd34a8a740012a9","honor":308,"honor_delta":2}}`,
			expect: user.Event{
				HonorUpgraded: &user.HonorUpgraded{},
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			var data user.Event
			if err := json.NewDecoder(strings.NewReader(tt.input)).Decode(&data); (err != nil) != tt.wantErr {
				t.Fatalf("Event.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.expect, data) {
				t.Fatalf("Expected = %v, got = %v", tt.expect, data)
			}
		})
	}
}
