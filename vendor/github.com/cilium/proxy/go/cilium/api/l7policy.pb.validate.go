// Code generated by protoc-gen-validate
// source: cilium/api/l7policy.proto
// DO NOT EDIT!!!

package cilium

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// Validate checks the field values on L7Policy with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *L7Policy) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for AccessLogPath

	// no validation rules for PolicyName

	// no validation rules for Denied_403Body

	if v, ok := interface{}(m.GetIsIngress()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return L7PolicyValidationError{
				Field:  "IsIngress",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// L7PolicyValidationError is the validation error returned by
// L7Policy.Validate if the designated constraints aren't met.
type L7PolicyValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e L7PolicyValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sL7Policy.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = L7PolicyValidationError{}
