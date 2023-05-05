package libs

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"

	mydao "github.com/wayne011872/systemMonitorServer/dao"
)

var (
	netUnit = []string{"Kbps", "Mbps", "Gbps", "Tbps"}
)

func TransferNetworkUnit(netData float64, unitIndex int) (float64, string) {
	if int(netData/1000) == 0 {
		return netData, netUnit[unitIndex]
	} else {
		unitIndex += 1
		return TransferNetworkUnit(netData/1000, unitIndex)
	}
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemoryPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPartitions() []disk.PartitionStat {
	parts, err := disk.Partitions(true)
	if err == nil {
		return parts
	} else {
		return nil
	}
}

func GetDiskUsageState(GetPartitions func() []disk.PartitionStat) []*disk.UsageStat {
	parts := GetPartitions()
	diskUsages := make([]*disk.UsageStat, 0)
	for _, p := range parts {
		diskUsage, _ := disk.Usage(p.Mountpoint)
		diskUsages = append(diskUsages, diskUsage)
	}
	return diskUsages
}

func GetDiskPercent(GetPartitions func() []disk.PartitionStat, GetUsage func(func() []disk.PartitionStat) []*disk.UsageStat) []float64 {
	diskUsages := GetUsage(GetPartitions)
	diskPercents := make([]float64, 0)
	for _, d := range diskUsages {
		diskPercents = append(diskPercents, d.UsedPercent)
	}
	return diskPercents
}

func GetNetInfo(networkName string) (float64, float64) {
	info, _ := net.IOCounters(true)
	for _, v := range info {
		if v.Name == networkName {
			return float64(v.BytesRecv), float64(v.BytesSent)
		}
	}
	return 0, 0
}

func GetNetPerSecond(GetNet func(string) (float64, float64), networkName string) (float64, float64) {
	oldRecv, oldSent := GetNet(networkName)
	time.Sleep(1 * time.Second)
	nowRecv, nowSent := GetNet(networkName)
	netIn := (nowRecv - oldRecv) / 1024
	netOut := (nowSent - oldSent) / 1024
	return netIn, netOut
}
func GetNetworkName() string {
	sysType := runtime.GOOS
	networkName := ""
	if sysType == "linux" {
		networkName = os.Getenv(("LINUX_NETWORK_NAME"))
	}
	if sysType == "windows" {
		networkName = os.Getenv(("WINDOWS_NETWORK_NAME"))
	}
	return networkName
}

func GetLocalIP() string {
	networkName := GetNetworkName()
	if networkName == "" {
		panic("取不到NETWORK_NAME")
	}
	addrs, _ := net.Interfaces()
	for _, v := range addrs {
		if v.Name == networkName {
			for _, addr := range v.Addrs {
				if len(strings.Split(addr.Addr, ".")) > 1 {
					return strings.Split(addr.Addr, "/")[0]
				}
			}
		}
	}
	return ""
}

func GetProcessesInfo(p *process.Process) (float64, *process.MemoryInfoStat, string) {
	pc, _ := p.CPUPercent()
	pm, _ := p.MemoryInfo()
	pn, _ := p.Name()
	return pc, pm, pn
}

func GetProcessesCPU() []*mydao.Process {
	processes, _ := process.Processes()
	topTenProcess := make(map[int]*mydao.Process, 10)
	for _, p := range processes {
		pc, _, pn := GetProcessesInfo(p)
		if pc != 0 {
			if len(topTenProcess) < 10 {
				proc := &mydao.Process{Pid: int(p.Pid), Cpu: pc, Mem: 0, Name: pn, MemRate: 0}
				topTenProcess[int(p.Pid)] = proc
			} else if len(topTenProcess) >= 10 {
				topTenProcessSorted := SortCPUProcesses(topTenProcess)
				delete(topTenProcess, topTenProcessSorted[9].Pid)
				proc := &mydao.Process{Pid: int(p.Pid), Cpu: pc, Mem: 0, Name: pn, MemRate: 0}
				topTenProcess[int(p.Pid)] = proc
			}
		}
	}
	topTenProcessSorted := SortCPUProcesses(topTenProcess)
	return topTenProcessSorted
}
func SortCPUProcesses(processes map[int]*mydao.Process) []*mydao.Process {
	var listProcess []*mydao.Process
	for _, v := range processes {
		listProcess = append(listProcess, v)
	}
	sort.Slice(listProcess, func(i, j int) bool {
		return listProcess[i].Cpu > listProcess[j].Cpu
	})
	return listProcess
}
func GetProcessesMemory() []*mydao.Process {
	processes, _ := process.Processes()
	var totalMemory uint64
	topTenProcess := make(map[int]*mydao.Process, 10)
	for _, p := range processes {
		_, pm, pn := GetProcessesInfo(p)
		if pm != nil {
			totalMemory += pm.RSS
			if len(topTenProcess) < 10 {
				proc := &mydao.Process{Pid: int(p.Pid), Cpu: 0, Mem: pm.RSS, Name: pn, MemRate: 0}
				topTenProcess[int(p.Pid)] = proc
			} else if len(topTenProcess) >= 10 {
				topTenProcessSorted := SortMemoryProcesses(topTenProcess)
				delete(topTenProcess, topTenProcessSorted[9].Pid)
				proc := &mydao.Process{Pid: int(p.Pid), Cpu: 0, Mem: pm.RSS, Name: pn, MemRate: 0}
				topTenProcess[int(p.Pid)] = proc
			}
		}
	}
	topTenProcessSorted := SortMemoryProcesses(topTenProcess)
	topTenProcessSorted = AddMemoryRateProcesses(topTenProcessSorted, totalMemory)
	return topTenProcessSorted
}

func AddMemoryRateProcesses(processes []*mydao.Process, totalMemory uint64) []*mydao.Process {
	for _, p := range processes {
		p.MemRate = (float64(p.Mem) / float64(totalMemory)) * 100
	}
	return processes
}

func SortMemoryProcesses(processes map[int]*mydao.Process) []*mydao.Process {
	var listProcess []*mydao.Process
	for _, v := range processes {
		listProcess = append(listProcess, v)
	}
	sort.Slice(listProcess, func(i, j int) bool {
		return listProcess[i].Mem > listProcess[j].Mem
	})
	return listProcess
}

func DetectError() (string,bool){
	ipAddress := GetLocalIP()
	if ipAddress == "" {
		panic("找不到ip位置")
	}
	errorRate, _ := strconv.Atoi(os.Getenv(("ERROR_RATE")))
	netErrorKbps, _ := strconv.Atoi(os.Getenv(("NETWORK_ERROR_KPBS")))
	mailContent := ""
	if int(GetCpuPercent()) >= errorRate {
		topTenCPUProcess := GetProcessesCPU()
		var cpuMessage string
		for _, p := range topTenCPUProcess {
			cpuMessage += fmt.Sprintf("<p>Pid:%-10s 程序名稱: %-30s CPU使用率:%.2f</p></br>", strconv.FormatInt(int64(p.Pid), 10), p.Name, p.Cpu)
		}

		cpuMailContent := fmt.Sprintf("<h3><strong>警告!!! %s 主機CPU使用率大於%d%%</strong></h3></br><p>以下是CPU使用率前10高的程序:</p></br><p>%s</p>", ipAddress, errorRate, cpuMessage)
		mailContent += cpuMailContent
	}
	if int(GetMemoryPercent()) >= errorRate {
		topTenMemoryProcess := GetProcessesMemory()
		var memoryMessage string
		for _, p := range topTenMemoryProcess {
			memoryMessage += fmt.Sprintf("<p>Pid:%-10s 程序名稱: %-30s 記憶體使用率:%.2f</p></br>", strconv.FormatInt(int64(p.Pid), 10), p.Name, p.MemRate)
		}
		memoryMailContent := fmt.Sprintf("<h3><strong>警告!!! %s 主機記憶體使用率大於%d%%</strong></h3></br><p>以下是記憶體使用率前10高的程序:</p></br><p>%s</p>", ipAddress, errorRate, memoryMessage)
		mailContent += memoryMailContent
	}
	diskPercents := GetDiskPercent(GetDiskPartitions, GetDiskUsageState)
	for k, d := range diskPercents {
		if int(d) >= errorRate {
			diskMailContent := fmt.Sprintf("<h3><strong>警告!!! %s硬碟%d使用率大於%d%%</strong></h3></br>", ipAddress, k, errorRate)
			mailContent += diskMailContent
		}
	}
	networkName := GetNetworkName()
	if networkName == "" {
		panic("取不到NETWORK_NAME")
	}
	netIn, _ := GetNetPerSecond(GetNetInfo, networkName)
	if int(netIn) >= netErrorKbps {
		NetMailContent := fmt.Sprintf("<h3><strong>警告!!! %s 網路輸入量大於%dKB/秒 可能為惡意攻擊</strong></h3></br>", ipAddress, netErrorKbps)
		mailContent += NetMailContent
	}
	if mailContent != "" {
		return mailContent,true
	}
	return "",false
}