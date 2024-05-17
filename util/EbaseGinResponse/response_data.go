package EbaseGinResponse

type Response struct {
	EbaseRequestID string `json:"EbaseRequestID,omitempty"`
	Code           int32  `json:"code,omitempty"`
	Info           string `json:"info,omitempty"`
	Msg            string `json:"msg,omitempty"`
	Status         string `json:"status,omitempty"`
}

type response struct {
	Response
	Data any `json:"data"`
}

func (e *response) SetEbaseRequestID(id string) {
	e.EbaseRequestID = id
}

func (e *response) SetMsg(msg string) {
	e.Msg = msg
}

func (e *response) SetInfo(info string) {
	e.Info = info
}

func (e *response) SetCode(code int32) {
	e.Code = code
}

func (e *response) SetSuccess(success bool) {
	if !success {
		e.Status = "error"
	}
}

type Page struct {
	Count     int `json:"count"`
	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`
}

type page struct {
	Page
	List any `json:"list"`
}

func (e *response) SetData(data any) {
	e.Data = data
}

func (e response) Clone() Responses {
	return &e
}
