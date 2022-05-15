package flagorenv

import (
	"os"
	"testing"
)

type emptyStruct struct{}

type testExpectedTypeError struct {
	field emptyStruct
}

func TestConfigDefaults(t *testing.T) {
	c := &Config{}

	LoadFlagsOrEnv[emptyStruct](c)

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

	res, err := LoadFlagsOrEnv[testExpectedTypeError](c)

	if err == nil {
		t.Errorf("Expected error, got success: %v", res)
	}
}

type testExpectedTypeSuccess struct {
	TypeSuccessField int64
}

func TestExpectedTypeSuccessFromEnv(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"

	os.Setenv("TEST_TYPE_SUCCESS_FIELD", "1")
	defer os.Unsetenv("TEST_TYPE_SUCCESS_FIELD")

	res, err := LoadFlagsOrEnv[testExpectedTypeSuccess](c)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if res.TypeSuccessField != 1 {
		t.Errorf("Expected TypeSuccessField to be 1, got %v", res.TypeSuccessField)
	}
}

type testStringSlice struct {
	StringSlice []string
}

func TestStringSlice(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"

	os.Setenv("TEST_STRING_SLICE", "a,b,c")
	defer os.Unsetenv("TEST_STRING_SLICE")

	res, err := LoadFlagsOrEnv[testStringSlice](c)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if res.StringSlice[0] != "a" {
		t.Errorf("Expected StringSlice[0] to be 'a', got %v", res.StringSlice[0])
	}

	if res.StringSlice[1] != "b" {
		t.Errorf("Expected StringSlice[1] to be 'b', got %v", res.StringSlice[1])
	}

	if res.StringSlice[2] != "c" {
		t.Errorf("Expected StringSlice[2] to be 'c', got %v", res.StringSlice[2])
	}
}

type testInt64_Slice struct {
	Int64_Slice []int64
}

func TestInt64_Slice(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"

	os.Setenv("TEST_INT64_SLICE", "1,2,3")
	defer os.Unsetenv("TEST_INT64_SLICE")

	res, err := LoadFlagsOrEnv[testInt64_Slice](c)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if res.Int64_Slice[0] != 1 {
		t.Errorf("Expected Int64_Slice[0] to be 1, got %v", res.Int64_Slice[0])
	}

	if res.Int64_Slice[1] != 2 {
		t.Errorf("Expected Int64_Slice[1] to be 2, got %v", res.Int64_Slice[1])
	}

	if res.Int64_Slice[2] != 3 {
		t.Errorf("Expected Int64_Slice[2] to be 3, got %v", res.Int64_Slice[2])
	}
}

type testInt32_Slice struct {
	Int32_Slice []int32
}

func TestInt32_Slice(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"

	os.Setenv("TEST_INT32_SLICE", "1,2,3")
	defer os.Unsetenv("TEST_INT32_SLICE")

	res, err := LoadFlagsOrEnv[testInt32_Slice](c)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if res.Int32_Slice[0] != 1 {
		t.Errorf("Expected Int32_Slice[0] to be 1, got %v", res.Int32_Slice[0])
	}

	if res.Int32_Slice[1] != 2 {
		t.Errorf("Expected Int32_Slice[1] to be 2, got %v", res.Int32_Slice[1])
	}

	if res.Int32_Slice[2] != 3 {
		t.Errorf("Expected Int32_Slice[2] to be 3, got %v", res.Int32_Slice[2])
	}
}

type testBoolSlice struct {
	BoolSlice []bool
}

func TestBoolSlice(t *testing.T) {
	c := &Config{}
	c.Prefix = "TEST"

	os.Setenv("TEST_BOOL_SLICE", "true,false")
	defer os.Unsetenv("TEST_BOOL_SLICE")

	res, err := LoadFlagsOrEnv[testBoolSlice](c)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if res.BoolSlice[0] != true {
		t.Errorf("Expected BoolSlice[0] to be true, got %v", res.BoolSlice[0])
	}

	if res.BoolSlice[1] != false {
		t.Errorf("Expected BoolSlice[1] to be false, got %v", res.BoolSlice[1])
	}
}
