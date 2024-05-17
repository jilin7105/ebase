package EbaseGinResponse

type Responses interface {
	SetCode(int32)
	SetEbaseRequestID(string)
	SetMsg(string)
	SetInfo(string)
	SetData(interface{})
	SetSuccess(bool)
	Clone() Responses // 初始化/重置
}
