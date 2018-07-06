package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rivo/tview"
	"github.com/sankt-petersbug/httptop"
	"github.com/spf13/cobra"
	//"github.com/andreyvit/diff"
)

type flagOptions struct {
	fpath     string
	rateLimit int
}

func NewRootCommand() *cobra.Command {
	opt := flagOptions{}

	cmd := &cobra.Command{
		Use:   "httptop",
		Short: "Monitor HTTP logs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(opt.fpath, opt.rateLimit); err != nil {
				return err
			}

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opt.fpath, "fpath", "f", "/var/log/access.log", "Log filepath (default is /var/log/access.log)")
	flags.IntVarP(&opt.rateLimit, "ratelimit", "r", 10, "Request rate limit (default is 10 req/s)")

	return cmd
}

func run(fpath string, rateLimit int) error {
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("File not found at %s", fpath))
	}

	t, err := httptop.ReadFile(fpath)
	if err != nil {
		return err
	}

	messageStream := httptop.BatchRead(httptop.ToRecord(t.Lines), 5*time.Second)
	app := tview.NewApplication()
	root := httptop.NewLayout(rateLimit, fpath)

	go func() {
		for msg := range messageStream {
			root.Update(msg)

			app.Draw()
		}
	}()

	if err := app.SetRoot(root.GetView(), true).Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	rootCmd := NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	return
}
