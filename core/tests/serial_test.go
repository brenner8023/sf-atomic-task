package tests

import (
	"fmt"
	"testing"
	"time"

	. "github.com/brenner8023/sf-atomic-task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSfAtomicTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SfAtomicTask Suite")
}

var _ = Describe("serial process A->B->C", func() {
	var (
		taskACount int
		taskBCount int
		taskCCount int
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
	deps := map[string][]string{
		"A": {},
		"B": {"A"},
		"C": {"B"},
	}
	tasks := map[string]TaskFunc{
		"A": taskA,
		"B": taskB,
		"C": taskC,
	}

	BeforeEach(func() {
		taskACount = 0
		taskBCount = 0
		taskCCount = 0
	})

	It("when run C task, task ABC should be called", func() {
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		result, err2 := run([]string{"C"})
		Expect(err2).To(BeNil())
		Expect(result).NotTo(HaveKey("A"))
		Expect(result).NotTo(HaveKey("B"))
		Expect(result["C"]).To(Equal("DataC"))
		Expect(taskACount).To(Equal(1))
		Expect(taskBCount).To(Equal(1))
		Expect(taskCCount).To(Equal(1))
	})

	It("when run B and C task, task ABC should be called", func() {
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		result, err2 := run([]string{"B", "C"})
		Expect(err2).To(BeNil())
		Expect(result).NotTo(HaveKey("A"))
		Expect(result["B"]).To(Equal("DataB"))
		Expect(result["C"]).To(Equal("DataC"))
		Expect(taskACount).To(Equal(1))
		Expect(taskBCount).To(Equal(1))
		Expect(taskCCount).To(Equal(1))
	})

	It("when C depends on A and B, C can get result of AB", func() {
		deps = map[string][]string{
			"C": {"A", "B"},
		}
		tasks := map[string]TaskFunc{
			"A": taskA,
			"B": taskB,
			"C": func(depContext map[string]any, params map[string]any) (any, error) {
				taskCCount++
				data := fmt.Sprintf("C:%v_%v", depContext["A"], depContext["B"])
				return data, nil
			},
		}
		run, err := DefineTasks(deps, &tasks)
		Expect(err).To(BeNil())
		result, err2 := run([]string{"C"})
		Expect(err2).To(BeNil())
		Expect(result["C"]).To(Equal("C:DataA_DataB"))
	})
})
