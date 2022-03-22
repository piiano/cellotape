package example

import (
	"fmt"
	"strings"
	"time"
)

// This is the service layer.
// The signatures of the functions here is pure to represent their purpose and is decoupled from any REST characteristics.
// These services are exposed as REST APIs only by their usage in a REST controller.
// They can be potentially be exposed in additional ways with other protocols, CLI, Golang SDK, etc.

func GreetV1(name string) (string, error) {
	return fmt.Sprintf("Hello %s!", name), nil
}

func GreetV2(name string, dayOfBirth time.Time, greetTemplate string) (string, error) {
	today := time.Now()
	if dayOfBirth.Month() == today.Month() && dayOfBirth.Day() == today.Day() {
		return fmt.Sprintf("Happy Birthday %s!", name), nil
	}
	if greetTemplate == "" {
		return fmt.Sprintf("Hello %s!", name), nil
	}
	if strings.Contains(greetTemplate, "%s") {
		return fmt.Sprintf(greetTemplate, name), nil
	}
	return greetTemplate, nil
}
