package notify

import (
	"fmt"
	"time"

	"github.com/wayne011872/goSterna/notify"
)

func SendLine(myLine notify.Line,lineContent string) {
	fmt.Printf("[%s] Send Error Line\n",time.Now().Format("2006-01-02 15:04:05"))
	myLine.SetLineBody(lineContent)
	myLine.SendLine()
	fmt.Printf("[%s] Send Error Line Complete\n",time.Now().Format("2006-01-02 15:04:05"))
}