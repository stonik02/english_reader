package morphology

import (
	"reflect"
	"testing"
)

func TestServiceCandidates(t *testing.T) {
	service := New()
	for _, test := range []struct {
		word string
		want []string
	}{
		{"took", []string{"took", "take"}},
		{"crept", []string{"crept", "creep"}},
		{"is", []string{"is", "be"}},
		{"got", []string{"got", "get"}},
		{"hummed", []string{"hummed", "hum", "humm", "humme"}},
		{"putting", []string{"putting", "put", "putt", "putte"}},
		{"studied", []string{"studied", "studi", "study", "studie"}},
		{"book", []string{"book"}},
	} {
		if got := service.Candidates(test.word); !reflect.DeepEqual(got, test.want) {
			t.Fatalf("Candidates(%q)=%#v, want %#v", test.word, got, test.want)
		}
	}
}
