package gh

import "runtime"

const ghRepo = "Azure/communication-ui-library"

func maxParallelFetch() int {
	n := runtime.NumCPU()
	if n > 8 {
		return n
	}
	return 8
}
