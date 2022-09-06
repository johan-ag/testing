package examples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	ErrLanguageNotFound = fmt.Errorf("error language not found")
)

func sayHello(language string) (string, error) {
	languages := map[string]string{
		"es": "hola",
		"en": "hello",
		"fr": "bonjour",
	}

	hello, exist := languages[language]

	if !exist {
		return "", ErrLanguageNotFound
	}

	return hello, nil
}

func TestSayHello(t *testing.T) {
	tests := []struct {
		name     string
		language string
		want     string
	}{
		{
			name:     "success say hello in english",
			language: "en",
			want:     "hello",
		},
		{
			name:     "success say hello in french",
			language: "fr",
			want:     "bonjour",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got, err := sayHello(tt.language)
			if err != nil {
				t.Log()
				t.Fail()
			}

			// then
			require.Equal(t, tt.want, got)
		})
	}
}
