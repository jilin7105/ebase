package LinkTracking

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/kafka/ProducerAbout"
	"github.com/jilin7105/ebase/logger"
)

// 全局配置
type trackBaseInfo struct {
	config     config.LinkTracking
	serverName string
	serverType string
	producer   *ProducerAbout.KafkaProducer
}

var linkTrackConfig trackBaseInfo

// LinkTracking
type linkTrackLogData struct {
	LinkTrackID         string `json:"link_track_id"`          //追踪 id
	LinkTrackParentID   string `json:"link_track_parent_id"`   //父级 id
	LinkTrackSpan       string `json:"link_track_span"`        //事件类型
	LinkTrackDesc       string `json:"link_track_desc"`        //事件描述
	LinkTrackTime       string `json:"link_track_time"`        //触发时间
	LinkTrackActionTime string `json:"link_track_action_time"` //执行时间
	ServerName          string `json:"server_name"`            //服务名称
	ServerType          string `json:"server_type"`            //服务类型
}

type Option func(*linkTrackLogData)

func GetIsOpen() bool {
	return linkTrackConfig.config.IsOpen
}

func GetIsLog() bool {
	return linkTrackConfig.config.IsLog
}

// 初始化配置
func InitLinkTracking(conf config.Config, producer *ProducerAbout.KafkaProducer) error {
	linkTrackConfig.config = conf.LinkTrack
	linkTrackConfig.serverName = conf.ServicesName
	linkTrackConfig.serverType = conf.AppType
	//检测 kafka 是否配置
	if linkTrackConfig.config.IsOpen {
		if producer != nil {
			linkTrackConfig.producer = producer
		} else {
			return fmt.Errorf("初始化错误，开启链路追踪失败，未配置 kafka 相关信息 ，请检查配置文件 ")
		}
	}
	return nil
}

func NewLinkTrackLogData(options ...Option) (*linkTrackLogData, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Info("链路追踪发送失败 [%v]", err)
		}
	}()
	if !linkTrackConfig.config.IsOpen {
		//如果配置文件未开启，返回错误信息
		return nil, fmt.Errorf("未开启链路追踪")
	}
	l := &linkTrackLogData{

		ServerName: linkTrackConfig.serverName,
		ServerType: linkTrackConfig.serverType,
	}
	for _, opt := range options {
		opt(l)
	}
	return l, nil
}

func LinkTrackID(id string) Option {
	return func(l *linkTrackLogData) {
		l.LinkTrackID = id
	}
}

func LinkTrackParentID(id string) Option {
	return func(l *linkTrackLogData) {
		l.LinkTrackParentID = id
	}
}

func LinkTrackSpan(span string) Option {
	return func(l *linkTrackLogData) {
		l.LinkTrackSpan = span
	}
}

func LinkTrackDesc(desc string) Option {
	return func(l *linkTrackLogData) {
		l.LinkTrackDesc = desc
	}
}

func LinkTrackActionTime(time string) Option {
	return func(l *linkTrackLogData) {
		l.LinkTrackActionTime = time
	}
}

func LinkTrackTime(time string) Option {
	return func(l *linkTrackLogData) {
		l.LinkTrackTime = time
	}
}

func GetEbaseRequestID(ctx interface{}) string {

	switch ctx.(type) {
	case gin.Context:
		c, _ := ctx.(gin.Context)
		return c.GetString("EbaseRequestID")
	case context.Context:
		c, _ := ctx.(context.Context)
		return c.Value("EbaseRequestID").(string)
	default:
		return ""

	}
}

func (l *linkTrackLogData) Send() error {
	defer func() {
		if err := recover(); err != nil {
			logger.Info("链路追踪发送失败 [%v]", err)
		}
	}()
	if !linkTrackConfig.config.IsOpen {
		return nil
	}
	kp := *linkTrackConfig.producer
	marshal, err := json.Marshal(l)
	if err != nil {
		return err
	}
	_, _, err = kp.Send(string(marshal))
	if err != nil {
		return err
	}
	return nil
}
