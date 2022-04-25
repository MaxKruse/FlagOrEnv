package flagorenv

import (
	"os"
	"testing"
)

type EmptyStruct struct{}

type TestExpectedTypeError struct {
	TypeErrorField int
}

type TestExpectedTypeSuccess struct {
	TypeSuccessField int64
}

func TestConfigDefaults(t *testing.T) {
	c := &Config{}

	LoadFlagsOrEnv[EmptyStruct](c)

	if c.PreferFlag != false {
		t.Errorf("PreferFlag should be false by default, was %v", c.PreferFlag)
	}

	if c.Prefix != "FLAGENV" {
		t.Errorf("Prefix should be 'FLAGENV' by default, was %v", c.Prefix)
	}
}

func TestExpectedTypeErrorFromEnv(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"

	os.Setenv("TEST_TYPE_ERROR_FIELD", "1")
	defer os.Unsetenv("TEST_TYPE_ERROR_FIELD")

	res, err := LoadFlagsOrEnv[TestExpectedTypeError](c)

	if err == nil {
		t.Errorf("Expected error, got success: %v", res)
	}
}

func TestExpectedTypeSuccessFromEnv(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"

	os.Setenv("TEST_TYPE_SUCCESS_FIELD", "1")
	defer os.Unsetenv("TEST_TYPE_SUCCESS_FIELD")

	res, err := LoadFlagsOrEnv[TestExpectedTypeSuccess](c)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if res.TypeSuccessField != 1 {
		t.Errorf("Expected TypeSuccessField to be 1, got %v", res.TypeSuccessField)
	}
}

// Always run tests with a "-test-type-error-field 1" flag
func TestExpectedTypeErrorFromFlag(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"
	c.PreferFlag = true

	res, err := LoadFlagsOrEnv[TestExpectedTypeError](c)

	if err == nil {
		t.Errorf("Expected error, got success: %v", res)
	}
}

// Always run tests with a "-test-type-success-field 1" flag
func TestExpectedTypeSuccessFromFlag(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"
	c.PreferFlag = true

	res, err := LoadFlagsOrEnv[TestExpectedTypeSuccess](c)

	if err == nil {
		t.Errorf("Expected error, got success: %v", res)
	}
}
