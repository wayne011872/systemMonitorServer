package libs

import (
	"fmt"
	"strconv"

	"github.com/wayne011872/systemMonitorServer/dao"
	"github.com/wayne011872/systemMonitorServer/mail"
)

func DetectError(in *dao.SysInfo) (bool){
	ipAddress := in.Ip
	cpuUsedPercent,_ := strconv.Atoi(in.CpuUsage)
	memUsedPercent,_ := strconv.Atoi(in.MemoryStatus.MemUsedPercent)
	errorRate := in.ErrorRate
	netErrorKbps := in.NetErrorKbps
	mailHead := fmt.Sprintf("<p><strong>主機IP: %s</strong></p>",ipAddress)
	mailSeparate := "<p>-------------------------------------------<p></br>"
	mailTail := "<p><strong>請通知系統管理員處理!!!</strong></p></br>"
	mailContent := ""
	if cpuUsedPercent >= errorRate {
		infoMessage := fmt.Sprintf("<p>目前CPU使用率: %s%%</p></br>",in.CpuUsage)
		cpuMailContent := fmt.Sprintf("<h3><strong>警告!!!主機CPU使用率大於%d%%</strong></h3></br><p>以下為目前CPU資訊: </p>%s</br>%s",errorRate,infoMessage,mailSeparate)
		mailContent += cpuMailContent
	}
	if memUsedPercent >= errorRate {
		var memoryMessage string
		for _, p := range in.MemoryProcess {
			memoryMessage += fmt.Sprintf("<p>Pid:%-10s 程序名稱: %-30s 記憶體使用率:%s%%</p></br>", strconv.FormatInt(int64(p.Pid), 10), p.Name, p.MemRate)
		}
		infoMessage := fmt.Sprintf("<p>目前記憶體使用率: %s%%</p></br>",in.MemoryStatus.MemUsedPercent)
		memoryMailContent := fmt.Sprintf("<h3><strong>警告!!!主機記憶體使用率大於%d%%</strong></h3></br><p>以下為目前記憶體資訊: </p>%s</br><p>以下是記憶體使用率前10高的程序:</p></br><p>%s</p></br>%s",errorRate,infoMessage, memoryMessage,mailSeparate)
		mailContent += memoryMailContent
	}
	for _, d := range in.DiskStatus {
		diskUsedPercent,_ := strconv.Atoi(d.UsedRate)
		if diskUsedPercent >= errorRate {
			infoMessage := fmt.Sprintf("<p>目前硬碟使用率: %s%%</p></br><p>總空間: %s</p></br><p>剩餘空間: %s</p></br>",d.UsedRate,d.TotalSize,d.AvailableSize)
			diskMailContent := fmt.Sprintf("<h3><strong>警告!!!主機硬碟%s使用率大於%d%%</strong></h3></br><p>以下為目前硬碟資訊: </p>%s</br>%s",d.Drive, errorRate,infoMessage,mailSeparate)
			mailContent += diskMailContent
		}
	}
	if int(in.NetworkIn) >= netErrorKbps {
		infoMessage := fmt.Sprintf("<p>目前網路輸入量: %fKbps</p></br>",in.NetworkIn)
		NetMailContent := fmt.Sprintf("<h3><strong>警告!!!主機網路輸入量大於%dKB/秒</strong></h3></br><p>以下為目前網路流量資訊: </p>%s</br>%s",netErrorKbps,infoMessage,mailSeparate)
		mailContent += NetMailContent
	}
	if mailContent != "" {
		mailContent = fmt.Sprintf("%s\n%s\n%s",mailHead,mailContent,mailTail)
		mail.SendMail("主機資源監控異常通知", mailContent)
		return true
	}
	return false
}