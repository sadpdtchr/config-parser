package tcp

import (
	"fmt"

	"github.com/haproxytech/config-parser/common"
	"github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/parsers/tcp/actions"
)

type TCPRequests struct {
	Mode string //frontent, backend
	data []TCPAction
}

func (h *TCPRequests) Init() {
	h.data = []TCPAction{}
}

func (h *TCPRequests) GetParserName() string {
	return "tpc-request"
}

func (h *TCPRequests) Clear() {
	h.Init()
}

func (h *TCPRequests) Get(createIfNotExist bool) (common.ParserData, error) {
	if len(h.data) == 0 && !createIfNotExist {
		return nil, &errors.FetchError{}
	}
	return h.data, nil
}

func (h *TCPRequests) Set(data common.ParserData) error {
	switch newValue := data.(type) {
	case []TCPAction:
		h.data = newValue
	case TCPAction:
		h.data = append(h.data, newValue)
	case *TCPAction:
		h.data = append(h.data, *newValue)
	}
	return fmt.Errorf("casting error")
}

func (h *TCPRequests) SetStr(data string) error {
	parts, comment := common.StringSplitWithCommentIgnoreEmpty(data, ' ')
	oldData, _ := h.Get(false)
	h.Clear()
	_, err := h.Parse(data, parts, []string{}, comment)
	if err != nil {
		h.Set(oldData)
	}
	return err
}

func (f *TCPRequests) ParseTCPRequest(request TCPAction, parts []string, comment string) error {
	err := request.Parse(parts, comment)
	if err != nil {
		return &errors.ParseError{Parser: "HTTPRequestLines", Line: ""}
	}
	f.data = append(f.data, request)
	return nil
}

func (h *TCPRequests) Parse(line string, parts, previousParts []string, comment string) (changeState string, err error) {
	if len(parts) >= 2 && parts[0] == "tcp-request" {
		var err error
		switch parts[1] {
		case "connection":
			if h.Mode == "backend" {
				return "", &errors.ParseError{Parser: "HTTPRequestLines", Line: line}
			}
			err = h.ParseTCPRequest(&actions.Connection{}, parts, comment)
		case "session":
			if h.Mode == "backend" {
				return "", &errors.ParseError{Parser: "HTTPRequestLines", Line: line}
			}
			err = h.ParseTCPRequest(&actions.Session{}, parts, comment)
		case "content":
			err = h.ParseTCPRequest(&actions.Content{}, parts, comment)
		case "inspect-delay":
			err = h.ParseTCPRequest(&actions.InspectDelay{}, parts, comment)
		default:
			return "", &errors.ParseError{Parser: "HTTPRequestLines", Line: line}
		}
		if err != nil {
			return "", err
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "HTTPRequestLines", Line: line}
}

func (h *TCPRequests) Result(AddComments bool) ([]common.ReturnResultLine, error) {
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "tpc-request " + req.String(),
			Comment: req.GetComment(),
		}
	}
	return result, nil
}
