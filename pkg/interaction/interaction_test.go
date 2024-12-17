package interaction

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mocks "github.com/sumup-challenges/coding-challenge-op-go-neda1985/mocks"
	"github.com/sumup-challenges/coding-challenge-op-go-neda1985/pkg/validator"
)

func TestInteraction_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockValidator(ctrl)

	// mock expectations
	mockValidator.EXPECT().CardNumberValidator("valid").Return(true, nil).Times(1)
	mockValidator.EXPECT().CardSchemaValidator("valid").Return("Visa", nil).Times(1)
	mockValidator.EXPECT().CardNumberValidator("invalid").Return(false, nil).Times(1)
	mockValidator.EXPECT().CardSchemaValidator("invalid").Return("", validator.ErrUnknownScheme).Times(1)

	interaction := NewInteraction(mockValidator)

	input := "valid\n1\ninvalid\nn\n"

	stdinReader, stdinWriter, _ := os.Pipe()
	stdoutReader, stdoutWriter, _ := os.Pipe()

	oldStdin := os.Stdin
	oldStdout := os.Stdout

	os.Stdin = stdinReader
	os.Stdout = stdoutWriter

	// Ensure original stdin and stdout are restored after the test
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	// Writing the input to stdin
	go func() {
		stdinWriter.Write([]byte(input))
		stdinWriter.Close()
	}()

	// Capture the program output
	var output bytes.Buffer
	done := make(chan struct{})
	go func() {
		io.Copy(&output, stdoutReader)
		close(done)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	interaction.Run(ctx)

	stdoutWriter.Close()

	<-done
	expected := "Welcome to the Card Validator!\n" +
		"Please enter a card number for validation or enter n to exit:\n" +
		"Card Number: Card validation results:\n" +
		"  - Valid: true\n" +
		"  - Card Scheme: Visa\n" +
		"Do you want to validate another card? Enter 1 to continue, or press 'n' to exit:\n" +
		"Choice: Please enter a card number for validation or enter n to exit:\n" +
		"Card Number: Card scheme detection failed: unknown card scheme\n" +
		"Please enter a card number for validation or enter n to exit:\n" +
		"Card Number: Thank you for using the Card Validator!\n"

	if output.String() != expected {
		t.Errorf("expected %q, got %q", expected, output.String())
	}
}
