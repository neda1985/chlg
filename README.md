# Description

Your goal is to create two functions in Golang that verifies the correctness of a supplied credit card number and determine its card scheme. You can decide on the function's signature and types. It is mandatory that you add automated tests.

Create a work branch for yourself. Create a pull request to the `main` branch when your solution is complete.

## Function #1: Validity of the card number

The following algorithm can be used to check validity of a card number:

1. Starting from the right, replace each **second** digit of the card number with its doubled value
2. When doubling a digit produces a 2-digit number (e.g 6 produces 12), then add those 2 digits (1+2 = 3)
3. Sum up all the digits

The card number is valid if the sum is divisible by 10

**Example**: Let's check if `5237 2516 2477 8133` is a valid credit card number.

1. Double each second digit: **10** 2 **6** 7 **4** 5 **2** 6 **4** 4 **14** 7 **16** 1 **6** 3
2. Add 2-digit numbers: **1** 2 6 7 4 5 2 6 4 4 **5** 7 **7** 1 6 3
3. Sum up all the digits: 70

70 is divisible by 10, so `5237 2516 2477 8133` is a **valid** credit card number

Please implement a function that given a credit card number returns if it is valid

## Function #2: Known/supported card schemes

Card Scheme (Visa, MasterCard, JCB, etc) can be detected by the first digits of the card and the length of the card.

**Example**

| Scheme           | Ranges           | Number of Digits | Example number   |
|---               |---               |---               |---
| American Express | 34,37            | 15               | 378282246310005  |
| JCB              | 3528-3589        | 16-19            | 3530111333300000 |
| Maestro          | 50, 56-58, 6     | 12-19            | 6759649826438453 |
| Visa             | 4                | 13,16,19         | 4012888888881881 |
| MasterCard       | 2221-2720, 51-55 | 16               | 5105105105105100 |


Please implement a function that given a credit card number returns its card scheme.

# How the tiny app work
**I was suppose to build two functions to do some validation, so I added some extra salt so it look more friendly :D**
 **Card Validator** is a Go CLI application that:
1. **Validates credit card numbers** using the Luhn algorithm.
2. **Detects the card scheme** (e.g., Visa, MasterCard, JCB).

---

## Features
- Validate credit card numbers.
- Identify card schemes based on card number ranges and length.
- Supports both **interactive mode** and **CLI mode** for batch processing.
- Graceful shutdown support.
---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/card-validator.git
   cd card-validator
   ```

2. Build the binary:
   ```bash
   go build -o card-validator
   ```

3. Run the application:
   ```bash
   ./card-validator
   ```
---
## Usage

### CLI Mode
Validate one or more card numbers directly from the command line:

```bash
./card-validator -l 4111111111111111,5105105105105100
```

**Example Output:**
```plaintext
Card 4111111111111111 validation is: true and card scheme is Visa
Card 5105105105105100 validation is: true and card scheme is MasterCard
```

### Interactive Mode
Run the program without arguments to validate cards interactively:

```bash
./card-validator
```

**Example Interaction:**
```plaintext
Welcome to the Card Validator!
Please enter a card number for validation or enter n to exit:
Card Number: 4111111111111111
Card validation results:
  - Valid: true
  - Card Scheme: Visa
Do you want to validate another card? Enter 1 to continue, or press 'n' to exit:
```

### Graceful Shutdown
Press `Ctrl+C` at any time to exit the program cleanly.


---

## Tests
Run unit tests to ensure functionality:

```bash
go test ./... -v
```

---