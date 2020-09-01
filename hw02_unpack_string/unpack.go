package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(text string) (string, error) {
	// Place your code here

	builder := strings.Builder{}

	var buffer = strings.Builder{}

	for _, r := range text { //если цифра, то повторяем buffer
		//если есть символ экранирования
		if buffer.String() == "\\" {
			//экранировать можно только цифры и слэш
			if !(unicode.IsDigit(r) || string(r) == "\\") {
				return "", ErrInvalidString
			}
			if unicode.IsDigit(r) {
				buffer.Reset()
			}
			buffer.WriteRune(r)
			continue
		}

		if unicode.IsDigit(r) {
			//если повторять нечего значит ошибка - два подряд числа
			if len(getNextLetter(buffer)) == 0 {
				return "", ErrInvalidString
			}
			//делаем набор
			countRepeat, _ := strconv.Atoi(string(r))
			builder.WriteString(strings.Repeat(getNextLetter(buffer), countRepeat))
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

func getNextLetter(buffer strings.Builder) string {
	nextLetter := buffer.String()
	if nextLetter == "\\\\" {
		return "\\"
	}
	return nextLetter
}
