package tests

import (
	"fmt"
	"time"

	. "github.com/brenner8023/sf-atomic-task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parallel process", func() {
	var (
		taskACount int
		taskBCount int
		taskCCount int
		taskDCount int
	)

	taskA := func(depContext map[string]any, params map[string]any) (any, error) {
		taskACount++
		return "DataA", nil
	}
	taskB := func(depContext map[string]any, params map[string]any) (any, error) {
		taskBCount++
		time.Sleep(100 * time.Microsecond)
		return "DataB", nil
	}
	taskC := func(depContext map[string]any, params map[string]any) (any, error) {
		taskCCount++
		return "DataC", nil
	}
	taskD := func(depContext map[string]any, params map[string]any) (any, error) {
		taskDCount++
		return "DataD", nil
	}
	deps := map[string][]string{
		"B": {"A"},
		"C": {"A"},
		"D": {"A"},
	}
	tasks := map[string]TaskFunc{
		"A": taskA,
		"B": taskB,
		"C": taskC,
		"D": taskD,
	}

	BeforeEach(func() {
		taskACount = 0
		taskBCount = 0
		taskCCount = 0
		taskDCount = 0
	})

	It("run task A->B|C|D", func() {
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		result, err2 := run([]string{"B", "C", "D"})
		Expect(err2).To(BeNil())
		Expect(result).NotTo(HaveKey("A"))
		Expect(result["B"]).To(Equal("DataB"))
		Expect(result["C"]).To(Equal("DataC"))
		Expect(result["D"]).To(Equal("DataD"))
		Expect(taskACount).To(Equal(1))
		Expect(taskBCount).To(Equal(1))
		Expect(taskCCount).To(Equal(1))
		Expect(taskDCount).To(Equal(1))
	})

	It("with no deps, run task B|C|D", func() {
		run, err := DefineTasks(map[string][]string{}, &tasks)
		Expect(err).To(BeNil())
		result, err2 := run([]string{"B", "C", "D"})
		Expect(err2).To(BeNil())
		Expect(result["B"]).To(Equal("DataB"))
		Expect(result["C"]).To(Equal("DataC"))
		Expect(result["D"]).To(Equal("DataD"))
		Expect(taskBCount).To(Equal(1))
		Expect(taskCCount).To(Equal(1))
		Expect(taskDCount).To(Equal(1))
	})

	It("run task A->B|C->D, with no deps", func() {
		deps = map[string][]string{
			"B": {"A"},
			"C": {"A"},
			"D": {"B", "C"},
		}
		taskA = func(depContext map[string]any, params map[string]any) (any, error) {
			taskACount++
			return "DataA", nil
		}
		taskB = func(depContext map[string]any, params map[string]any) (any, error) {
			taskBCount++
			time.Sleep(100 * time.Microsecond)
			return fmt.Sprintf("%v->DataB", depContext["A"]), nil
		}
		taskC = func(depContext map[string]any, params map[string]any) (any, error) {
			taskCCount++
			return fmt.Sprintf("%v->DataC", depContext["A"]), nil
		}
		taskD = func(depContext map[string]any, params map[string]any) (any, error) {
			taskDCount++
			return fmt.Sprintf("DataD:%v_%v", depContext["B"], depContext["C"]), nil
		}
		tasks = map[string]TaskFunc{
			"A": taskA,
			"B": taskB,
			"C": taskC,
			"D": taskD,
		}
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		result, err2 := run([]string{"D"})
		Expect(err2).To(BeNil())
		Expect(result).NotTo(HaveKey("A"))
		Expect(result).NotTo(HaveKey("B"))
		Expect(result).NotTo(HaveKey("C"))
		Expect(result["D"]).To(Equal("DataD:DataA->DataB_DataA->DataC"))
		Expect(taskACount).To(Equal(1))
		Expect(taskBCount).To(Equal(1))
		Expect(taskCCount).To(Equal(1))
		Expect(taskDCount).To(Equal(1))
	})
})
