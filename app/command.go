package app

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Command configures a command to run.
type Command struct {
	Name         string
	key          key.Binding
	Cmd          func(string) tea.Cmd
	ValidateArgs func(string) bool
	async        bool
}

// ExecCmdErrMsg is a tea.Msg of a command's error.
type ExecCmdErrMsg struct {
	err error
}

// NewCommand initializes a Command.
func NewCommand(name string) *Command {
	return &Command{
		Name: name,
		Cmd:  NewStatusMessage,
	}
}

// BlockingCommand initializes a Command that will exec in a blocking fashion.
func BlockingCommand(name string) *Command {
	c := NewCommand(name)
	c.Cmd = c.Exec
	return c
}

// AsyncCommand initializes a Command that will exec in asynchronous fashion.
func AsyncCommand(name string) *Command {
	c := NewCommand(name)
	c.Cmd = c.Exec
	c.async = true
	return c
}

// Key sets the keybind for the command (if any).
func (c *Command) Key(k key.Binding) *Command {
	c.key = k
	return c
}

// Exec runs the command, optionally validating the args if ValidateArgs is set.
func (c *Command) Exec(args string) tea.Cmd {
	if !c.ArgsValid(args) {
		return func() tea.Msg {
			return ExecCallback(errors.New("invalid args"))
		}
	}
	if c.async {
		return c.asyncCmd(args)
	}
	return c.teaCmd(args)
}

func (c *Command) asyncCmd(args string) tea.Cmd {
	return func() tea.Msg {
		var (
			ste strings.Builder
			sto strings.Builder
			out string
		)

		cmd := exec.Command(c.Name, args)
		cmd.Stderr = &ste
		cmd.Stdout = &sto

		err := cmd.Run()
		if err != nil {
			out = ste.String()
		}
		if o := sto.String(); o != "" {
			out = o
		}

		o := strings.Join(strings.Split(out, "\n"), " ")

		return ExecCallback(errors.Join(err, errors.New(o)))
	}
}

func (c Command) teaCmd(args string) tea.Cmd {
	cmd := exec.Command(c.Name, args)
	return tea.ExecProcess(cmd, ExecCallback)
}

// ArgsValid validates the args using the Command's ValidateArgs function.
func (c Command) ArgsValid(args string) bool {
	if c.ValidateArgs != nil {
		return c.ValidateArgs(args)
	}
	return true
}

// ExecCallback returns a tea.Msg for the Command's error.
func ExecCallback(err error) tea.Msg {
	return ExecCmdErrMsg{
		err: err,
	}
}
