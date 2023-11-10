package config

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	allowCharactor = "abcdevghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_=+"
)

var allowCharactorMapper map[rune]bool

func init() {
	for _, c := range allowCharactor {
		allowCharactorMapper[c] = true
	}
}

func VariableValidate(input string, variableType VariableType) error {
	if variableType == VariableTypeString {
		for _, c := range input {
			if _, b := allowCharactorMapper[c]; !b {
				return errors.Errorf("invalid charactor %c", c)
			}
		}
		return nil
	}
	if variableType == VariableTypeBool {
		input = strings.ToLower(input)
		if input == "1" || input == "0" || input == "false" || input == "true" {
			return nil
		} else {
			return errors.New("invalid bool: " + input)
		}
	} else if variableType == VariableTypeInteger || variableType == VariableTypeInt64 {
		_, err := strconv.ParseInt(input, 10, 64)
		return err
	} else if variableType == VariableTypeNumber {
		_, err := strconv.ParseFloat(input, 64)
		return err
	}
	return nil
}
