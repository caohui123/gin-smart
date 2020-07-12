package service

import (
	"testing"

	"github.com/jangozw/gin-smart/pkg/app"
)

func TestA(t *testing.T) {
	if _, err := app.ParseUserByToken(""); err != nil {
		// debug.PrintStack()
	}
}
