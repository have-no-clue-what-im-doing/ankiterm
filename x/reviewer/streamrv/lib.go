package streamrv

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"golang.org/x/term"

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
		card, err := am.NextCard()
		if err != nil {
			if strings.Contains(err.Error(), "Gui review is not currently active") {
				fmt.Println("Congratulations! You have finished all cards.")
				return
			}
			panic(err)
		}

		clearScreen()

		fmt.Printf("\n[REVIEW MODE]\n")
		fmt.Println(format(card.Question))
		fmt.Println("\n[Press any key to Show Answer]")

		awaitAnyKey()
		fmt.Print("\n---\n")
		fmt.Println(format(card.Answer))

		lookup := []string{"Again", "Hard", "Good", "Easy"}
		for i, button := range card.Buttons {
			fmt.Printf("[%d] %s (%s)\n", button, lookup[i], card.NextReviews[i])
		}

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

// Clears the screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Reads a single key press without requiring ENTER
func awaitAnyKey() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	var b [1]byte
	os.Stdin.Read(b[:]) // Read one key
}

// Reads a single key press for selecting 1-4 without requiring ENTER
func awaitAction(validRange []int) reviewer.Action {
	fmt.Print("Enter choice (1-4): ")

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	var b [1]byte
	os.Stdin.Read(b[:]) // Read single character
	input := string(b[:])

	// Convert input to integer
	i, err := strconv.Atoi(input)
	if err != nil || !xslices.Contains(validRange, i) {
		fmt.Printf("\nInvalid input \"%s\", try again.\n", input)
		return awaitAction(validRange)
	}

	return reviewer.ActionFromString(input)
}

func format(text string) string {
	text = xmisc.PurgeStyle(text)
	text = xmisc.TtyColor(text)
	return text
}
