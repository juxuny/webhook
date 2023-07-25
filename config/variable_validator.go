package config

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func VariableValidate(input string, variableType VariableType) error {
	if variableType == VariableTypeString {
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
