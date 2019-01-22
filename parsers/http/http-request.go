package http

import (
	"strings"

	"github.com/haproxytech/config-parser/common"
	"github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/parsers/http/actions"
)

type HTTPRequests struct {
	HTTPRequests []HTTPAction
}

func (h *HTTPRequests) Init() {
	h.HTTPRequests = []HTTPAction{}
}

func (h *HTTPRequests) GetParserName() string {
	return "http-request"
}

func (f *HTTPRequests) ParseHTTPRequest(request HTTPAction, parts []string, comment string) error {
	err := request.Parse(parts, comment)
	if err != nil {
		return &errors.ParseError{Parser: "HTTPRequestLines", Line: ""}
	}
	f.HTTPRequests = append(f.HTTPRequests, request)
	return nil
}

func (h *HTTPRequests) Parse(line string, parts, previousParts []string, comment string) (changeState string, err error) {
	if len(parts) >= 2 && parts[0] == "http-request" {
		var err error
		switch parts[1] {
		case "add-header":
			err = h.ParseHTTPRequest(&actions.AddHeader{}, parts, comment)
		case "del-header":
			err = h.ParseHTTPRequest(&actions.DelHeader{}, parts, comment)
		case "set-header":
			err = h.ParseHTTPRequest(&actions.SetHeader{}, parts, comment)
		case "set-var":
			err = h.ParseHTTPRequest(&actions.SetVar{}, parts, comment)
		case "allow":
			err = h.ParseHTTPRequest(&actions.Allow{}, parts, comment)
		case "deny":
			err = h.ParseHTTPRequest(&actions.Deny{}, parts, comment)
		case "redirect":
			err = h.ParseHTTPRequest(&actions.Redirect{}, parts, comment)
		case "auth":
			err = h.ParseHTTPRequest(&actions.Auth{}, parts, comment)
		default:
			if strings.HasPrefix(parts[1], "add-acl(") {
				err = h.ParseHTTPRequest(&actions.AddAcl{}, parts, comment)
			} else if strings.HasPrefix(parts[1], "del-acl(") {
				err = h.ParseHTTPRequest(&actions.DelAcl{}, parts, comment)
			} else if strings.HasPrefix(parts[1], "set-var(") {
				err = h.ParseHTTPRequest(&actions.SetVar{}, parts, comment)
			} else {
				return "", &errors.ParseError{Parser: "HTTPRequestLines", Line: line}
			}
		}
		if err != nil {
			return "", err
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "HTTPRequestLines", Line: line}
}

func (h *HTTPRequests) Valid() bool {
	if len(h.HTTPRequests) > 0 {
		return true
	}
	return false
}

func (h *HTTPRequests) Result(AddComments bool) []common.ReturnResultLine {
	result := make([]common.ReturnResultLine, len(h.HTTPRequests))
	for index, req := range h.HTTPRequests {
		result[index] = common.ReturnResultLine{
			Data:    "http-request " + req.String(),
			Comment: req.GetComment(),
		}
	}
	return result
}
