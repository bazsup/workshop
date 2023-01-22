//go:build unit

package cloudpocket_test

import (
	"testing"

	"github.com/kkgo-software-engineering/workshop/cloudpocket"
	"github.com/stretchr/testify/assert"
)

func TestValidAmount(t *testing.T) {
	transfer := cloudpocket.Transfer{Amount: 0.1}
	assert.True(t, transfer.IsValidAmount())
}


func TestInvalidAmount(t *testing.T) {
	transfer := cloudpocket.Transfer{Amount: 0.001}
	assert.False(t, transfer.IsValidAmount())
}
