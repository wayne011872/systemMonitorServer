package mail

import (
	"fmt"
	"time"
	"strconv"
	"os"

	"github.com/wayne011872/goSterna/mail"
)
func IsSendMail(sendTime time.Time,isSend bool) bool{
	sendInterval, _ := strconv.Atoi(os.Getenv(("SEND_MAIL_INTERVAL_TIME")))
	nowTime := time.Now()
	duration := int(nowTime.Sub(sendTime).Minutes())
	if duration >= sendInterval {
		isSend = false
	}
	return isSend
}

func SendSystemErrorMail(mailContent string) (time.Time, bool) {
	sendInterval, _ := strconv.Atoi(os.Getenv(("SEND_MAIL_INTERVAL_TIME")))
	mailTitle := "主機異常通知"
	myMail := &mail.Mail{}
	myMail.MailInit()
	myMail.SetMailTitle(mailTitle)
	fmt.Printf("[%s] Send Error E-mail\n",time.Now().Format("2006-01-02 15:04:05"))
	myMail.SetMailBody(mailContent)
	myMail.SendMail()
	sendTime := time.Now()
	isSend := true
	fmt.Printf("--------------------------------------per %d minutes------------------------------------------\n",sendInterval)
	return sendTime, isSend
}