package model

import (
	"testing"

	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {
	user := &SampleUser{}
	err := app.Db.Where("id=?", 1).First(&user).Error
	assert.Nil(t, err)
}
