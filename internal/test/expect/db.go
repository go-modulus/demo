package expect

import (
	"boilerplate/internal/test"
)

func HasInDb(table string, conditions map[string]any) Expectation {
	return True(test.HasInDb(table, conditions))
}

func HasOneInDb(table string, conditions map[string]any) Expectation {
	return True(test.HasOneInDb(table, conditions))
}

func HasNotInDb(table string, conditions map[string]any) Expectation {
	return False(test.HasInDb(table, conditions))
}
