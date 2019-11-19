package main

import (
	"fmt"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"log"
	"strconv"
)

func main() {
	if err := nvml.Init(); err != nil {
		log.Printf("Failed to initialize NVML: %s.", err)
	}
	defer func() { log.Println("Shutdown of NVML returned:", nvml.Shutdown()) }()
	topo := getGpuTopo()
	distance := ""
	gpus := ""
	for k, v := range topo {
		gpus += k + ","
		for _, v1 := range v {
			distance += strconv.Itoa(v1) + " "
		}
		distance += ","
	}
	log.Printf("%v", gpus)
	log.Printf("%v", distance)
}
func getNumaNode() {
	devices := getAllDevices()
	for _, dev := range devices {
		CPUAffinity := dev.CPUAffinity
		fmt.Println(CPUAffinity)
	}
}

// 获取gpu拓扑结构
func getGpuTopo() map[string][]int {
	// gpu0 0,1,2,2,6,6,6,6
	topo := make(map[string][]int)
	devices := getAllDevices()
	for _, ldev := range devices {
		for _, rdev := range devices {
			link := getP2PLink(ldev, rdev)
			distance := 0
			switch link {
			case nvml.P2PLinkCrossCPU:
				distance = 6
			case nvml.P2PLinkSameCPU:
				distance = 5
			case nvml.P2PLinkHostBridge:
				distance = 4
			case nvml.P2PLinkMultiSwitch:
				distance = 3
			case nvml.P2PLinkSingleSwitch:
				distance = 2
			case nvml.P2PLinkSameBoard:
				distance = 1
			case nvml.P2PLinkUnknown:
				distance = 0
			}
			key := fmt.Sprintf("%v", ldev.UUID)
			if links, found := topo[key]; found {
				links = append(links, distance)
				topo[key] = links
			} else {
				links = []int{distance}
				topo[key] = links
			}
		}
	}
	return topo
}
func getP2PLink(dev1, dev2 *nvml.Device) nvml.P2PLinkType {
	log.Printf("get p2p link")
	link, err := nvml.GetP2PLink(dev1, dev2)
	if err != nil {
		log.Printf("Error %v\n", err)
	}
	log.Printf("dev %v dev %v ,p2p link: %v", dev1.UUID, dev2.UUID, link)
	return link
}
func getAllDevices() []*nvml.Device {
	log.Printf("get all device")
	devs := []*nvml.Device{}
	n, err := nvml.GetDeviceCount()
	if err != nil {
		log.Printf("Error %v\n", err)
	}
	for i := uint(0); i < n; i++ {
		d, err := nvml.NewDeviceLite(i)
		if err != nil {
			log.Printf("Error %v\n", err)
		}
		log.Printf("device %v", d)
		devs = append(devs, d)
	}
	return devs
}

