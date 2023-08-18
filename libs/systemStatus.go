package libs

import (
	"fmt"
	"strconv"

	"github.com/wayne011872/systemMonitorServer/dao"
	"github.com/wayne011872/goSterna/notify"
	myNotify "github.com/wayne011872/systemMonitorServer/notify"
)

func DetectError(in *dao.SysInfo,mailIns notify.Mail,lineIns notify.Line) (bool){
	ipAddress := in.Ip
	cpuUsedPercent,_ := strconv.ParseFloat(in.CpuUsage,64)
	memUsedPercent,_ := strconv.ParseFloat(in.MemoryStatus.MemUsedPercent,64)
	errorRate := float64(in.ErrorRate)
	netErrorKbps := in.NetErrorKbps
	lineHead := fmt.Sprintf("主機IP: %s",ipAddress)
	lineSeparate := "-------------------------------------------\n"
	lineTail := "請通知系統管理員處理!!!\n"
	mailHead := fmt.Sprintf("<p><strong>主機IP: %s</strong></p>",ipAddress)
	mailSeparate := "<p>-------------------------------------------<p></br>"
	mailTail := "<p><strong>請通知系統管理員處理!!!</strong></p></br>"
	mailContent := ""
	lineContent := ""
	if cpuUsedPercent >= errorRate {
		infoMessage := fmt.Sprintf("<p>目前CPU使用率: %s%%</p></br>",in.CpuUsage)
		infoLineMessage := fmt.Sprintf("目前CPU使用率: %s%%\n",in.CpuUsage)
		cpuMailContent := fmt.Sprintf("<h3><strong>警告!!!主機CPU使用率大於%0.f%%</strong></h3></br><p>以下為目前CPU資訊: </p>%s</br>%s",errorRate,infoMessage,mailSeparate)
		cpuLineContent := fmt.Sprintf("警告!!!主機CPU使用率大於%0.f%%\n以下為目前CPU資訊: %s\n%s",errorRate,infoLineMessage,lineSeparate)
		mailContent += cpuMailContent
		lineContent += cpuLineContent
	}
	if memUsedPercent >= errorRate {
		var memoryMessage string
		var memoryLineMessage string
		for _, p := range in.MemoryProcess {
			memoryMessage += fmt.Sprintf("<p>Pid:%-10s 程序名稱: %-30s 記憶體使用率:%s%%</p></br>", strconv.FormatInt(int64(p.Pid), 10), p.Name, p.MemRate)
			memoryLineMessage += fmt.Sprintf("Pid:%-10s 程序名稱: %-30s 記憶體使用率:%s%%\n", strconv.FormatInt(int64(p.Pid), 10), p.Name, p.MemRate)
		}
		infoMessage := fmt.Sprintf("<p>目前記憶體使用率: %s%%</p></br>",in.MemoryStatus.MemUsedPercent)
		infoLineMessage := fmt.Sprintf("目前記憶體使用率: %s%%\n",in.MemoryStatus.MemUsedPercent)
		memoryMailContent := fmt.Sprintf("<h3><strong>警告!!!主機記憶體使用率大於%0.f%%</strong></h3></br><p>以下為目前記憶體資訊: </p>%s</br><p>以下是記憶體使用率前10高的程序:</p></br><p>%s</p></br>%s",errorRate,infoMessage, memoryMessage,mailSeparate)
		memoryLineContent := fmt.Sprintf("警告!!!主機記憶體使用率大於%0.f%%\n以下為目前記憶體資訊: %s\n以下是記憶體使用率前10高的程序:\n%s\n%s",errorRate,infoLineMessage, memoryLineMessage,lineSeparate)
		mailContent += memoryMailContent
		lineContent += memoryLineContent
	}
	for _, d := range in.DiskStatus {
		diskUsedPercent,_ := strconv.ParseFloat(d.UsedRate,64)
		if diskUsedPercent >= errorRate {
			infoMessage := fmt.Sprintf("<p>目前硬碟使用率: %s%%</p></br><p>總空間: %s</p></br><p>剩餘空間: %s</p></br>",d.UsedRate,d.TotalSize,d.AvailableSize)
			infoLineMessage := fmt.Sprintf("目前硬碟使用率: %s%%\n總空間: %s\n剩餘空間: %s\n",d.UsedRate,d.TotalSize,d.AvailableSize)
			diskMailContent := fmt.Sprintf("<h3><strong>警告!!!主機硬碟%s使用率大於%0.f%%</strong></h3></br><p>以下為目前硬碟資訊: </p>%s</br>%s",d.Drive, errorRate,infoMessage,mailSeparate)
			diskLineContent := fmt.Sprintf("警告!!!主機硬碟%s使用率大於%0.f%%\n以下為目前硬碟資訊: %s\n%s",d.Drive, errorRate,infoLineMessage,lineSeparate)
			mailContent += diskMailContent
			lineContent += diskLineContent
		}
	}
	if int(in.NetworkIn) >= netErrorKbps {
		infoMessage := fmt.Sprintf("<p>目前網路輸入量: %fKbps</p></br>",in.NetworkIn)
		infoLineMessage := fmt.Sprintf("目前網路輸入量: %fKbps\n",in.NetworkIn)
		NetMailContent := fmt.Sprintf("<h3><strong>警告!!!主機網路輸入量大於%dKB/秒</strong></h3></br><p>以下為目前網路流量資訊: </p>%s</br>%s",netErrorKbps,infoMessage,mailSeparate)
		NetLineContent := fmt.Sprintf("警告!!!主機網路輸入量大於%dKB/秒\n以下為目前網路流量資訊: %s\n%s",netErrorKbps,infoLineMessage,lineSeparate)
		mailContent += NetMailContent
		lineContent += NetLineContent
	}
	if mailContent != "" {
		mailContent = fmt.Sprintf("%s\n%s\n%s",mailHead,mailContent,mailTail)
		lineContent = fmt.Sprintf("%s\n%s\n%s",lineHead,lineContent,lineTail)
		myNotify.SendMail(mailIns,"主機資源監控異常通知", mailContent)
		myNotify.SendLine(lineIns,lineContent)
		return true
	}
	return false
}