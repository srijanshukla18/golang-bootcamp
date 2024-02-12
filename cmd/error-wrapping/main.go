package main

import (
	"errors"
	"fmt"
)

// functionC generates an error
func functionC() error {
	return errors.New("original error from C")
}

// functionB wraps the error from functionC
func functionB() error {
	err := functionC()
	if err != nil {
		// Wrap the error with additional context
		return fmt.Errorf("functionB failed: %w", err)
	}
	return nil
}

// functionA calls functionB and unwraps the error
func functionA() {
	err := functionB()
	if err != nil {
		fmt.Println("Error received in functionA:", err)

		// Unwrap the error to get the original error from functionC
		unwrappedErr := errors.Unwrap(err)
		if unwrappedErr != nil {
			fmt.Println("Unwrapped error:", unwrappedErr)
		}

		// Check if the unwrapped error is the specific error from functionC
		if errors.Is(err, unwrappedErr) {
			fmt.Println("The unwrapped error is the original error from functionC")
		}
	}
}

func main() {
	functionA()
}
