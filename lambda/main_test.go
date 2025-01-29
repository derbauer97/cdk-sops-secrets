package main

import (
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-test/deep"
)

func TestMain(t *testing.M) {

	//log.SetOutput(ioutil.Discard)

	v := t.Run()
	// After all tests have run `go-snaps` can check for not used snapshots
	snaps.Clean(t)

	os.Exit(v)
}

func Test_GetS3FileContent(t *testing.T) {
	mocks := getMocks(t)
	data, err := mocks.getS3FileContent(SopsS3File{
		Bucket: "..",
		Key:    "../test-secrets/json/sopsfile.enc-age.json",
	})
	check(err)
	snaps.MatchSnapshot(t, string(data))
}

func Test_UpdateSecret(t *testing.T) {
	mocks := getMocks(t)
	fileName := "4547532a137611d83958d17095c6c2d38ae0036a760c3b79c9dd5957d1c20cf2.yaml"
	inputArn := "arn:${Partition}:secretsmanager:${Region}:${Account}:secret:${SecretId}"
	secretValue := []byte("some-secret-data")

	response, err := mocks.updateSecret(fileName, inputArn, secretValue)
	check(err)

	snaps.MatchSnapshot(t, response)
}

func Test_UpdateSSMParameter(t *testing.T) {
	mocks := getMocks(t)

	paramterName := "/foo/bar"
	parameterValue := []byte("some-secret-data")

	response, err := mocks.updateSSMParameter(paramterName, parameterValue, "key")
	check(err)

	snaps.MatchSnapshot(t, response)
}

func Test_DecryptSopsFileContent(t *testing.T) {

	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-1EFUWJ0G2XJTJFWTAM2DGMA4VCK3R05W58FSMHZP3MZQ0ZTAQEAFQC6T7T3")

	sopsEncrypted, err := os.ReadFile("../test-secrets/json/sopsfile.enc-age.json")
	check(err)
	sopsDecrypted, err := decryptSopsFileContent(sopsEncrypted, "json")
	check(err)
	sopsExpected, err := os.ReadFile("../test-secrets/json/sopsfile.json")
	check(err)

	sopsDecryptedJ := UnmarshalAny(sopsDecrypted)
	sopsExpectedJ := UnmarshalAny(sopsExpected)

	if diff := deep.Equal(sopsDecryptedJ, sopsExpectedJ); diff != nil {
		t.Error(diff)
	}
}
func Test_IsHumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "Human readable string",
			input:    []byte("Hello, World!"),
			expected: true,
		},
		{
			name:     "String with non-printable characters",
			input:    []byte("Hello\x00World"),
			expected: false,
		},
		{
			name:     "String with only spaces",
			input:    []byte("     "),
			expected: true,
		},
		{
			name:     "Empty string",
			input:    []byte(""),
			expected: true,
		},
		{
			name:     "String with null byte",
			input:    []byte{0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHumanReadable(tt.input)
			if result != tt.expected {
				t.Errorf("isHumanReadable(%q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