/**
2019/09/06 17:30:30 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 ,p2p link: N/A
2019/09/06 17:30:30 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 dev GPU-cb4c85c1-04cb-7a37-10e5-71f9766e7871 ,p2p link: Single PCI switch
2019/09/06 17:30:30 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 dev GPU-30d810f3-2170-3faf-f597-dd32467e50dc ,p2p link: Same CPU socket
2019/09/06 17:30:30 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 dev GPU-edea8dd7-fb72-c7f8-09fc-19757672e417 ,p2p link: Same CPU socket
2019/09/06 17:30:30 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 dev GPU-5ee3f1b6-dd9c-3d32-02b6-91a8daf2dd0b ,p2p link: Cross CPU socket
2019/09/06 17:30:30 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 dev GPU-0243ee05-6f49-fbae-87d2-f8b6dccab18f ,p2p link: Cross CPU socket
2019/09/06 17:30:30 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 dev GPU-39ea8922-a1da-4e77-8b3f-08dbd00fa53a ,p2p link: Cross CPU socket
2019/09/06 17:30:30 dev GPU-e021c6fc-b7da-289b-b504-24d355822ab3 dev GPU-2ce8723e-3023-bfca-2b54-e6c6f193b07c ,p2p link: Cross CPU socket
leinao@XP001:~$ nvidia-smi topo -m
        GPU0    GPU1    GPU2    GPU3    GPU4    GPU5    GPU6    GPU7    mlx5_0    CPU Affinity
GPU0     X     PIX    NODE    NODE    SYS    SYS    SYS    SYS    SYS    0-11
GPU1    PIX     X     NODE    NODE    SYS    SYS    SYS    SYS    SYS    0-11
GPU2    NODE    NODE   X     PIX    SYS    SYS    SYS    SYS    SYS    0-11
GPU3    NODE    NODE  PIX     X     SYS    SYS    SYS    SYS    SYS    0-11
GPU4    SYS    SYS    SYS    SYS     X     PIX    NODE    NODE    NODE    12-23
GPU5    SYS    SYS    SYS    SYS    PIX     X     NODE    NODE    NODE    12-23
GPU6    SYS    SYS    SYS    SYS    NODE    NODE     X     PIX    NODE    12-23
GPU7    SYS    SYS    SYS    SYS    NODE    NODE    PIX     X     NODE    12-23
mlx5_0  SYS    SYS    SYS    SYS    NODE    NODE    NODE    NODE     X

Legend:

  X    = Self
  SYS  = Connection traversing PCIe as well as the SMP interconnect between NUMA nodes (e.g., QPI/UPI)
  NODE = Connection traversing PCIe as well as the interconnect between PCIe Host Bridges within a NUMA node
  PHB  = Connection traversing PCIe as well as a PCIe Host Bridge (typically the CPU)
  PXB  = Connection traversing multiple PCIe switches (without traversing the PCIe Host Bridge)
  PIX  = Connection traversing a single PCIe switch
  NV#  = Connection traversing a bonded set of # NVLinks
*/
/**


# select 3 gpu
0 2 5 5 6 6 6 6
2 0 5 5 6 6 6 6
5 5 0 2 6 6 6 6
5 5 2 0 6 6 6 6

6 6 6 6 0 2 5 5
6 6 6 6 2 0 5 5
6 6 6 6 5 5 0 2
6 6 6 6 5 5 2 0

# select 3 gpu
   G0 G4 G5
G0 0  6  6
G4 6  0  2
G5 6  2  0
distance=6

   G4 G5 G6
G4 0  2  5
G5 2  0  2
G6 5  2  0
distance=5

   G5 G6 G7
G5 0  5  5
G6 5  0  2
G7 5  2  0
distance=5

# allocate
   G0 G4
G0 0  6
G4 6  0
distance=6
   G4 G5
G4 0  2
G5 2  0
distance=2
   G5 G6
G5 0  5
G6 5  0
distance=5
   G6 G7
G6 0  2
G7 2  0
distance=2

# backfit
   G0 G6 G7
G0 0  6  6
G6 6  0  2
G7 6  2  0
{(G0,G6),(G0,G7)} distance=6
=> (G0,G6,G7)
{(G6,G7)} distance=2
=> (G6,G7)
(G0,G6,G7) - (G6,G7) = G0

{(G0,G6),(G0,G7)} 6*2
{(G6,G7)} 2*1


nvidia.com/gpu-topo: "0 2 5 5 6 6 6 6,2 0 5 5 6 6 6 6,5 5 0 2 6 6 6 6,5 5 2 0 6 6 6 6,6 6 6 6 0 2 5 5,6 6 6 6 2 0 5 5,6 6 6 6 5 5 0 2,6 6 6 6 5 5 2 0"
nvidia.com/gpu-capacity: "GPU-e021c6fc-b7da-289b-b504-24d355822ab3,GPU-cb4c85c1-04cb-7a37-10e5-71f9766e7871,GPU-30d810f3-2170-3faf-f597-dd32467e50dc,GPU-edea8dd7-fb72-c7f8-09fc-19757672e417,GPU-5ee3f1b6-dd9c-3d32-02b6-91a8daf2dd0b,GPU-0243ee05-6f49-fbae-87d2-f8b6dccab18f,GPU-39ea8922-a1da-4e77-8b3f-08dbd00fa53a,GPU-2ce8723e-3023-bfca-2b54-e6c6f193b07c"
nvidia.com/gpu-allocated: "GPU-e021c6fc-b7da-289b-b504-24d355822ab3,GPU-cb4c85c1-04cb-7a37-10e5-71f9766e7871,GPU-30d810f3-2170-3faf-f597-dd32467e50dc"
nvidia.com/gpu-isolation
# Allocated
availabel: 3 4 5 6 7
1、 3 4 =6
2、 4 5 =2         min
5、 5 6 =5
6、 6 7 =2         min
gpunum >=2 取距离最远的
# backfit
1、 (3 4)(3 5)(3 6)(3 7) =6    24
2、 (4 5) =2 (4 6)(4 7)=5      12
3、 (5 6)(5 7) =5              10
4、 (6 7)=2                    2
backfit暂时只针对gpunum=1
*/
/**
root@nvidia-device-plugin-daemonset-jsx8p:/# nvidia-smi topo -m
	GPU0	GPU1	GPU2	GPU3	GPU4	GPU5	GPU6	GPU7	mlx5_0	mlx5_1	mlx5_2	mlx5_3	CPU Affinity
GPU0	 X 	NV2	NV2	NV1	NV1	NODE	NODE	NODE	PIX	PHB	NODE	NODE	18-35,54-71
GPU1	NV2	 X 	NV1	NV1	NODE	NV2	NODE	NODE	PIX	PHB	NODE	NODE	18-35,54-71
GPU2	NV2	NV1	 X 	NV2	NODE	NODE	NV1	NODE	PIX	PHB	NODE	NODE	18-35,54-71
GPU3	NV1	NV1	NV2	 X 	NODE	NODE	NODE	NV2	PIX	PHB	NODE	NODE	18-35,54-71
GPU4	NV1	NODE	NODE	NODE	 X 	NV2	NV2	NV1	NODE	NODE	PIX	PHB	18-35,54-71
GPU5	NODE	NV2	NODE	NODE	NV2	 X 	NV1	NV1	NODE	NODE	PIX	PHB	18-35,54-71
GPU6	NODE	NODE	NV1	NODE	NV2	NV1	 X 	NV2	NODE	NODE	PIX	PHB	18-35,54-71
GPU7	NODE	NODE	NODE	NV2	NV1	NV1	NV2	 X 	NODE	NODE	PIX	PHB	18-35,54-71
mlx5_0	PIX	PIX	PIX	PIX	NODE	NODE	NODE	NODE	 X 	PHB	NODE	NODE
mlx5_1	PHB	PHB	PHB	PHB	NODE	NODE	NODE	NODE	PHB	 X 	NODE	NODE
mlx5_2	NODE	NODE	NODE	NODE	PIX	PIX	PIX	PIX	NODE	NODE	 X 	PHB
mlx5_3	NODE	NODE	NODE	NODE	PHB	PHB	PHB	PHB	NODE	NODE	PHB	 X
*/
// 0 2 2 2 5 5 5 5 ,
// 2 0 2 2 5 5 5 5 ,
// 2 2 0 2 5 5 5 5 ,
// 2 2 2 0 5 5 5 5 ,
// 5 5 5 5 0 2 2 2 ,
// 5 5 5 5 2 0 2 2 ,
// 5 5 5 5 2 2 0 2 ,
// 5 5 5 5 2 2 2 0
