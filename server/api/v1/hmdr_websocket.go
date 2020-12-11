package v1

import (
	"encoding/json"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var GlobalData []byte

//webSocket请求ping 返回pong
func WebSocket(c *gin.Context) {
	// var sc model.SocketControl
	var dataChannel <-chan []byte
	var infoChannel chan model.SocketControl
	var writeChannel chan model.SocketControl
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("升级websocket失败", zap.Any("err", err))
		return
	}
	// defer ws.Close()

	go func() {
		var sc model.SocketControl
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				global.GVA_LOG.Error("read message error", zap.Any("err", err))
			}
			if err := json.Unmarshal(message, &sc); err != nil {
				global.GVA_LOG.Error("json covert faild", zap.Any("err", err))
			}
			fmt.Println("sc:", sc)
			infoChannel <- sc
		}
	}()
	go func() {
		var sc model.SocketControl
		for {
			sc = <-infoChannel
			fmt.Println("sc2:", sc)
			if sc.Status {
				if ok := service.SelectLogWatcher(sc); ok {
					service.CalcLogWatcherCount(sc, 1)
					dataChannel, _ = service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
				} else {
					service.CreateLogWatcher(sc)
					dataChannel, _ = service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
				}
			} else {
				service.CalcLogWatcherCount(sc, -1)
			}
			if ok := service.SelectLogWatcher(sc); ok {
				writeChannel <- sc
				// dataChannel, _ = service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
			}
		}
	}()
	go func() {
		var sc model.SocketControl
		for {
			sc = <-writeChannel
			dataChannel, _ := service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
			GlobalData = <-dataChannel
		}
	}()
	go func() {
		var sc model.SocketControl
		for {
			sc = <-writeChannel
			err = ws.WriteMessage(websocket.TextMessage, GlobalData)
			if err != nil {
				break
			}
		}
	}()
	/*	for {

		mt, message, err := ws.ReadMessage()
		if err != nil {
			global.GVA_LOG.Error("read message error", zap.Any("err", err))
		}
		if err := json.Unmarshal(message, &sc); err != nil {
			global.GVA_LOG.Error("json covert faild", zap.Any("err", err))
		}
		fmt.Println("sc:", sc)
		if sc.Status {
			if ok := service.SelectLogWatcher(sc); ok {
				service.CalcLogWatcherCount(sc, 1)
				dataChannel, _ = service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
			} else {
				service.CreateLogWatcher(sc)
				dataChannel, _ = service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
			}
		} else {
			service.CalcLogWatcherCount(sc, -1)
		}
		fmt.Println("datachannel:", dataChannel)

		for {
			data := <-dataChannel
			err = ws.WriteMessage(mt, data)
			if err != nil {
				break
			}
		}
	}*/
}
