package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// easyjson -all stats.go.
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)

	user := User{}
	pattern := "." + domain
	for scanner.Scan() {
		bytes := scanner.Bytes()

		err := unmarshal(bytes, &user)
		if err != nil {
			return nil, err
		}

		if user.Email == "" {
			return nil, fmt.Errorf("not found email")
		}

		domain, err = getDomain(&user, pattern)

		if err != nil {
			return nil, err
		}

		if domain == "" {
			continue
		}

		result[domain]++

		resetEmail(&user)
	}

	return result, nil
}

func resetEmail(u *User) {
	u.Email = ""
}

func getDomain(user *User, pattern string) (string, error) {
	containsPoint := strings.Contains(user.Email, ".")
	if !containsPoint {
		return "", fmt.Errorf("not contain '.'")
	}

	hasSuffix := strings.HasSuffix(user.Email, pattern)
	if !hasSuffix {
		return "", nil
	}

	splitted := strings.Split(user.Email, "@")
	if len(splitted) != 2 {
		return "", fmt.Errorf("not contain @")
	}
	return strings.ToLower(splitted[1]), nil
}

func unmarshal(line []byte, u json.Unmarshaler) error {
	return u.UnmarshalJSON(line)
}
