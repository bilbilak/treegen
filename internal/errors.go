package internal

import (
	"fmt"
	"os"

	app "github.com/bilbilak/treegen/config"
)

func Help() {
	fmt.Printf("%s v%s\n", app.Name, app.Version)
	fmt.Println("Run with --help for usage instructions.")
}

func FatalError(message ...error) {
	if len(message) > 0 {
		_, err := fmt.Fprintln(os.Stderr, message[0])

		if err != nil {
			fmt.Println("Something went wrong!")
		}

		os.Exit(1)
	} else {
		Help()
		os.Exit(2)
	}
}
