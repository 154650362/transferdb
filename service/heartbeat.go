package service

import "google.golang.org/genproto/googleapis/type/datetime"

//  todo 用来上报心跳至server

type updata struct {
	addr string
	date datetime.DateTime
}
