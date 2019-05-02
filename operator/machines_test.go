package operator

import "testing"

type TestCase struct {
	Arg1 *SparkClusterApiOperatorRequest
	Arg2  *SparkClusterApiOperatorRequest
	Solution int
}

func TestComputeNumberOfExpectedInstances(t *testing.T) {

	cases := []TestCase {

		// [1] 1:5
		{
			Arg1: &SparkClusterApiOperatorRequest{
				CPUCount: 1,
				MemoryMBs: 1000,
				ContainerCount: 1,
			},
			Arg2: &SparkClusterApiOperatorRequest{
				CPUCount: 5,
				MemoryMBs: 5000,
				ContainerCount: 5,
			},
			Solution: 5,
		},

		// [2] 5:5
		{
			Arg1: &SparkClusterApiOperatorRequest{
				CPUCount: 5,
				MemoryMBs: 5000,
				ContainerCount: 5,
			},
			Arg2: &SparkClusterApiOperatorRequest{
				CPUCount: 5,
				MemoryMBs: 5000,
				ContainerCount: 5,
			},
			Solution: 2,
		},


		// [3] 1:5
		{
			Arg1: &SparkClusterApiOperatorRequest{
				CPUCount: 1,
				MemoryMBs: 1000,
				ContainerCount: 1,
			},
			Arg2: &SparkClusterApiOperatorRequest{
				CPUCount: 1,
				MemoryMBs: 5000,
				ContainerCount: 9,
			},
			Solution: 5,
		},

		// [3] 1:22
		{
			Arg1: &SparkClusterApiOperatorRequest{
				CPUCount: 1,
				MemoryMBs: 1000,
				ContainerCount: 1,
			},
			Arg2: &SparkClusterApiOperatorRequest{
				CPUCount: 1,
				MemoryMBs: 22000,
				ContainerCount: 9,
			},
			Solution: 22,
		},
	}

	for _, test := range cases {
		x := ComputeNumberOfExpectedInstances(test.Arg1, test.Arg2)
		if x != test.Solution {
			t.Errorf("Error asserting { CPU: %d, MemoryMB: %d, Containers: %d}  [%d] == solution [%d]", test.Arg2.CPUCount, test.Arg2.MemoryMBs, test.Arg2.ContainerCount, x, test.Solution)
		}
	}

}
