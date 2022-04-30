package formatter

import (
	"fmt"
	"testing"
	"time"

	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/stretchr/testify/require"
)

func TestTextFormat(t *testing.T) {
	level := contract.PanicLevel
	msg := "test for panic"
	bs, err := TextFormatter(level, time.Now(), msg, nil)
	require.NoError(t, err)
	require.NotEmpty(t, bs)
	fmt.Println(string(bs))
}
