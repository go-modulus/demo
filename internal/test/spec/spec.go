package spec

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/test/expect"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/nsf/jsondiff"
	"testing"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var cyan = "\033[36m"
var blue = "\033[34m"

var assertionPrefix = "    "

func When(t *testing.T, msg string) {
	t.Helper()
	t.Log(cyan + "When " + msg + reset)
}

func Context(t *testing.T, msg string) {
	t.Helper()
	t.Log(cyan + msg + reset)
}

func Then(t *testing.T, message string, expectations ...expect.Expectation) bool {
	t.Helper()
	for _, expectation := range expectations {
		isPassed, errMsg := expectation()
		if !isPassed {
			t.Errorf(red + assertionPrefix + message + reset)
			t.Log(assertionPrefix + assertionPrefix + errMsg)
			return false
		}
	}

	t.Log(green + assertionPrefix + message + reset)

	return true
}

func ThenJsonContains(
	t *testing.T,
	message string,
	expected map[string]interface{},
	actualJson []byte,
) bool {
	t.Helper()
	expectedJson, _ := json.Marshal(expected)
	opt := jsondiff.DefaultConsoleOptions()
	opt.SkipMatches = true
	opt.SkippedObjectProperty = nil
	opt.SkippedArrayElement = nil

	diff, str := jsondiff.Compare(
		actualJson,
		expectedJson,
		&opt,
	)
	if jsondiff.SupersetMatch != diff {
		t.Errorf(red + assertionPrefix + message + reset + "\n" + str)
		return false
	}
	return true
}

func HasCommonError(t *testing.T, msg string, err error, expected *framework.CommonError) bool {
	t.Helper()
	cerr, ok := err.(*framework.CommonError)
	if !ok {
		t.Errorf(red + assertionPrefix + msg + reset)
		t.Log(
			assertionPrefix+assertionPrefix+
				"Expected error to be of type CommonError, but got: %v", err,
		)
		return false
	}
	expectedIdent := expected.Identifier()
	if cerr.Identifier() != expectedIdent {
		t.Errorf(red + assertionPrefix + msg + reset)
		t.Log(
			assertionPrefix+assertionPrefix+
				"Expected error identifier to be %v, but got: %v", expectedIdent, cerr.Identifier(),
		)
		return false
	}
	return true
}

func Auth(t *testing.T, ctx context.Context) {
	t.Helper()
	testUserId, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	id := framework.GetCurrentUserId(ctx)
	msg := ""
	if id == nil {
		msg = "I am not authenticated"
	} else if *id == testUserId.String() {
		msg = "I am authenticated as default test user"
	} else {
		msg = "I am authenticated as user " + *id
	}

	t.Log(cyan + msg + reset)
}

func Given(t *testing.T, descriptions ...string) {
	t.Helper()
	//msg := strings.Join(descriptions, "\n"+assertionPrefix)
	//if msg != "" {
	//	msg = "\n" + assertionPrefix + msg
	//}

	t.Log(cyan + "Given:" + reset)
	for _, desc := range descriptions {
		t.Log(blue + assertionPrefix + desc + reset)
	}
}
