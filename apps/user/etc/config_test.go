package etc_test

import (
	"testing"

	"moj/apps/user/etc"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	conf := etc.NewAppConfig()
	require.Equal(t, 8080, conf.AppPort)
}
