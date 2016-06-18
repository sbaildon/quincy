package commands

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Execute(*discordgo.Session, *discordgo.Message)
}

// Commands tracks loaded commands.
var Commands map[string]Command

func init() {
	Commands = make(map[string]Command)
}

func HandleCommand(session *discordgo.Session, message *discordgo.Message) {
	input := strings.Split(message.Content, " ")[0]

	if !validCommand(input) {
		return
	}

	command, err := buildCommand(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	command.Execute(session, message)
}

func validCommand(input string) bool {
	if input[0] == '!' {
		return true
	}

	return false
}

func buildCommand(input string) (Command, error) {
	cmd := input[1:len(input)]

	comm, err := ExpandCommand(cmd)
	if err != nil {
		return nil, err
	}

	builder := Commands[comm]

	return builder, nil
}

// ExpandCommand expands the passed in command to the full value.
func ExpandCommand(command string) (string, error) {
	names := CommandNames()
	r := regexp.MustCompile(`^` + command)

	validCommands := []string{}
	for _, n := range names {
		// Exact match returns immediately.
		if n == command {
			return n, nil
		}

		if r.Match([]byte(n)) {
			validCommands = append(validCommands, n)
		}
	}

	switch len(validCommands) {
	case 0:
		return "", fmt.Errorf("No command found for %q", command)
	case 1:
		return validCommands[0], nil
	default:
		return "", fmt.Errorf("Multiple commands matched %q: %v", command, validCommands)
	}
}

// AddCommand should be called within your commands's init() func.
// This will register the command so it can be used.
func AddCommand(name string, command Command) {
	Commands[name] = command
}

// DisplayCommands displays all the loaded commands.
func DisplayCommands() string {
	names := CommandNames()
	return fmt.Sprintf("%s\n", strings.Join(names, "\n"))
}

// CommandNames returns a sorted slice of command names.
func CommandNames() []string {
	names := []string{}

	for key, _ := range Commands {
		names = append(names, key)
	}

	sort.Strings(names)
	return names
}

