package main

import (
	"fmt"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/mitchellh/go-ps"
	"log"
	"strings"
)

func main() {
	if err := nvml.Init(); err != nil {
		log.Printf("Failed to initialize NVML: %s.", err)
	}
	defer func() { log.Println("Shutdown of NVML returned:", nvml.Shutdown()) }()
	//getGpuProcessTree()
	topo := getGpuTopo()
	for k, v := range topo {
		log.Printf("topo %v : %v", k, v)
	}
}

func getGpuProcessTree() {
	processMap := make(map[string][]nvml.ProcessInfo)
	// get all device
	devices := getAllDevices()
	for _, device := range devices {
		// get all running gpu process
		process := getAllRunningProcess(device)
		for _, p := range process {
			// get container-shim process map[ppid]process
			dockerProcess, err := getDockerPpid(int(p.PID))
			if err != nil {
				log.Printf("find process error %v", err)
			}
			key := fmt.Sprintf("<%v/%v>", dockerProcess.PPid(), dockerProcess.Pid())
			if plist, found := processMap[key]; found {
				plist = append(plist, p)
			} else {
				processMap[key] = []nvml.ProcessInfo{p}
			}
		}
	}
	for k, v := range processMap {
		log.Printf("ppid %v , process %v", k, v)
	}
}

func mergeContainer(pid int, ppid int) {
	// get all container
	// filter container pid == pid
	// check pid's ppid == ppid
	// merge pid -> containerName

}
func getProcessInfo(pid int) {
	p, err := ps.FindProcess(pid)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	log.Printf("process %v", p)
}
func getDockerPpid(pid int) (ps.Process, error) {
	process, _ := ps.FindProcess(pid)
	for {
		pprocess, err := ps.FindProcess(process.PPid())
		if err != nil {
			log.Printf("error %v", err)
			return nil, err
		}
		if strings.Contains(pprocess.Executable(), "containerd-shim") {
			log.Printf("process pid %v, containerd-shim pid %v", process.Pid(), pprocess.Pid())
			return process, nil
		}
		process = pprocess
	}
}
func checkErr(err error) {
	log.Printf("Error %v", err)
}
func getAllRunningProcess(dev *nvml.Device) []nvml.ProcessInfo {
	allRunningP, err := dev.GetAllRunningProcesses()
	checkErr(err)
	log.Printf("all %v", allRunningP)
	return allRunningP
}
//func getP2PLink(dev1, dev2 *nvml.Device) nvml.P2PLinkType {
//	log.Printf("get p2p link")
//	link, err := nvml.GetP2PLink(dev1, dev2)
//	if err != nil {
//		log.Printf("Error %v\n", err)
//	}
//	log.Printf("dev %v dev %v ,p2p link: %v", dev1.UUID, dev2.UUID, link)
//	return link
//}

//// 获取gpu拓扑结构
//func getGpuTopo() map[string][]int {
//	// gpu0 0,1,2,2,6,6,6,6
//	topo := make(map[string][]int)
//	devices := getAllDevices()
//	for idx, ldev := range devices {
//		for _, rdev := range devices {
//			link := getP2PLink(ldev, rdev)
//			distance := 0
//			switch link {
//			case nvml.P2PLinkCrossCPU:
//				distance = 6
//			case nvml.P2PLinkSameCPU:
//				distance = 5
//			case nvml.P2PLinkHostBridge:
//				distance = 4
//			case nvml.P2PLinkMultiSwitch:
//				distance = 3
//			case nvml.P2PLinkSingleSwitch:
//				distance = 2
//			case nvml.P2PLinkSameBoard:
//				distance = 1
//			case nvml.P2PLinkUnknown:
//				distance = 0
//			}
//			key := fmt.Sprintf("%v/%v", ldev.UUID, idx)
//			if links, found := topo[key]; found {
//				links = append(links, distance)
//				topo[key] = links
//			} else {
//				links = []int{distance}
//				topo[key] = links
//			}
//		}
//	}
//	return topo
//}
//
//func getAllDevices() []*nvml.Device {
//	log.Printf("get all device")
//	devs := []*nvml.Device{}
//	n, err := nvml.GetDeviceCount()
//	if err != nil {
//		log.Printf("Error %v\n", err)
//	}
//	for i := uint(0); i < n; i++ {
//		d, err := nvml.NewDeviceLite(i)
//		if err != nil {
//			log.Printf("Error %v\n", err)
//		}
//		log.Printf("device %v", d)
//		devs = append(devs, d)
//	}
//	return devs
//}

