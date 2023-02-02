package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFastAesKey(t *testing.T) {
	testCases := []struct {
		name string
		data string

		wantDecryptKey string
	}{
		{
			name:           "813b28aeede3e3fc1daa2fce885a4b8a",
			data:           "813b28aeede3e3fc1daa2fce885a4b8a:3sAIbKUZjBF28VZRcFNIdNceE91GEYa4MDcKsy4Jfog=",
			wantDecryptKey: "fa417156b9a34368",
		},
		{
			name:           "afe96771a8df7154f6aa2a9587484c63",
			data:           "afe96771a8df7154f6aa2a9587484c63:a+89LQn8j98H7n/+uoNRTwRyXZbkkflkYOrztZZmkq8=",
			wantDecryptKey: "fa417156b9a34368",
		},
		{
			name:           "9ea0ba752acc77b742772464da5f4d14",
			data:           "9ea0ba752acc77b742772464da5f4d14:baVwjtEpEK/qSUKXAO+gGzgXDTEHCRwjlHFPy9ZIWmY=",
			wantDecryptKey: "fa417156b9a34368",
		},
		{
			name:           "96ef2d4c35c0eb6b20a606a7f2294769",
			data:           "96ef2d4c35c0eb6b20a606a7f2294769:QLpttSKbOoPEDyWSor3NdmkVUxXCzDLjZYaMsEG9PNs=",
			wantDecryptKey: "fa417156b9a34368",
		},
		{
			name:           "96ef2d4c35c0eb6b20a606a7f2294769",
			data:           "96ef2d4c35c0eb6b20a606a7f2294769:QLpttSKbOoPEDyWSor3NdmkVUxXCzDLjZYaMsEG9PNs=",
			wantDecryptKey: "fa417156b9a34368",
		},
		{
			name:           "a9e892ead614da64975031d98abd48e1",
			data:           "a9e892ead614da64975031d98abd48e1:prg62z5UadxONGxKHZNCGFJ4Rse2POTix1n4ZqFHLbI=",
			wantDecryptKey: "5fa20ccda77445bd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decryptKey, err := FastAesKey(tc.data)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantDecryptKey, decryptKey)
		})
	}
}

func TestPlayAuthDecrypt(t *testing.T) {
	testCases := []struct {
		name       string
		encryptKey string

		wantDecryptKey string
	}{
		{
			name:           "p7wpy2KMLMRhiBjYYZIt2WKXKoyM",
			encryptKey:     "p7wpy2KMLMRhiBjYYZIt2WKXKoyM",
			wantDecryptKey: "fa417156b9a34368",
		},
		{
			name:           "kLwewGaIIfJaiCLGbY0hwFq6EpmZ",
			encryptKey:     "kLwewGaIIfJaiCLGbY0hwFq6EpmZ",
			wantDecryptKey: "5fa20ccda77445bd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fastAesKey := PlayAuthDecrypt(tc.encryptKey)
			assert.Equal(t, tc.wantDecryptKey, fastAesKey)
		})
	}
}
