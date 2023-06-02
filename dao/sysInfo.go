package dao
type ProcessStatus struct {
	Rank		uint8	
	Pid			int	
	Name		string		`bson:"name,omitempty"`
	CpuRate		string  	`bson:"cpu_rate,omitempty"`
	MemRate 	string 		`bson:"mem_rate,omitempty"`
}

type MemoryStatus struct {
	MemTotalStorage    	string	`bson:"total"`
	MemUsedStorage		string	`bson:"used"`
	MemUsedPercent		string	`bson:"usedrate"`
}
type DiskStatus struct {
	Drive			string
	TotalSize     	string
	AvailableSize 	string
	UsedSize		string
	UsedRate		string
}

type SysInfo struct {
	Ip            	string
	CpuUsage      	string
	MemoryStatus 	*MemoryStatus
	MemoryProcess 	[]*ProcessStatus
	DiskStatus  	[]*DiskStatus
	NetworkIn     	float64
	NetworkOut    	float64
	DataTime      	string
	SendTime      	string
	ErrorRate     	int
	NetErrorKbps  	int
}
