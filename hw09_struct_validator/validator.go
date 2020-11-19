package hw09_struct_validator //nolint:golint,stylecheck
import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ErrValidate struct {
	Err error
	Msg string
}

func (v ErrValidate) Error() string {
	var sb strings.Builder
	sb.WriteString(v.Msg)
	sb.WriteString(":")
	sb.WriteString(v.Err.Error())
	return sb.String()
}

var (
	ErrNotStruct      = errors.New("input value is not a struct")
	ErrInvalidLength  = errors.New("length is invalid")
	ErrNotMatchRegexp = errors.New("string is not matched for regexp")
	ErrNotInSet       = errors.New("not included in set")
	ErrLess           = errors.New("less than the minimum")
	ErrMax            = errors.New("more than maximum")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder

	for _, err := range v {
		sb.WriteString("field: ")
		sb.WriteString(err.Field)
		sb.WriteString(", error: ")
		sb.WriteString(err.Err.Error())
	}
	return sb.String()
}

func Validate(v interface{}) (ValidationErrors, error) {
	var errs ValidationErrors
	valueV := reflect.ValueOf(v)
	if valueV.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	for i := 0; i < valueV.NumField(); i++ {
		value := valueV.Field(i)
		field := valueV.Type().Field(i)
		valErr, err := validateField(value, field)
		if err != nil {
			return nil, err
		}
		errs = append(errs, valErr...)
	}
	return errs, nil
}

func validateField(value reflect.Value, typeValue reflect.StructField) (ValidationErrors, error) {
	var errs ValidationErrors
	// потому что поддерживаем только 2 типа
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		return validateString(value, typeValue)
	case reflect.Int:
		return validateInt(value, typeValue)
	case reflect.Slice:
		return validateSlice(value, typeValue)
	default:
		return errs, nil
	}
}

func validateString(value reflect.Value, typeValue reflect.StructField) (ValidationErrors, error) {
	//   * `len:32` - длина строки должна быть ровно 32 символа;
	//    * `regexp:\\d+` - согласно регулярному выражению строка должна состоять из цифр
	//    (`\\` - экранирование слэша);
	//    * `in:foo,bar` - строка должна входить в множество строк {"foo", "bar"}.
	var errs ValidationErrors
	strValue := value.String()

	tag, ok := getTags(typeValue.Tag)
	if !ok {
		return nil, ErrValidate{Msg: "Ошибка парсинга тэгов", Err: nil}
	}

	lenString, err := validateLenString(tag, typeValue, strValue)
	if err != nil {
		return nil, err
	}
	errs = append(errs, lenString...)
	regexpString, err := validateRegexpString(tag, strValue, typeValue)
	if err != nil {
		return nil, err
	}
	errs = append(errs, regexpString...)

	inString, err := validateInString(tag, strValue, typeValue)
	if err != nil {
		return nil, err
	}
	errs = append(errs, inString...)

	return errs, nil
}

func validateInString(tag map[string]string, strValue string, typeValue reflect.StructField) (ValidationErrors, error) {
	var errs ValidationErrors
	inValueString, ok := tag["in"]
	if !ok {
		return errs, nil
	}
	inSlice := strings.Split(inValueString, ",")
	var findIn = false
	for _, inString := range inSlice {
		if strValue == inString {
			findIn = true
		}
	}
	if !findIn {
		errs = append(errs,
			ValidationError{Field: typeValue.Name, Err: ErrNotInSet})
	}

	return errs, nil
}

func validateRegexpString(tag map[string]string, strValue string, typeValue reflect.StructField) (ValidationErrors, error) {
	var errs ValidationErrors
	r, ok := tag["regexp"]
	if !ok {
		return errs, nil
	}
	var re, err = regexp.Compile(r)
	if err != nil {
		return nil, ErrValidate{Msg: "error compile", Err: err}
	}
	if !re.MatchString(strValue) {
		errs = append(errs,
			ValidationError{Field: typeValue.Name, Err: ErrNotMatchRegexp})
	}

	return errs, nil
}

