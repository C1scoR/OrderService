package resilience

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

var myArr []string

func ProcessMessage(message string) error {
	// Simulate random failures
	if rand.Float32() < 0.3 {
		return errors.New("simulated processing failed")
	}
	fmt.Printf("Processed event: %v\n", message)
	//In the call func if err == nil, then use MoveToDLQ
	return nil
}

func MoveToDLQ(message string) {
	myArr = append(myArr, message)
}
