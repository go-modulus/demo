package spec

import (
	baseAssert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"testing"
)

func Equal(msg string, t *testing.T, expected any, actual any) bool {
	t.Helper()
	res := baseAssert.Equal(t, expected, actual, red+msg+reset)
	if res {
		t.Log(green + assertionPrefix + msg + reset)
	}
	return res
}

func NotEqual(msg string, t *testing.T, expected any, actual any) bool {
	t.Helper()
	res := baseAssert.NotEqual(t, expected, actual, red+msg+reset)
	if res {
		t.Log(green + assertionPrefix + msg + reset)
	}
	return res
}

func Nil(msg string, t *testing.T, val any) bool {
	t.Helper()
	res := baseAssert.Nil(t, val, red+msg+reset)
	if res {
		t.Log(green + assertionPrefix + msg + reset)
	}
	return res
}

func NotNil(msg string, t *testing.T, val any) bool {
	t.Helper()
	res := baseAssert.NotNil(t, val, red+msg+reset)
	if res {
		t.Log(green + assertionPrefix + msg + reset)
	}
	return res
}

func False(msg string, t *testing.T, val bool) bool {
	t.Helper()
	res := baseAssert.False(t, val, red+msg+reset)
	if res {
		t.Log(green + assertionPrefix + msg + reset)
	}
	return res
}

func True(msg string, t *testing.T, val bool) bool {
	t.Helper()
	//res := baseAssert.True(t, val, red+msg+reset)
	if val {
		t.Log(green + assertionPrefix + msg + reset)
		return true
	}

	t.Errorf(red + assertionPrefix + msg + reset)
	return false
}

func NoError(t *testing.T, msg string, err error) bool {
	t.Helper()
	res := baseAssert.Nil(t, err, red+msg+reset)
	res2 := baseAssert.NoError(t, err, red+msg+reset)
	if res && res2 {
		t.Log(green + assertionPrefix + msg + reset)
	}
	return res
}

func HasError(t *testing.T, msg string, err error) bool {
	t.Helper()
	res := baseAssert.Error(t, err, red+msg+reset)
	if res {
		t.Log(green + assertionPrefix + msg + reset)
	}
	return res
}

func HasParticularError(msg string, t *testing.T, err error, expectedErrorMsg string) bool {
	t.Helper()
	res := baseAssert.Error(t, err, red+msg+reset)
	if !res {
		return false
	}
	return Equal(msg, t, expectedErrorMsg, err.Error())
}

func MockedMethodIsCalled(
	msg string,
	t *testing.T,
	m *mock.Call,
	arguments ...interface{},
) bool {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Helper()
	isCalled := m.Parent.AssertCalled(t, m.Method, arguments...)
	m.Unset()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	if !isCalled {
		t.Log(red + assertionPrefix + msg + reset + "\n" + string(out))
	}
	return isCalled
}

func MockedMethodIsNotCalled(
	msg string,
	t *testing.T,
	m *mock.Call,
	arguments ...interface{},
) bool {
	t.Helper()
	isCalled := m.Parent.AssertNotCalled(t, m.Method, arguments...)
	m.Unset()
	True(msg, t, isCalled)
	return isCalled
}
