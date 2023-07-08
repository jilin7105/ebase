package ebase

import (
	"github.com/jilin7105/ebase/logger"
	"time"
)

/**
记录微服务相关配置
*/

// 设置心跳推送方法
func SetHeartbeatPush(f func() error) {
	ebInstance.heartbeatPush = f
}

// 设置服务注册方法
func SetRegfunc(f func() error) {
	ebInstance.regfunc = f
}

func (e *Eb) Initmicro() {
	if e.Config.Micro.IsReg {
		if e.regfunc == nil {
			panic("未找到 regfunc 方法 \n" +
				"micro.is_reg 为 true 时  需要设置  regfunc 方法 \n" +
				"使用 ebase.SetRegfunc() 设置自己的服务注册方法")
		}
		go e.runRegfunc()
	}

	if e.Config.Micro.IsHeartPush {
		if e.heartbeatPush == nil {
			panic("未找到 heartbeatPush 方法 \n" +
				"micro.is_heart_push 为 true 时  需要设置  heartbeatPush 方法 \n" +
				"使用 ebase.SetHeartbeatPush() 设置自己的心跳推送方法")
		}
		go func(speed int64) {
			if speed == 0 {
				e.runHeartbeatPush()
			} else {
				for {
					e.runHeartbeatPush()
					duration := time.Duration(speed) * time.Second
					time.Sleep(duration)
				}
			}
		}(e.Config.Micro.HeartPushSpeed)
	}
}

//执行注册方法
func (e *Eb) runRegfunc() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("注册服务失败", r)
		}
	}()
	err := e.regfunc()
	if err != nil {
		panic(err)
	}
}

//执行心跳检测方法
func (e *Eb) runHeartbeatPush() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("心跳推送服务失败", r)
		}
	}()
	err := e.heartbeatPush()
	if err != nil {
		panic(err)
	}
}
