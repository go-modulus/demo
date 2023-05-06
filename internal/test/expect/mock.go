package expect

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

func MockedMethodIsCalled(
	t *testing.T,
	m *mock.Call,
	arguments ...interface{},
) Expectation {
	isCalled := m.Parent.AssertCalled(t, m.Method, arguments...)
	m.Unset()
	return True(isCalled)
}

func MockedMethodIsNotCalled(
	t *testing.T,
	m *mock.Call,
	arguments ...interface{},
) Expectation {
	isCalled := m.Parent.AssertNotCalled(t, m.Method, arguments...)
	m.Unset()
	return True(isCalled)
}
