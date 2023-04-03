package expect

import (
	"github.com/onsi/gomega"
)

type Expectation = func() (bool, string)

func Equal(expected any, actual any) Expectation {
	return func() (bool, string) {
		matcher := gomega.Equal(expected)
		res, err := matcher.Match(actual)
		if err != nil {
			return false, err.Error()
		}
		return res, matcher.FailureMessage(actual)
	}
}

func NotEqual(expected any, actual any) Expectation {
	return func() (bool, string) {
		matcher := gomega.Not(gomega.Equal(expected))
		res, err := matcher.Match(actual)
		if err != nil {
			return false, err.Error()
		}
		return res, matcher.FailureMessage(actual)
	}
}

func Nil(val any) Expectation {
	return func() (bool, string) {
		matcher := gomega.BeNil()
		res, err := matcher.Match(val)
		if err != nil {
			return false, err.Error()
		}
		return res, matcher.FailureMessage(val)
	}
}

func NotNil(val any) Expectation {
	return func() (bool, string) {
		matcher := gomega.Not(gomega.BeNil())
		res, err := matcher.Match(val)
		if err != nil {
			return false, err.Error()
		}
		return res, matcher.FailureMessage(val)
	}
}

func True(actual any) Expectation {
	return func() (bool, string) {
		matcher := gomega.BeTrue()
		res, err := matcher.Match(actual)
		if err != nil {
			return false, err.Error()
		}
		return res, matcher.FailureMessage(actual)
	}
}

func False(actual any) Expectation {
	return func() (bool, string) {
		matcher := gomega.BeFalse()
		res, err := matcher.Match(actual)
		if err != nil {
			return false, err.Error()
		}
		return res, matcher.FailureMessage(actual)
	}
}
