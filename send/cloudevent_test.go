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
		{"empty string", []byte(""), send.MessageTo{""}},
		{"non empty string", []byte("non empty string"), send.MessageTo{"non empty string"}},
		{"JSON object", []byte("{\"hello\":\"world\"}"), send.MessageTo{"{\"hello\":\"world\"}"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := send.MessageTo{}
			err := actual.UnmarshalJSON(tt.data)

			if err != nil {
				t.Error(err)
			}
			if len(tt.expected) != len(actual) {
				t.Errorf("expected %v to match %v", actual, tt.expected)
			}
			for i, a := range actual {
				expected := tt.expected[i]
				if expected != a {
					t.Errorf("expected %s to match %s", a, tt.expected[i])
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
		{"non empty array", []byte("[\"hello\", \"world\"]"), send.MessageTo{"hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := send.MessageTo{}
			err := actual.UnmarshalJSON(tt.data)

			if err != nil {
				t.Error(err)
			}
			if len(tt.expected) != len(actual) {
				t.Errorf("expected %v to match %v", actual, tt.expected)
			}
			for i, a := range actual {
				expected := tt.expected[i]
				if expected != a {
					t.Errorf("expected %s to match %s", a, tt.expected[i])
				}
			}
		})
	}
}

func TestMessageTo_UnmarshalJSON_ErrorsIfInvalidJson(t *testing.T) {
	t.Parallel()
	actual := send.MessageTo{}
	err := actual.UnmarshalJSON([]byte(",a,jii3{"))

	if err == nil {
		return
	}

	t.Errorf("error was expected but instead actual value is %v", actual)
}
