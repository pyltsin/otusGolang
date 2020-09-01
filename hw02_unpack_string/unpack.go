package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const escaped = "\\"
const recordedEscaped = "\\\\"
const emptyString = ""

func Unpack(text string) (string, error) {
	// Place your code here

	builder := strings.Builder{}

	var buffer = &strings.Builder{}

	for _, r := range text { //если цифра, то повторяем buffer
		//если есть символ экранирования
		if isEscaped(buffer) {
			//экранировать можно только цифры и слэш
			if !isDigitOrEscaped(r) {
				return emptyString, ErrInvalidString
			}
			if isDigit(r) {
				buffer.Reset()
			}
			buffer.WriteRune(r)
			continue
		}

		if unicode.IsDigit(r) {
			//если повторять нечего значит ошибка - два подряд числа
			if isEmpty(buffer) {
				return emptyString, ErrInvalidString
			}
			//делаем набор
			countRepeat, _ := strconv.Atoi(string(r))
			repeatedSymbol := strings.Repeat(getNextLetter(buffer), countRepeat)
			builder.WriteString(repeatedSymbol)
			buffer.Reset()
			continue
		}

		//если буква, то запихиваем в буффер, а оттуда заполняем builder
		builder.WriteString(getNextLetter(buffer))
		buffer.Reset()
		buffer.WriteRune(r)
	}

	builder.WriteString(buffer.String())

	return builder.String(), nil
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isDigitOrEscaped(r rune) bool {
	return unicode.IsDigit(r) || string(r) == escaped
}

func isEscaped(buffer *strings.Builder) bool {
	return buffer.String() == escaped
}

func isEmpty(buffer *strings.Builder) bool {
	return len(getNextLetter(buffer)) == 0
}

func getNextLetter(buffer *strings.Builder) string {
	nextLetter := buffer.String()
	if nextLetter == recordedEscaped {
		return escaped
	}
	return nextLetter
}
