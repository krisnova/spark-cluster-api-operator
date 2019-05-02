package operator

import "github.com/kris-nova/logger"

type SparkClusterApiOperatorRequest struct {
	CPUCount int
	MemoryMBs int
	ContainerCount int
}

// ComputeNumberOfExpectedInstances will returned a signed integer based on a solution to the input
func ComputeNumberOfExpectedInstances (serverSize *SparkClusterApiOperatorRequest, mutation *SparkClusterApiOperatorRequest) int {

	// Validate we have values everywhere or else bust with -1
	if mutation.ContainerCount == 0 {
		return -1
	}
	if mutation.MemoryMBs == 0 {
		return -1
	}
	if mutation.CPUCount == 0 {
		return -1
	}

	// For the world's simplest impl for today
	// Let's make sure we have enough resources based on a given server size for this requested mutation
	// Let n denote the number of servers of size serverSize we need
	n := 0

	// Check CPU First
	cpudiv := mutation.CPUCount / serverSize.CPUCount
	cpumod := mutation.CPUCount % serverSize.CPUCount
	//fmt.Println(1)
	n = cpudiv
	if cpumod != 0 {
		//fmt.Println(2)
		n++
	}

	//Memory next
	memdiv := mutation.MemoryMBs / serverSize.MemoryMBs
	memmod := mutation.MemoryMBs % serverSize.MemoryMBs
	if memdiv > n {
		//fmt.Println(3)
		n = memdiv
		if memmod != 0 {
			//fmt.Println(4)
			n++
		}
	}
	if memdiv == memmod && memmod != 0 {
		//fmt.Println(5)
		n++
	}

	if n > mutation.ContainerCount {
		logger.Warning("Number of machines greater than number of containers")
	}


	// Next let's try to get close to 2 containers per server as a minimum
	halfrequesteddiv := mutation.ContainerCount / 2
	//halfrequestedmod := mutation.ContainerCount % 2
	if n < halfrequesteddiv {
		//fmt.Println(6)
		n = halfrequesteddiv
	}

	return n
}