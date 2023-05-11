package dao

type Process struct {
	Pid     int      
	Cpu     float64  `bson:"cpu,omitempty"`
	Mem     uint64	 `bson:"mem,omitempty"`
	Name    string
	MemRate float64	 `bson:"memrate,omitempty"`
}

type SysInfo struct {
	Ip				string
	CpuUsage    	float64
	CpuProcess  	[]*Process
	MemoryUsage 	float64
	MemoryProcess 	[]*Process
	DiskUsage   	[]float64
	NetworkIn   float64
	NetworkOut  float64
	DataTime    string
	SendTime    string
	ErrorRate   int
}