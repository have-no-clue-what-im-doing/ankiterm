package streamrv

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"os/exec"

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

	// Set up signal handling to exit cleanly on Ctrl+C
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)

	for {
		select {
		case <-signalChannel:
			fmt.Println("\nExiting...")
			return
		default:
			card, err := am.NextCard()
			if err != nil {
				if strings.Contains(err.Error(), "Gui review is not currently active") {
					fmt.Println("Congratulations! You have finished all cards.")
					return
				}
				panic(err)
			}
			if err != nil {
				fmt.Println("No more cards.")
				return
			}

			// Clear the screen after each card
			clearScreen()

			fmt.Printf("\n[REVIEW MODE]\n")
			fmt.Println(format(card.Question))
			fmt.Println("\n[Press any key to Show Answer]")
			awaitEnter()

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
}

// Function to clear the terminal screen
func clearScreen() {
	cmd := exec.Command("clear") // for Unix/Linux
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func awaitEnter() {
	// Wait for any key press to continue
	var input string
	fmt.Scanf("%s", &input)
}

func awaitAction(validRange []int) reviewer.Action {
	// No extra Enter key needed, input is captured directly
	fmt.Print("Enter choice (1-4): ")
	var input string
	fmt.Scanf("%s", &input)

	// Try to parse the input into an integer
	i, err := strconv.Atoi(input)
	if err != nil || !xslices.Contains(validRange, i) {
		fmt.Printf("Invalid input \"%s\" out of range, try again: \n", input)
		return awaitAction(validRange)
	}
	return reviewer.ActionFromString(input)
}

func format(text string) string {
	text = xmisc.PurgeStyle(text)
	text = xmisc.TtyColor(text)
	return text
}
