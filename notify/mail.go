package notify

import (
	"fmt"
	"time"

	"github.com/wayne011872/goSterna/notify"
)

func SendMail(myMail notify.Mail,mailTitle,mailContent string) {
	myMail.SetMailTitle(mailTitle)
	fmt.Printf("[%s] Send Error E-mail\n",time.Now().Format("2006-01-02 15:04:05"))
	myMail.SetMailBody(mailContent)
	myMail.SendMail()
	fmt.Printf("[%s] Send Error E-mail Complete\n",time.Now().Format("2006-01-02 15:04:05"))
}