package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// easyjson -all stats.go
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
	rd := bufio.NewReader(r)
	user := User{}
	pattern := "." + domain
	var finish bool
	for !finish {
		bytes, err := rd.ReadBytes('\n')

		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err //nolint:wrapcheck
		}

		if errors.Is(err, io.EOF) {
			finish = true
		}

		err = unmarshal(bytes, &user)
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

		i, ok := result[domain]
		if !ok {
			result[domain] = 1
		} else {
			result[domain] = i + 1
		}
	}

	return result, nil
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
