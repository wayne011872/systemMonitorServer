package dao

type Process struct {
	Pid     int
	Cpu     float64 `bson:"cpu,omitempty"`
	Mem     uint64  `bson:"mem,omitempty"`
	Name    string
	MemRate float64 `bson:"memrate,omitempty"`
}
type DiskStatus struct {
	DiskTotalStorage	string 	`bson:"total"`
	DiskUsedStorage		string	`bson:"used"`
	DiskUsedPercent		float64	`bson:"usedrate"`
}

type SysInfo struct {
	Ip            string
	CpuUsage      float64
	CpuProcess    []*Process
	MemoryUsage   float64
	MemoryProcess []*Process
	DiskStatus  	[]*DiskStatus
	NetworkIn     float64
	NetworkOut    float64
	DataTime      string
	SendTime      string
	ErrorRate     int
	NetErrorKbps  int
}
