package base

type JsonResult struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	HEART_BEAT_PACKET = 0x00
	MESSAGE_PACKET    = 0x01
)
