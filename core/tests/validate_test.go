package tests

import (
	"fmt"
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

		deps = map[string][]string{
			"A": {"A"},
		}
		tasks = map[string]TaskFunc{
			"A": taskA,
		}
		_, err = DefineTasks(deps, &tasks)
		Expect(err).To(MatchError("sf-atomic-task: validateDeps - deps[A] has a circular dependency"))

		deps = map[string][]string{
			"A": {"C"},
		}
		tasks = map[string]TaskFunc{
			"A": taskA,
		}
		_, err = DefineTasks(deps, &tasks)
		Expect(err).To(MatchError("sf-atomic-task: validateDeps - deps[A]C is not defined in tasks"))
	})

	It("function validateRunningTasks should work", func() {
		deps := map[string][]string{}
		tasks := map[string]TaskFunc{
			"A": taskA,
			"B": taskB,
		}
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		_, err = run([]string{"C"})
		Expect(err).To(MatchError("sf-atomic-task: validateRunningTasks - C is not defined in tasks"))
	})

	It("debug mode should work", func() {
		deps := map[string][]string{}
		tasks := map[string]TaskFunc{
			"A": taskA,
			"B": taskB,
		}
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		_, err = run([]string{"A"}, map[string]any{"debug": true})
		Expect(err).To(BeNil())
	})

	It("when task throw error, should return error", func() {
		deps := map[string][]string{}
		tasks := map[string]TaskFunc{
			"A": func(depContext map[string]any, params map[string]any) (any, error) {
				return nil, fmt.Errorf("errorA")
			},
			"B": func(depContext map[string]any, params map[string]any) (any, error) {
				return "DataB", nil
			},
		}
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		_, err = run([]string{"A", "B"})
		Expect(err).To(MatchError("errorA"))
	})

	It("validate params", func() {
		deps := map[string][]string{}
		tasks := map[string]TaskFunc{
			"A": taskA,
			"B": taskB,
		}
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		_, err = run([]string{"B"}, map[string]any{"a1": 1}, map[string]any{"a2": 2})
		Expect(err).To(MatchError("sf-atomic-task: run - too many parameters, expected at most 1"))
	})
})
