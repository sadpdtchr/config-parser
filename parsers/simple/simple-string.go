package simple

import (
	"fmt"

	"github.com/haproxytech/config-parser/common"
	"github.com/haproxytech/config-parser/errors"
)

type SimpleString struct {
	Enabled    bool
	Value      string
	Name       string
	SearchName string
	Comment    string
}

func (s *SimpleString) Init() {
	s.Enabled = false
	s.SearchName = s.Name
}

func (s *SimpleString) GetParserName() string {
	return s.SearchName
}

func (s *SimpleString) Parse(line string, parts, previousParts []string, comment string) (changeState string, err error) {
	if parts[0] == s.SearchName {
		if len(parts) < 2 {
			return "", &errors.ParseError{Parser: "SimpleString", Line: line, Message: "Parse error"}
		}
		s.Enabled = true
		s.Value = parts[1]
		s.Comment = comment
		return "", nil
	}
	return "", &errors.ParseError{Parser: s.SearchName, Line: line}
}

func (s *SimpleString) Valid() bool {
	if s.Enabled {
		return true
	}
	return false
}

func (s *SimpleString) Result(AddComments bool) []common.ReturnResultLine {
	if s.Enabled {
		return []common.ReturnResultLine{
			common.ReturnResultLine{
				Data:    fmt.Sprintf("%s %s", s.SearchName, s.Value),
				Comment: s.Comment,
			},
		}
	}
	return []common.ReturnResultLine{}
}
