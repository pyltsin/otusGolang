package hw09_struct_validator //nolint:golint,stylecheck
import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder

	for _, err := range v {
		sb.WriteString(fmt.Sprintf("Error %s for field %s \n", err.Err, err.Field))
	}
	return sb.String()
}

func Validate(v interface{}) ValidationErrors {
	var errors ValidationErrors
	valueV := reflect.ValueOf(v)
	for i := 0; i < valueV.NumField(); i++ {
		value := valueV.Field(i)
		field := valueV.Type().Field(i)
		err := validateField(value, field)
		errors = append(errors, err...)
	}
	return errors
}

func validateField(value reflect.Value, typeValue reflect.StructField) ValidationErrors {
	var errors ValidationErrors
	switch value.Kind() { //nolint:exhaustive потому что поддерживаем только 2 типа
	case reflect.String:
		return validateString(value, typeValue)
	case reflect.Int:
		return validateInt(value, typeValue)
	case reflect.Slice:
		return validateSlice(value, typeValue)
	default:
		return errors
	}
}

func validateString(value reflect.Value, typeValue reflect.StructField) ValidationErrors {
	//   * `len:32` - длина строки должна быть ровно 32 символа;
	//    * `regexp:\\d+` - согласно регулярному выражению строка должна состоять из цифр
	//    (`\\` - экранирование слэша);
	//    * `in:foo,bar` - строка должна входить в множество строк {"foo", "bar"}.
	var errors ValidationErrors
	strValue := value.String()

	tag, ok := getTags(typeValue.Tag)
	if !ok {
		return errors
	}

	errors = append(errors, validateLenString(tag, typeValue, strValue)...)
	errors = append(errors, validateRegexpString(tag, strValue, typeValue)...)
	errors = append(errors, validateInString(tag, strValue, typeValue)...)

	return errors
}

func validateInString(tag map[string]string, strValue string, typeValue reflect.StructField) ValidationErrors {
	var errors ValidationErrors
	inValueString, ok := tag["in"]
	if ok {
		inSlice := strings.Split(inValueString, ",")
		var findIn = false
		for _, inString := range inSlice {
			if strValue == inString {
				findIn = true
			}
		}
		if !findIn {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("%v not in %s", strValue, inValueString)})
		}
	}
	return errors
}

func validateRegexpString(tag map[string]string, strValue string, typeValue reflect.StructField) ValidationErrors {
	var errors ValidationErrors
	r, ok := tag["regexp"]
	if ok {
		var re = regexp.MustCompile(r)
		if !re.MatchString(strValue) {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("%s can re: %s", strValue, r)})
		}
	}
	return errors
}

func validateLenString(tag map[string]string, typeValue reflect.StructField, strValue string) ValidationErrors {
	var errors ValidationErrors
	reqLen, ok := tag["len"]
	if ok {
		lenInt, err := strconv.Atoi(reqLen)
		if err != nil {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("can not parse len value: %s", reqLen)})
		} else if len(strValue) != lenInt {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("len %v  != %v", strValue, lenInt)})
		}
	}
	return errors
}

func getTags(tag reflect.StructTag) (map[string]string, bool) {
	var tags = make(map[string]string)
	lookup, ok := tag.Lookup("validate")
	if !ok {
		return nil, false
	}

	for _, tag := range strings.Split(lookup, "|") {
		splittedTag := strings.Split(tag, ":")
		if len(splittedTag) != 2 {
			return nil, false
		}
		tags[splittedTag[0]] = splittedTag[1]
	}
	return tags, true
}

func validateInt(value reflect.Value, typeValue reflect.StructField) ValidationErrors {
	//    * `min:10` - число не может быть меньше 10;
	//    * `max:20` - число не может быть больше 20;
	//    * `in:256,1024` - число должно входить в множество чисел {256, 1024};
	var errors ValidationErrors
	iValue := value.Int()

	tag, ok := getTags(typeValue.Tag)
	if !ok {
		errors = append(errors,
			ValidationError{Field: "", Err: fmt.Errorf("can not parse validate tag: %s", typeValue.Tag)})
		return errors
	}

	errors = append(errors, validateMinInt(tag, typeValue, iValue)...)
	errors = append(errors, validateMaxInt(tag, typeValue, iValue)...)
	errors = append(errors, validateInInt(tag, typeValue, iValue)...)

	return errors
}

func validateInInt(tag map[string]string, typeValue reflect.StructField, iValue int64) ValidationErrors {
	var errors ValidationErrors
	inValueString, ok := tag["in"]
	if ok {
		inSlice := strings.Split(inValueString, ",")
		var findIn = false
		for _, inString := range inSlice {
			inVal, err := strconv.Atoi(inString)
			if err != nil {
				errors = append(errors,
					ValidationError{Field: typeValue.Name, Err: fmt.Errorf("can not parse in value: %s", inValueString)})
			} else if int64(inVal) == iValue {
				findIn = true
				break
			}
		}
		if !findIn {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("%v not in %s", iValue, inValueString)})
		}
	}
	return errors
}

func validateMaxInt(tag map[string]string, typeValue reflect.StructField, iValue int64) ValidationErrors {
	var errors ValidationErrors
	maxStr, ok := tag["max"]
	if ok {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("can not parse max value: %s", maxStr)})
		} else if iValue > int64(max) {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("%v less then %v", iValue, max)})
		}
	}
	return errors
}

func validateMinInt(tag map[string]string, typeValue reflect.StructField, iValue int64) ValidationErrors {
	var errors ValidationErrors
	minStr, ok := tag["min"]
	if ok {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("can not parse min value: %s", minStr)})
		} else if iValue < int64(min) {
			errors = append(errors,
				ValidationError{Field: typeValue.Name, Err: fmt.Errorf("%v less then %v", iValue, min)})
		}
	}
	return errors
}

func validateSlice(value reflect.Value, typeValue reflect.StructField) ValidationErrors {
	var errors ValidationErrors
	for i := 0; i < value.Len(); i++ {
		index := value.Index(i)
		err := validateField(index, typeValue)
		errors = append(errors, err...)
	}
	return errors
}
