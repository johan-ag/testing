package examples

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type anything struct {
	id       uint
	name     string
	deadline time.Duration
}

func TestAddition(t *testing.T) {
	want := &anything{
		id:       1,
		name:     "bari",
		deadline: time.Second * 2,
	}

	got := &anything{
		id:       1,
		name:     "bari",
		deadline: time.Second * 2,
	}

	require.Equal(t, want, got)
}
