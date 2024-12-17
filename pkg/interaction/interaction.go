package interaction

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/sumup-challenges/coding-challenge-op-go-neda1985/pkg/validator"
)

// Interaction handles user interaction for card validation
type Interaction struct {
	Validator validator.Validator
}

// NewInteraction creates a new Interaction instance
func NewInteraction(validator validator.Validator) *Interaction {
	return &Interaction{Validator: validator}
}

// Run starts the CLI interaction loop
func (i *Interaction) Run(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Card Validator!")

	for {
		select {
		case <-ctx.Done():
			// Exit immediately if the context is canceled
			fmt.Println("Session timed out or interrupted. Exiting...")
			return
		default:
			// Prompt user for input
			fmt.Println("Please enter a card number for validation or enter n to exit:")
			fmt.Print("Card Number: ")

			if err := waitForInputOrCancel(ctx, scanner); err != nil {
				fmt.Println(err)
				return
			}

			cardNumber := scanner.Text()
			if cardNumber == "n" || cardNumber == "N" {
				fmt.Println("Thank you for using the Card Validator!")
				return
			}

			isValid, err := i.Validator.CardNumberValidator(cardNumber)
			if err != nil {
				fmt.Printf("Card number validation failed: %v\n", err)
				continue
			}
			cardScheme, err := i.Validator.CardSchemaValidator(cardNumber)
			if err != nil {
				fmt.Printf("Card scheme detection failed: %v\n", err)
				continue
			}

			fmt.Println("Card validation results:")
			fmt.Printf("  - Valid: %v\n", isValid)
			fmt.Printf("  - Card Scheme: %s\n", cardScheme)

			fmt.Println("Do you want to validate another card? Enter 1 to continue, or press 'n' to exit:")
			fmt.Print("Choice: ")

			if err := waitForInputOrCancel(ctx, scanner); err != nil {
				fmt.Println(err)
				return
			}

			if scanner.Text() != "1" {
				fmt.Println("Thank you for using the Card Validator!")
				return
			}
		}
	}
}

func waitForInputOrCancel(ctx context.Context, scanner *bufio.Scanner) error {
	done := make(chan struct{})
	go func() {
		scanner.Scan()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("session timed out or interrupted. Exiting")
	case <-done:
		return nil
	}
}
