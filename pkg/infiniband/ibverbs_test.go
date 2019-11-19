package main

import (
	"fmt"
	"github.com/ruanxingbaozi/gpu-affinity/pkg/infiniband/ibverbs"
)

func main() {

	ibvDevList, err := ibverbs.IbvGetDeviceList()
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(ibvDevList)
}
