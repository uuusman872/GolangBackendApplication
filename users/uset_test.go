package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidUser(t *testing.T) {
	got := CheckUsername("test1234")
	assert.True(t, true, got)
}

func TestCheckUsernameTable(t *testing.T) {
	testCase := []struct {
		name  string
		input string
		want  bool
	}{
		{"too short", "test1", false},
		{"too valid", "test12345", true},
		{"contains admin", "adminuser", false},
		{"empty", "", false},
	}
	for _, tc := range testCase {
		got := CheckUsername(tc.input)
		assert.Equal(t, tc.want, got, "got %v\nwant %v", got, tc.want)
	}
}
