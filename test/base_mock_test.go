package test

import (
	"github.com/poblish/promenade/api"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMockApiSpec(t *testing.T) {
	test(t, &BaseMock{})
}

func test(t *testing.T, m api.PrometheusMetrics) {
	require.NotNil(t, m)
}
