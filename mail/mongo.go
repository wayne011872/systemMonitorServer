package mail

import (
	"fmt"

	"github.com/wayne011872/goSterna/mail"
)

func SendMgoErrorMail(ip string,err error) {
	mailTitle := "資料庫異常通知"
	myMail := &mail.Mail{}
	myMail.MailInit()
	myMail.SetMailTitle(mailTitle)
	mailContent := fmt.Sprintf("%s 主機資料庫異常，異常資訊如下:\n %v",ip,err)
	myMail.SetMailBody(mailContent)
	myMail.SendMail()
}