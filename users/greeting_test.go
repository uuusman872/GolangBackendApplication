package users

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Greeting(prefix, name string) {
	fmt.Printf("%s %s", prefix, name)
}

func TestGreeting(t *testing.T) {
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w
	defer func() {
		os.Stdout = originalStdout
	}()
	Greeting("Mr", "Joseph")
	w.Close()
	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	assert.NoError(t, err)

	want := "Mr. Joseph"
	got := buf.String()
	fmt.Printf("")
	assert.Equal(t, want, got)
}
