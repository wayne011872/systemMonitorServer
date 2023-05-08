package libs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wayne011872/systemMonitorServer/dao"
	"github.com/wayne011872/systemMonitorServer/mail"
)

func DetectError(in *dao.SysInfo) (bool){
	ipAddress := in.Ip
	errorRate, _ := strconv.Atoi(os.Getenv(("ERROR_RATE")))
	netErrorKbps, _ := strconv.Atoi(os.Getenv(("NETWORK_ERROR_KPBS")))
	mailContent := ""
	if int(in.CpuUsage) >= errorRate {
		var cpuMessage string
		for _,p := range in.CpuProcess {
			cpuMessage += fmt.Sprintf("<p>Pid:%-10s 程序名稱: %-30s CPU使用率:%.2f</p></br>", strconv.FormatInt(int64(p.Pid), 10), p.Name, p.Cpu)
		}
		
		cpuMailContent := fmt.Sprintf("<h3><strong>警告!!! %s 主機CPU使用率大於%d%%</strong></h3></br><p>以下是CPU使用率前10高的程序:</p></br><p>%s</p>", ipAddress, errorRate, cpuMessage)
		mailContent += cpuMailContent
	}
	if int(in.MemoryUsage) >= errorRate {
		var memoryMessage string
		for _, p := range in.MemoryProcess {
			memoryMessage += fmt.Sprintf("<p>Pid:%-10s 程序名稱: %-30s 記憶體使用率:%.2f</p></br>", strconv.FormatInt(int64(p.Pid), 10), p.Name, p.MemRate)
		}
		memoryMailContent := fmt.Sprintf("<h3><strong>警告!!! %s 主機記憶體使用率大於%d%%</strong></h3></br><p>以下是記憶體使用率前10高的程序:</p></br><p>%s</p>", ipAddress, errorRate, memoryMessage)
		mailContent += memoryMailContent
	}
	for k, d := range in.DiskUsage {
		if int(d) >= errorRate {
			diskMailContent := fmt.Sprintf("<h3><strong>警告!!! %s硬碟%d使用率大於%d%%</strong></h3></br>", ipAddress, k, errorRate)
			mailContent += diskMailContent
		}
	}
	if int(in.NetworkIn) >= netErrorKbps {
		NetMailContent := fmt.Sprintf("<h3><strong>警告!!! %s 網路輸入量大於%dKB/秒 可能為惡意攻擊</strong></h3></br>", ipAddress, netErrorKbps)
		mailContent += NetMailContent
	}
	if mailContent != "" {
		mail.SendMail("主機資源異常通知", mailContent)
		return true
	}
	return false
}