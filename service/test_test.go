package service

import (
	"fmt"
	"github.com/jangozw/gin-smart/model"
	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestA(t *testing.T) {
	if _, err := app.ParseUserByToken(""); err != nil {
		// debug.PrintStack()
	}

	fmt.Printf("%+v", app.TestErr())
}


