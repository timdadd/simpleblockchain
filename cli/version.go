package cli

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

const Major = "0"
const Minor = "10" // Transparent Db
const Fix = "0"
const Verbal = "Peer-to-Peer DB Sync"

var ErrIncorrectUsage = errors.New("Incorrect Usage")

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Version: %s.%s.%s-beta %s", Major, Minor, Fix, Verbal))
	},
}

func IncorrectUsageErr() error {
	return ErrIncorrectUsage
}
