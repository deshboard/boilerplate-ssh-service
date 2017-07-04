package app

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// Application wraps a list of commands and handles their execution.
type Application struct {
	session  ssh.Session
	term     *terminal.Terminal
	prompt   string
	commands map[string]*cobra.Command
}

// NewApplication returns a new Application.
func NewApplication(session ssh.Session) *Application {
	prompt := fmt.Sprintf("%s@deshboard:$ ", session.User())
	term := terminal.NewTerminal(session, prompt)

	app := &Application{
		session: session,
		term:    term,
		prompt:  prompt,
	}

	app.commands = map[string]*cobra.Command{
		"help": &cobra.Command{
			Use:   "help",
			Short: "Shows the list of available commands",
			Run:   (&helpCommand{app}).Run,
		},
		"man": &cobra.Command{
			Use:   "man command",
			Short: "Shows the manual for a command",
			Run:   (&manCommand{app}).Run,
		},
		"hello": &cobra.Command{
			Use:   "hello",
			Short: "Asks your name and welcomes you",
			Run:   (&helloCommand{app}).Run,
		},
	}

	return app
}

// Run handles the main loop.
func (a *Application) Run() {
	for {
		line, err := a.term.ReadLine()

		// Ctrl+D received
		if err == io.EOF {
			io.WriteString(a.session, "\n")
			a.session.Exit(0)
		} else if err == nil {
			if line != "" {
				args := strings.Split(line, " ")
				a.Execute(args)
			}
		}
	}
}

// Execute handles the command execution.
func (a *Application) Execute(args []string) {
	if cmd, ok := a.commands[args[0]]; !ok {
		io.WriteString(a.session, fmt.Sprintf("command not found: %s\n", args[0]))
	} else {
		cmd.SetArgs(args[1:])
		cmd.SetOutput(a.session)
		cmd.Execute()
	}
}

// ReadInput temporarily changes the prompt and reads a line.
func (a *Application) ReadInput() (input string, err error) {
	a.term.SetPrompt("")
	input, err = a.term.ReadLine()
	a.term.SetPrompt(a.prompt)
	return
}

// helpCommand lists the available commands.
type helpCommand struct {
	app *Application
}

// Run lists the available commands.
func (c *helpCommand) Run(cmd *cobra.Command, args []string) {
	var keys []string
	for key := range c.app.commands {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// To perform the opertion you want
	for _, k := range keys {
		io.WriteString(c.app.session, fmt.Sprintf("%s - %s\n", k, c.app.commands[k].Short))
	}
}

// manCommand shows the manual for a command.
type manCommand struct {
	app *Application
}

// Run shows the manual for a command.
func (c *manCommand) Run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		io.WriteString(c.app.session, "What manual page do you want?\n")

		return
	}

	// Without this man man would end up in always showing manual for man no matter the command argument
	if args[0] == "man" {
		io.WriteString(c.app.session, "Usage: man [command]\n")

		return
	}

	if cmd, ok := c.app.commands[args[0]]; !ok {
		io.WriteString(c.app.session, fmt.Sprintf("No manual entry for %s\n", args[0]))
	} else {
		cmd.SetOutput(c.app.session)
		cmd.SetArgs([]string{args[0], "--help"})
		cmd.Execute()
	}
}

// helloCommand asks your name and welcomes you.
type helloCommand struct {
	app *Application
}

// Run asks your name and welcomes you.
func (c *helloCommand) Run(cmd *cobra.Command, args []string) {
	io.WriteString(c.app.session, "What is your name? ")

	name, err := c.app.ReadInput()

	if err != nil {
		io.WriteString(c.app.session, fmt.Sprintf("%v\n", err))
		return
	}

	io.WriteString(c.app.session, fmt.Sprintf("Welcome, %s!\n", name))
}
