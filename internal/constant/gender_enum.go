// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.8
// Revision: 3d844c8ecc59661ed7aa17bfd65727bc06a60ad8
// Build Date: 2023-09-18T14:55:21Z
// Built By: goreleaser

package constant

import (
	"errors"
	"fmt"
)

const (
	// GenderFemale is a Gender of type Female.
	GenderFemale Gender = iota
	// GenderMale is a Gender of type Male.
	GenderMale
	// GenderUnknown is a Gender of type Unknown.
	GenderUnknown
)

var ErrInvalidGender = errors.New("not a valid Gender")

const _GenderName = "femalemaleunknown"

var _GenderMap = map[Gender]string{
	GenderFemale:  _GenderName[0:6],
	GenderMale:    _GenderName[6:10],
	GenderUnknown: _GenderName[10:17],
}

// String implements the Stringer interface.
func (x Gender) String() string {
	if str, ok := _GenderMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Gender(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Gender) IsValid() bool {
	_, ok := _GenderMap[x]
	return ok
}

var _GenderValue = map[string]Gender{
	_GenderName[0:6]:   GenderFemale,
	_GenderName[6:10]:  GenderMale,
	_GenderName[10:17]: GenderUnknown,
}

// ParseGender attempts to convert a string to a Gender.
func ParseGender(name string) (Gender, error) {
	if x, ok := _GenderValue[name]; ok {
		return x, nil
	}
	return Gender(0), fmt.Errorf("%s is %w", name, ErrInvalidGender)
}

// MarshalText implements the text marshaller method.
func (x Gender) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Gender) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseGender(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
