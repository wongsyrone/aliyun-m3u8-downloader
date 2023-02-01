package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFastAesKey(t *testing.T) {
	testCases := []struct {
		name string
		data string

		wantFastAesKey string
	}{
		{
			name:           "813b28aeede3e3fc1daa2fce885a4b8a",
			data:           "813b28aeede3e3fc1daa2fce885a4b8a:3sAIbKUZjBF28VZRcFNIdNceE91GEYa4MDcKsy4Jfog=",
			wantFastAesKey: "fa417156b9a34368",
		},
		{
			name:           "afe96771a8df7154f6aa2a9587484c63",
			data:           "afe96771a8df7154f6aa2a9587484c63:a+89LQn8j98H7n/+uoNRTwRyXZbkkflkYOrztZZmkq8=",
			wantFastAesKey: "fa417156b9a34368",
		},
		{
			name:           "9ea0ba752acc77b742772464da5f4d14",
			data:           "9ea0ba752acc77b742772464da5f4d14:baVwjtEpEK/qSUKXAO+gGzgXDTEHCRwjlHFPy9ZIWmY=",
			wantFastAesKey: "fa417156b9a34368",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fastAesKey, err := FastAesKey(tc.data)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantFastAesKey, fastAesKey)
		})
	}
}
