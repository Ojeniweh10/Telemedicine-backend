package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenerateUUID(data string) string {
	// This function generates a UUID based on the input data.
	prefix := data
	if len(data) > 3 {
		prefix = data[:3]
	}
	rand.Seed(time.Now().UnixNano())
	suffix := fmt.Sprintf("%06d", rand.Intn(1000000))
	uuid := prefix + suffix
	return uuid
}

func GenerateUsertag(firstname string) string {
	reg := regexp.MustCompile(`[^a-zA-Z]+`)
	firstname = reg.ReplaceAllString(firstname, "")
	var wisetag string
	if len(firstname) >= 3 {
		wisetag = strings.ToUpper(firstname[:3])
	} else {
		wisetag = strings.ToUpper(firstname + strings.Repeat("X", 3-len(firstname)))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := rng.Intn(900000) + 100000
	return wisetag + strconv.Itoa(randomNumber)
}
