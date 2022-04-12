package service

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func Checkauth(username, password string) bool {
	//todo 后期可以放在数据库中， 现在先从配置文件里读取或者写死
	return true
}
