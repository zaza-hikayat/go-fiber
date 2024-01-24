package utility

import (
	"math/rand"
	"strconv"
)

func GenerateRandomNumber() string {
	// Generate a random number between 1000 and 9999 (inclusive)
	num := rand.Intn(9000) + 1000
	return strconv.Itoa(num)
}
