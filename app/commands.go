package app

import (
	"fmt"
	"io"
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
		session:  session,
		term:     term,
		prompt:   prompt,
		commands: make(map[string]*cobra.Command),
	}

	helpCmd := &helpCommand{
		session:  session,
		commands: app.commands,
	}

	manCmd := &manCommand{
		session:  session,
		commands: app.commands,
	}

	helloCmd := &helloCommand{
		session: session,
		term:    term,
		prompt:  prompt,
	}

	app.commands["help"] = &cobra.Command{
		Use:   "help",
		Short: "Lists the commands",
		Run:   helpCmd.Run,
	}

	app.commands["man"] = &cobra.Command{
		Use:   "man command",
		Short: "Shows the manual for a command",
		Run:   manCmd.Run,
	}

	app.commands["hello"] = &cobra.Command{
		Use:   "hello",
		Short: "Asks your name and welcomes you",
		Run:   helloCmd.Run,
	}

	return app
}

// Execute handles the command execution.
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

// helpCommand lists the available commands.
type helpCommand struct {
	session  ssh.Session
	commands map[string]*cobra.Command
}

// Run lists the available commands.
func (c *helpCommand) Run(cmd *cobra.Command, args []string) {
	for name, cmd := range c.commands {
		io.WriteString(c.session, fmt.Sprintf("%s - %s\n", name, cmd.Short))
	}
}

// manCommand shows the manual for a command.
type manCommand struct {
	session  ssh.Session
	commands map[string]*cobra.Command
}

// Run shows the manual for a command.
func (c *manCommand) Run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		io.WriteString(c.session, "What manual page do you want?\n")

		return
	}

	if cmd, ok := c.commands[args[0]]; !ok {
		io.WriteString(c.session, fmt.Sprintf("No manual entry for %s\n", args[0]))
	} else {
		cmd.SetOutput(c.session)
		cmd.SetArgs([]string{args[0], "--help"})
		cmd.Execute()
	}
}

// helloCommand asks your name and welcomes you.
type helloCommand struct {
	session ssh.Session
	term    *terminal.Terminal
	prompt  string
}

// Run asks your name and welcomes you.
func (c *helloCommand) Run(cmd *cobra.Command, args []string) {
	io.WriteString(c.session, "What is your name? ")

	c.term.SetPrompt("")
	name, err := c.term.ReadLine()
	c.term.SetPrompt(c.prompt)

	if err != nil {
		io.WriteString(c.session, fmt.Sprintf("%v\n", err))
		return
	}

	io.WriteString(c.session, fmt.Sprintf("Welcome, %s!\n", name))
}
