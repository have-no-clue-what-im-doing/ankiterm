package streamrv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pluveto/ankiterm/x/automata"
	"github.com/pluveto/ankiterm/x/reviewer"
	"github.com/pluveto/ankiterm/x/xmisc"
	"github.com/pluveto/ankiterm/x/xslices"
)

func Execute(am *automata.Automata, deck string) {
	if am == nil {
		panic("am (automata.Automata) is nil")
	}
	if deck == "" {
		panic("deck is empty")
	}

	err := am.StartReview(deck)
	if err != nil {
		panic(err)
	}
	defer am.StopReview()

	for {
		// Clear the screen before displaying the next card
		clearScreen()

		card, err := am.NextCard()
		if err != nil {
			if strings.Contains(err.Error(), "Gui review is not currently active") {
				fmt.Println("Congratulations! You have finished all cards.")
				return
			}
			panic(err)
		}

		if card == nil {
			fmt.Println("No more cards.")
			return
		}

		// Print question only once
		fmt.Printf("\n[REVIEW MODE]\n")
		fmt.Println(format(card.Question))
		fmt.Println("\n[Press any key to Show Answer]")
		awaitEnter() // Wait for input to show answer

		// After input, print answer
		fmt.Print("\n---\n")
		fmt.Println(format(card.Answer))

		// Show answer options
		lookup := []string{"Again", "Hard", "Good", "Easy"}
		for i, button := range card.Buttons {
			fmt.Printf("[%d] %s (%s)\n", button, lookup[i], card.NextReviews[i])
		}

		// Get user action for card answer
		action := awaitAction(am.CurrentCard().Buttons)
		switch code := action.GetCode(); code {
		case reviewer.ActionAbort:
			return
		case reviewer.ActionSkip:
			continue
		case reviewer.ActionAnswer:
			am.AnswerCard(action.(reviewer.AnswerAction).CardEase)
		default:
			panic("unknown action code")
		}
	}
}

// Function to clear the terminal screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func awaitEnter() {
	// Wait for any key press to continue
	var input string
	fmt.Scanln(&input)
}

func awaitAction(validRange []int) reviewer.Action {
	// Wait for user to select action
	print("awaitAction")
	var input string
	fmt.Scanln(&input)

	// Try to parse input into an integer
	i, err := strconv.Atoi(input)
	if err != nil || !xslices.Contains(val
