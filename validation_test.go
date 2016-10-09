package validation

import (
	"testing"
)

// structToValidate это структура для тестов валидации
type structToValidate struct {
	checkRequired *string
	checkMin      int
	checkMax      int
	checkRange    int
}

// Validate валидирует структуру structToValidate для тестирования
func (s structToValidate) Validate() (bool, *Validation) {
	v := &Validation{}
	v.Check(s.checkRequired, Required{}).Message("Поле checkRequired не может быть пустым")
	v.Check(s.checkMin, Min{2}).Message("Значение поля checkMin слишком мало")
	v.Check(s.checkMax, Max{5}).Message("Значение поля checkMax слишком велико")
	v.Check(s.checkRange, Range{Min{2}, Max{5}}).Message("Значение поля checkRange вне диапазона")
	return !v.HasErrors(), v
}

// getValidStructToValidate возвращает валидную структуру для тестов валидации
func getValidStructToValidate() structToValidate {
	return structToValidate{
		checkRequired: pString("q"),
		checkMin:      2,
		checkMax:      5,
		checkRange:    3,
	}
}

// pString возвращает указатель на строку s
func pString(s string) *string {
	return &s
}

// TestValidateIfOK тестирует валидацию валидной структуры
func TestValidateIfOK(t *testing.T) {
	s := getValidStructToValidate()

	isValid, vContext := s.Validate()
	if !isValid {
		t.Fatal("Ожидалась валидная структура")
	}

	expectedErrorCount := 0
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateIfOK тестирует валидацию валидной структуры где поле,
// проверяемое валидаторов Range имеет граничные значения
func TestValidateRangeIfBorderOK(t *testing.T) {
	s := getValidStructToValidate()

	s.checkRange = 2
	isValid, vContext := s.Validate()
	if !isValid {
		t.Fatal("Ожидалась валидная структура")
	}

	expectedErrorCount := 0
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}

	vContext.Clear()
	s.checkRange = 5
	isValid, vContext = s.Validate()
	if !isValid {
		t.Fatal("Ожидалась валидная структура")
	}

	expectedErrorCount = 0
	receivedErrorCount = len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateRequired тестирует проверку на отсутствие необходимого поля
func TestValidateRequired(t *testing.T) {
	s := getValidStructToValidate()

	s.checkRequired = nil
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Ожидалась невалидная структура")
	}

	expectedErrorCount := 1
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}

	vContext.Clear()
	s.checkRequired = pString("")
	isValid, vContext = s.Validate()
	if isValid {
		t.Fatal("Ожидалась невалидная структура")
	}

	expectedErrorCount = 1
	receivedErrorCount = len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateMin тестирует проверку на значение меньше минимального
func TestValidateMin(t *testing.T) {
	s := getValidStructToValidate()

	s.checkMin = 1
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Ожидалась невалидная структура")
	}

	expectedErrorCount := 1
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateMax тестирует проверку на значение больше максимального
func TestValidateMax(t *testing.T) {
	s := getValidStructToValidate()

	s.checkMax = 6
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Ожидалась невалидная структура")
	}

	expectedErrorCount := 1
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateRange тестирует проверку на значение вне диапазона
func TestValidateRange(t *testing.T) {
	s := getValidStructToValidate()

	s.checkRange = 1
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Ожидалась невалидная структура")
	}

	expectedErrorCount := 1
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}

	vContext.Clear()
	s.checkRange = 6
	isValid, vContext = s.Validate()
	if isValid {
		t.Fatal("Ожидалась невалидная структура")
	}

	expectedErrorCount = 1
	receivedErrorCount = len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}
}

// TestValidateIfTwoErrors проверяет что есть сообщения об ошибках в случае нескольких ошибок
func TestValidateIfTwoErrors(t *testing.T) {
	s := getValidStructToValidate()

	s.checkMin = 1
	s.checkMax = 6
	isValid, vContext := s.Validate()
	if isValid {
		t.Fatal("Ожидалась невалидная структура")
	}

	expectedErrorCount := 2
	receivedErrorCount := len(vContext.Errors)
	if receivedErrorCount != expectedErrorCount {
		t.Fatalf("Ожидалось %d ошибок валидации, получено: %d", expectedErrorCount, receivedErrorCount)
	}
}
