package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type formerrors map[string][]string

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (e formerrors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e formerrors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

type Form struct {
	url.Values
	Errors formerrors
}

func NewForm(form url.Values) *Form {
	return &Form{
		Values: form,
		Errors: formerrors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) *Form {
	for _, field := range fields {
		value := f.Get(field)
		fmt.Printf("The key is %s and its value is %s", field, value)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
	return f
}

func (f *Form) MinLength(field string, n int) *Form {
	values := f.Get(field)
	if values == "" {
		return f
	}
	if utf8.RuneCountInString(values) < n {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum of %d characters are required)", n))
	}
	return f
}

func (f *Form) MaxLength(field string, n int) *Form {
	values := f.Get(field)
	if values == "" {
		return f
	}
	if utf8.RuneCountInString(values) > n {
		f.Errors.Add(field, fmt.Sprintf("This field is too Long (maximum of %d characters are required)", n))
	}
	return f
}

func (f *Form) IsEmail(field string, pattern *regexp.Regexp) *Form {
	values := strings.TrimSpace(f.Get(field))
	if values == "" {
		return f
	}
	if !pattern.MatchString(values) {
		f.Errors.Add(field, "This field is invalid")
	}
	return f
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
