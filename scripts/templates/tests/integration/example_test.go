package integration

import (
	"testing"

	"github.com/rancher/charts/tests/utils"
	"github.com/stretchr/testify/require"
)

func TestPass(t *testing.T) {
	utils.HelloWorld()
	require.Equal(t, "hello-world", "hello-world")
}
