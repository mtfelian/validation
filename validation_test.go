package validation

import (
	"testing"
)

type structToValidate struct {
	checkRequired *string
	checkMin      int
	checkMax      int
	checkRange    int
}

// Validate the s
func (s structToValidate) Validate() (bool, *Validation) {
	v := &Validation{}
	v.Check(s.checkRequired, Required{}).Message("checkRequired should not be empty")
	v.Check(s.checkMin, Min{2}).Message("checkMin field value too small")
	v.Check(s.checkMax, Max{5}).Message("checkMax field value too high")
	v.Check(s.checkRange, Range{Min{2}, Max{5}}).Message("checkRange is out of range")
	return !v.HasErrors(), v
}

func getValidStructToValidate() structToValidate {
	return structToValidate{
		checkRequired: pString("q"),
		checkMin:      2,
		checkMax:      5,
		checkRange:    3,
	}
}

func pString(s string) *string { return &s }

// TestValidateIfOK checks valid structure validation
func TestValidateIfOK(t *testing.T) {
	s := getValidStructToValidate()

	isValid, vContext := s.Validate()
	if !isValid {
		t.Fatal("Expected a valid struct")
	}

	expectedErrorCount := 0
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateIfOK checks valid structure validation if field checked by Range validator has border values
func TestValidateRangeIfBorderOK(t *testing.T) {
	s := getValidStructToValidate()

	s.checkRange = 2
	isValid, vContext := s.Validate()
	if !isValid {
		t.Fatal("Expected a valid struct")
	}

	expectedErrorCount := 0
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}

	vContext.Clear()
	s.checkRange = 5
	isValid, vContext = s.Validate()
	if !isValid {
		t.Fatal("Expected a valid struct")
	}

	expectedErrorCount = 0
	receivedErrorCount = len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateRequired checks validation for the required field
func TestValidateRequired(t *testing.T) {
	s := getValidStructToValidate()

	s.checkRequired = nil
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Expected invalid struct")
	}

	expectedErrorCount := 1
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}

	vContext.Clear()
	s.checkRequired = pString("")
	isValid, vContext = s.Validate()
	if isValid {
		t.Fatal("Expected invalid struct")
	}

	expectedErrorCount = 1
	receivedErrorCount = len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateMin checks validation of value lesser then min
func TestValidateMin(t *testing.T) {
	s := getValidStructToValidate()

	s.checkMin = 1
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Expected invalid struct")
	}

	expectedErrorCount := 1
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateMax checks validation of value higher then max
func TestValidateMax(t *testing.T) {
	s := getValidStructToValidate()

	s.checkMax = 6
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Expected invalid struct")
	}

	expectedErrorCount := 1
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateRange checks validation of out of range value
func TestValidateRange(t *testing.T) {
	s := getValidStructToValidate()

	s.checkRange = 1
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Expected invalid struct")
	}

	expectedErrorCount := 1
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}

	vContext.Clear()
	s.checkRange = 6
	isValid, vContext = s.Validate()
	if isValid {
		t.Fatal("Expected invalid struct")
	}

	expectedErrorCount = 1
	receivedErrorCount = len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateIfTwoErrors checks error details to contain multiple messages if such
func TestValidateIfTwoErrors(t *testing.T) {
	s := getValidStructToValidate()

	s.checkMin = 1
	s.checkMax = 6
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Expected invalid struct")
	}

	expectedErrorCount := 2
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Expected %d validation errors, received: %d", expectedErrorCount, receivedErrorCount)
	}
}
