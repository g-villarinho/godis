package protocol

import (
	"fmt"
	"strings"

	"github.com/g-villarinho/godis/internal/core/domain/command"
)

type Command struct {
	Type command.Type
	Args []string
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseCommand(line string) (*Command, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, fmt.Errorf("empty command")
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	commandType := command.Type(strings.ToUpper(parts[0]))
	if !commandType.IsValid() {
		return nil, fmt.Errorf("unknown command: %s", parts[0])
	}

	command := &Command{
		Type: commandType,
		Args: []string{},
	}

	if len(parts) > 1 {
		command.Args = parts[1:]
	}

	return command, nil
}

func (p *Parser) FormatResponse(result any) string {
	switch v := result.(type) {
	case string:
		return v
	case int, int64:
		return fmt.Sprintf("%d", v)
	case bool:
		if v {
			return p.FormatOK()
		}
		return "ERR Operation failed"
	case error:
		return p.FormatError(v.Error())
	default:
		return fmt.Sprintf("%v", result)
	}
}

func (p *Parser) FormatOK() string {
	return "OK"
}

func (p *Parser) FormatError(msg string) string {
	return fmt.Sprintf("ERR: %s", msg)
}

func (p *Parser) FormatNil() string {
	return "nil"
}
