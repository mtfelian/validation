package validation

import (
	"fmt"
	"regexp"
)

// ValidationError represents validation error message
type ValidationError struct{ Message string }

// String returns a validation error message
func (e *ValidationError) String() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Validation represents validation errors list
type Validation struct{ Errors []*ValidationError }

// String returns a formatted list of validation errors separated by new lines
func (v Validation) String() string {
	s := ""
	for _, message := range v.Errors {
		s += fmt.Sprintf("%s\n", message.String())
	}
	return s
}

// Clear all validation errors
func (v *Validation) Clear() { v.Errors = []*ValidationError{} }

// HasErrors returns true if v contains at least one error
func (v *Validation) HasErrors() bool { return len(v.Errors) > 0 }

// Error adds an error message to a validation context
func (v *Validation) Error(message string, args ...interface{}) *ValidationResult {
	result := (&ValidationResult{Ok: false, Error: &ValidationError{}}).Message(message, args...)
	v.Errors = append(v.Errors, result.Error)
	return result
}

// ValidationResult is returned for every validation method
// It has Ok if validation succeeded or an error pointer
type ValidationResult struct {
	Error *ValidationError
	Ok    bool
}

// Message writes an error message to ValidationResult,
// Returns itself. Can be called like Sprintf()
func (r *ValidationResult) Message(message string, args ...interface{}) *ValidationResult {
	if r.Error == nil {
		return r
	}
	r.Error.Message = message
	if len(args) > 0 {
		r.Error.Message = fmt.Sprintf(message, args...)
	}
	return r
}

// Required checks the obj is not nil and is not zero value (empty string or slice)
func (v *Validation) Required(obj interface{}) *ValidationResult { return v.apply(Required{}, obj) }

// Min checks that n >= min
func (v *Validation) Min(n int, min int) *ValidationResult { return v.apply(Min{min}, n) }

// Max checks that n <= max
func (v *Validation) Max(n int, max int) *ValidationResult { return v.apply(Max{max}, n) }

// Range checks that min <= n <= max
func (v *Validation) Range(n, min, max int) *ValidationResult {
	return v.apply(Range{Min{min}, Max{max}}, n)
}

// MinSize checks that size of obj >= min
func (v *Validation) MinSize(obj interface{}, min int) *ValidationResult {
	return v.apply(MinSize{min}, obj)
}

// MaxSize checks that size of obj <= max
func (v *Validation) MaxSize(obj interface{}, max int) *ValidationResult {
	return v.apply(MaxSize{max}, obj)
}

// Length checks that length of obj == n
func (v *Validation) Length(obj interface{}, n int) *ValidationResult { return v.apply(Length{n}, obj) }

// Match checks that str matches regex
func (v *Validation) Match(str string, regex *regexp.Regexp) *ValidationResult {
	return v.apply(Match{regex}, str)
}

// Email checks that e-mail is valid
func (v *Validation) Email(email string) *ValidationResult {
	return v.apply(Email{Match{emailPattern}}, email)
}

// apply validator to obj
// и возвращает результат валидации
func (v *Validation) apply(chk Validator, obj interface{}) *ValidationResult {
	if chk.IsSatisfied(obj) {
		return &ValidationResult{Ok: true}
	}

	// add error to a validation context
	err := &ValidationError{Message: chk.DefaultMessage()}
	v.Errors = append(v.Errors, err)

	// and return it
	return &ValidationResult{Ok: false, Error: err}
}

// Check apply validator checks to obj and returns a ValidationResult from first validation error,
// or from last successful validation
func (v *Validation) Check(obj interface{}, checks ...Validator) *ValidationResult {
	var result *ValidationResult
	for _, check := range checks {
		result = v.apply(check, obj)
		if !result.Ok {
			return result
		}
	}
	return result
}
