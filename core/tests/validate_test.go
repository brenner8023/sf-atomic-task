package tests

import (
	"time"

	. "github.com/brenner8023/sf-atomic-task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("validate process", func() {
	taskA := func(depContext map[string]any, params map[string]any) (any, error) {
		return "DataA", nil
	}
	taskB := func(depContext map[string]any, params map[string]any) (any, error) {
		time.Sleep(100 * time.Microsecond)
		return "DataB", nil
	}
	// taskC := func(depContext map[string]any, params map[string]any) (any, error) {
	// 	return "DataC", nil
	// }

	BeforeEach(func() {
		
	})

	It("function ValidateDeps should work", func() {
		deps := map[string][]string{
			"A": {},
			"C": {"A"},
		}
		tasks := map[string]TaskFunc{
			"A": taskA,
			"B": taskB,
		}
		_, err := DefineTasks(deps, &tasks)
		Expect(err).To(MatchError("sf-atomic-task: validateDeps - deps[C] is not defined in tasks"))
	})
})
