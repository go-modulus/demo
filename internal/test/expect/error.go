package expect

import (
	"boilerplate/internal/framework"
	"fmt"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func MatchUserError(expected framework.ErrorIdentifier) types.GomegaMatcher {
	return &userErrorMatcher{expected: expected}
}

type userErrorMatcher struct {
	expected framework.ErrorIdentifier
}

func (m *userErrorMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, fmt.Errorf("Expected an error, got nil")
	}

	_, ok := actual.(error)
	if !ok {
		return false, fmt.Errorf("Expected an error. Got:\n%s", format.Object(actual, 1))
	}

	userError, ok := actual.(framework.UserError)
	if !ok {
		return false, fmt.Errorf(
			"MatchUserError matcher expects an framework.UserError.Got:\n%s",
			format.Object(actual, 1),
		)
	}

	return m.expected == userError.Identifier(), nil
}

func (m *userErrorMatcher) FailureMessage(actual interface{}) (message string) {
	userError := actual.(framework.UserError)
	return format.Message(userError.Identifier(), "to be equal", m.expected)
}

func (m *userErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	userError := actual.(framework.UserError)
	return format.Message(userError.Identifier(), "not to be equal", m.expected)
}

func MatchError(identifier framework.ErrorIdentifier, err error) Expectation {
	return func() (bool, string) {
		matcher := MatchUserError(identifier)
		res, matchErr := matcher.Match(err)
		if matchErr != nil {
			return false, matchErr.Error()
		}
		return res, matcher.FailureMessage(err)
	}
}

func HasError(err error) Expectation {
	return func() (bool, string) {
		matcher := gomega.HaveOccurred()
		res, err2 := matcher.Match(err)
		if err2 != nil {
			return false, err2.Error()
		}
		return res, matcher.FailureMessage(err)
	}
}

func NoError(err error) Expectation {
	return func() (bool, string) {
		matcher := gomega.BeNil()
		res, err2 := matcher.Match(err)
		if err2 != nil {
			return false, err2.Error()
		}
		return res, matcher.FailureMessage(err)
	}
}
