package performance

import (
	"reflect"
	"testing"
)

func TestParsingCast(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"장사익, 하고운, 윤덕현, 권순찬, 윤부식, 유용진 등", []string{"장사익", "하고운", "윤덕현", "권순찬", "윤부식", "유용진"}},
		{"황인경 등", []string{"황인경"}},
		{"박용규, 최현준, 김경록", []string{"박용규", "최현준", "김경록"}},
		{"", []string{}},
	}

	for _, tt := range tests {
		got, _ := ParsingCast(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParsingCast(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestParsingKeyword(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"K-pop|Ballad|Indie", []string{"K-pop", "Ballad", "Indie"}},
		{"HipHop", []string{"HipHop"}},
		{"", []string{}},
	}

	for _, tt := range tests {
		got, _ := ParsingKeyword(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParsingKeyword(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}