func validateLenString(tag map[string]string, typeValue reflect.StructField, strValue string) (ValidationErrors, error) {
	var errs ValidationErrors
	reqLen, ok := tag["len"]
	if !ok {
		return errs, nil
	}

	lenInt, err := strconv.Atoi(reqLen)
	if err != nil {
		return nil, ErrValidate{Msg: "can not parse len value", Err: err}
	}

	if len(strValue) != lenInt {
		errs = append(errs,
			ValidationError{Field: typeValue.Name, Err: ErrInvalidLength})
	}

	return errs, nil
}

func getTags(tag reflect.StructTag) (map[string]string, bool) {
	var tags = make(map[string]string)
	lookup, ok := tag.Lookup("validate")
	if !ok {
		return nil, true
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

func validateInt(value reflect.Value, typeValue reflect.StructField) (ValidationErrors, error) {
	//    * `min:10` - число не может быть меньше 10;
	//    * `max:20` - число не может быть больше 20;
	//    * `in:256,1024` - число должно входить в множество чисел {256, 1024};
	var errs ValidationErrors
	iValue := value.Int()

	tag, ok := getTags(typeValue.Tag)
	if !ok {
		return nil, ErrValidate{Msg: "can not parse validate tag", Err: nil}
	}

	minError, err := validateMinInt(tag, typeValue, iValue)
	if err != nil {
		return nil, err
	}
	errs = append(errs, minError...)
	maxInt, err := validateMaxInt(tag, typeValue, iValue)
	if err != nil {
		return nil, err
	}
	errs = append(errs, maxInt...)
	inInt, err := validateInInt(tag, typeValue, iValue)
	if err != nil {
		return nil, err
	}
	errs = append(errs, inInt...)

	return errs, nil
}

func validateInInt(tag map[string]string, typeValue reflect.StructField, iValue int64) (ValidationErrors, error) {
	var errs ValidationErrors
	inValueString, ok := tag["in"]
	if !ok {
		return errs, nil
	}
	inSlice := strings.Split(inValueString, ",")
	var findIn = false
	for _, inString := range inSlice {
		inVal, err := strconv.ParseInt(inString, 10, 64)
		if err != nil {
			return nil, ErrValidate{Msg: "can not parse len value", Err: err}
		}
		if inVal == iValue {
			findIn = true
			break
		}
	}
	if !findIn {
		errs = append(errs, ValidationError{Field: typeValue.Name, Err: ErrNotInSet})
	}
	return errs, nil
}

func validateMaxInt(tag map[string]string, typeValue reflect.StructField, iValue int64) (ValidationErrors, error) {
	var errs ValidationErrors
	maxStr, ok := tag["max"]
	if !ok {
		return errs, nil
	}

	max, err := strconv.ParseInt(maxStr, 10, 64)
	if err != nil {
		return nil, ErrValidate{Msg: "can not parse max value", Err: err}
	}
	if iValue > max {
		errs = append(errs, ValidationError{Field: typeValue.Name, Err: ErrMax})
	}

	return errs, nil
}

func validateMinInt(tag map[string]string, typeValue reflect.StructField, iValue int64) (ValidationErrors, error) {
	var errs ValidationErrors
	minStr, ok := tag["min"]
	if !ok {
		return errs, nil
	}
	min, err := strconv.ParseInt(minStr, 10, 64)
	if err != nil {
		return nil, ErrValidate{Msg: "can not parse min value", Err: err}
	} else if iValue < min {
		errs = append(errs, ValidationError{Field: typeValue.Name, Err: ErrLess})
	}
	return errs, nil
}

func validateSlice(value reflect.Value, typeValue reflect.StructField) (ValidationErrors, error) {
	var errs ValidationErrors
	for i := 0; i < value.Len(); i++ {
		index := value.Index(i)
		valErr, err := validateField(index, typeValue)
		if err != nil {
			return nil, err
		}
		errs = append(errs, valErr...)
	}
	return errs, nil
}
