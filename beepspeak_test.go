package beepspeak

import (
	"os"
	"testing"
)

func TestInitSay(t *testing.T) {
	cipherText := ""
	passPhrase := ""

	// Invalid cipherText / passPhrase
	err := InitSay(cipherText, passPhrase)
	if err == nil {
		t.Errorf("ERROR: Expected err != nil")
	}

	cipherText = "hlwjlfFMLUIWbbjyphbr4b1mbCsGYS3Ciy4gXvSnZg=="
	passPhrase = "foo"

	env := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if env != "" {
		t.Errorf("ERROR: env var is already set, %v", env)
	}

	err = InitSay(cipherText, passPhrase)
	if err != nil {
		t.Errorf("ERROR: Expected err == nil, got %v", err)
	}

	env = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if env == "" {
		t.Errorf("ERROR: env var is not set")
	}
}

func TestReadable(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"12345678", "12345678"},
		{"^^^^", ""},
		{"a_b/c[d]e\"f^g", "a b c d efg"},
	}

	for _, testCase := range testCases {
		answer := readable(testCase.input)
		if answer != testCase.expected {
			t.Errorf("ERROR: For %v expected %v, got %v", testCase.input, testCase.expected, answer)
		}
	}
}
