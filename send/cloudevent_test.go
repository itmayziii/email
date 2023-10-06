package send_test

import (
	"github.com/itmayziii/email/send"
	"testing"
)

func TestMessageTo_UnmarshalJSON_UnmarshalsStrings(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected send.MessageTo
	}{
		{"empty string", []byte(""), send.MessageTo{}},
		{"empty value", []byte{}, send.MessageTo{}},
		{"non empty string", []byte(`"non empty string"`), send.MessageTo{"non empty string"}},
	}

	for _, tt := range tests {
		ttCopy := tt
		t.Run(ttCopy.name, func(t *testing.T) {
			t.Parallel()
			actual := send.MessageTo{}
			err := actual.UnmarshalJSON(ttCopy.data)

			if err != nil {
				t.Errorf("case: \"%s\", unexpected error: %v", ttCopy.name, err)
			}
			if len(ttCopy.expected) != len(actual) {
				t.Errorf("expected %s to match %s", actual, ttCopy.expected)
			}
			for i, a := range actual {
				expected := ttCopy.expected[i]
				if expected != a {
					t.Errorf("expected %s to match %s", a, ttCopy.expected[i])
				}
			}
		})
	}
}

func TestMessageTo_UnmarshalJSON_UnmarshalsStringArrays(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected send.MessageTo
	}{
		{"empty array", []byte("[]"), send.MessageTo{}},
		{"non empty array", []byte(`["hello", "world"]`), send.MessageTo{"hello", "world"}},
	}

	for _, tt := range tests {
		ttCopy := tt
		t.Run(ttCopy.name, func(t *testing.T) {
			t.Parallel()
			actual := send.MessageTo{}
			err := actual.UnmarshalJSON(ttCopy.data)

			if err != nil {
				t.Error(err)
			}
			if len(ttCopy.expected) != len(actual) {
				t.Errorf("expected %v to match %v", actual, ttCopy.expected)
			}
			for i, a := range actual {
				expected := ttCopy.expected[i]
				if expected != a {
					t.Errorf("expected %s to match %s", a, ttCopy.expected[i])
				}
			}
		})
	}
}

func TestMessageTo_UnmarshalJSON_ErrorsIfJsonIsNotStringOrStringArray(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"syntax error", []byte("{,a,jii3{")},
		{"JSON object does not fit into string or []string", []byte(`{"hello":"world"}`)},
	}

	for _, tt := range tests {
		ttCopy := tt
		t.Run(ttCopy.name, func(t *testing.T) {
			t.Parallel()
			actual := send.MessageTo{}
			err := actual.UnmarshalJSON(ttCopy.data)

			if err != nil {
				return
			}
			t.Errorf("case: \"%s\", error was expected but actual value is %s", ttCopy.name, actual)
		})
	}
}
