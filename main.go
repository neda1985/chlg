package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sumup-challenges/coding-challenge-op-go-neda1985/pkg/interaction"
	"github.com/sumup-challenges/coding-challenge-op-go-neda1985/pkg/validator"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("loop number %v\n", i)
		}()
	}
	wg.Wait()

	//newValidator := validator.NewValidator()
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//var target string
	//flag.StringVar(&target, "l", "", "Enter a card number to validate (comma-separated for multiple)")
	//flag.Parse()
	//
	//cardList := strings.Split(target, ",")
	//
	//// Handle the CLI mode with -l flag
	//if target != "" {
	//	handleCLIMode(ctx, newValidator, cardList)
	//	return
	//}
	//
	//// Handle interactive mode
	//handleInteractiveMode(ctx, newValidator, cancel)
}

// handleCLIMode processes card numbers provided via CLI
func handleCLIMode(ctx context.Context, validator validator.Validator, cardList []string) {
	var wg sync.WaitGroup

	for _, card := range cardList {
		wg.Add(1)
		go func(card string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				fmt.Printf("Validation for card %s interrupted.\n", card)
				return
			default:
				isValid, err := validator.CardNumberValidator(card)
				if err != nil {
					fmt.Printf("Validation failed for card %s: %v\n", card, err)
					return
				}
				cardScheme, err := validator.CardSchemaValidator(card)
				if err != nil {
					fmt.Printf("Card scheme detection failed for card %s: %v\n", card, err)
					return
				}
				fmt.Printf("Card %s validation is: %v and card scheme is %s\n", card, isValid, cardScheme)
			}
		}(card)
	}

	wg.Wait()
}

// handleInteractiveMode starts the interactive user validation process
func handleInteractiveMode(ctx context.Context, validator validator.Validator, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	go func() {
		<-signalChan
		fmt.Println("\nProgram interrupted. Exiting gracefully...")
		cancel()
	}()
	newInteraction := interaction.NewInteraction(validator)
	newInteraction.Run(ctx)
}
