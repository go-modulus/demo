package framework

import "github.com/urfave/cli/v2"

type Commands struct {
	commands []*cli.Command
}

func NewCommands() *Commands {
	return &Commands{
		commands: make([]*cli.Command, 0),
	}
}

func (c *Commands) Add(command *cli.Command) {
	c.commands = append(c.commands, command)
}

func (c *Commands) GetAll() []*cli.Command {
	return c.commands
}

func (c *Commands) GetCommandByName(name string) *cli.Command {
	for _, command := range c.commands {
		if command.Name == name {
			return command
		}
	}

	return nil
}
