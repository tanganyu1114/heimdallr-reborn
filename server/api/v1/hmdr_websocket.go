package v1

import (
	"encoding/json"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"
	"gin-vue-admin/utils/pipe/log_watcher_pipe"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocket(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("升级websocket失败", zap.Any("err", err))
		return
	}

	defer func() {
		err := ws.Close()
		fmt.Println("关闭ws")
		if err != nil {
			global.GVA_LOG.Error("关闭websocket错误", zap.Any("err", err))
		}
	}()

	// 生成唯一码
	md5sc := utils.MD5V([]byte(time.Now().String()))
	wsChan := make(chan []byte, 0)

	// 等待组
	wg := new(sync.WaitGroup)
	defer wg.Wait()
	go func(ch chan<- []byte) {
		wg.Add(1)
		defer wg.Done()
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				global.GVA_LOG.Error("read message error", zap.Any("err", err))
				wsChan <- []byte("stop")
				return
			}
			wsChan <- message
		}
	}(wsChan)

	sc := new(model.SocketControl)
	msg := <-wsChan
	err = json.Unmarshal(msg, sc)
	if err != nil {
		global.GVA_LOG.Error("json covert failed", zap.Any("err", err))
		return
	}
	outerChannel := make(chan []byte, 0)
	if ok := service.SelectLogWatcher(*sc); ok {
		// 加入LogWatcher
		service.CalcLogWatcherCount(*sc, 1)
	} else {
		// 新增LogWatcher
		service.CreateLogWatcher(*sc)
	}
	// 确认pipe
	LogPipe := sc.GetLogPipe(service.GlobalSocketInfo)
	if LogPipe != nil {
		// 输出管道接入已有pipe
		err := LogPipe.InsertOuterChannel(md5sc, outerChannel)
		if err != nil {
			global.GVA_LOG.Error("log watch pipe error", zap.Any("err", err))
		}

	} else {
		// 新增pipe
		outerMap := make(map[string]chan<- []byte)
		outerMap[md5sc] = outerChannel
		innerChannel, _ := service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
		LogPipe, err = log_watcher_pipe.NewLogWatcherPipe(innerChannel, outerMap)
		if err != nil {
			global.GVA_LOG.Error("log watch pipe error", zap.Any("err", err))
		}
		service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogPipe = LogPipe
		// 初始化pipe后开始引流
		LogPipe.Watching()
	}

	defer func() {
		// 移除当前outerchannel
		if service.SelectLogWatcher(*sc) {
			LogPipe = sc.GetLogPipe(service.GlobalSocketInfo)
			LogPipe.Remove(md5sc)
			// 移除logwater
			service.CalcLogWatcherCount(*sc, -1)
		}
	}()

	// 处理输出管道数据及收集前端会话关闭请求
	for {
		select {
		// 从sc中获取到outerchannel然后输出数据信息
		case data := <-outerChannel:
			err := ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				global.GVA_LOG.Warn("write message err", zap.Any("err", err))
				return
			}
		case msg = <-wsChan:
			if strings.EqualFold(string(msg), "stop") {
				return
			}
		}

	}
}

////webSocket请求ping 返回pong
//func WebSocket(c *gin.Context) {
//	// var sc model.SocketControl
//	//var readChannel chan model.SocketControl
//	//var writeChannel chan model.SocketControl
//
//	readChannel := make(chan model.SocketControl, 0)
//	writeChannel := make(chan model.SocketControl, 0)
//	innerChannel := make(<-chan []byte, 0)
//	outerChannel := make(chan []byte, 0)
//	//升级get请求为webSocket协议
//
//	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		global.GVA_LOG.Error("升级websocket失败", zap.Any("err", err))
//		return
//	}
//	defer func() {
//		err := ws.Close()
//		fmt.Println("关闭ws")
//		if err != nil {
//			global.GVA_LOG.Error("关闭websocket错误", zap.Any("err", err))
//		}
//	}()
//	// defer ws.Close()
//
//	// 生成唯一码
//	ts := time.Now().String()
//	md5sc := utils.MD5V([]byte(ts))
//	wg := new(sync.WaitGroup)
//	defer wg.Wait()
//
//	// 后台读消息
//	go func() {
//		wg.Add(1)
//		defer wg.Done()
//		var sc model.SocketControl
//		for {
//			_, message, err := ws.ReadMessage()
//			if err != nil {
//				global.GVA_LOG.Error("read message error", zap.Any("err", err))
//				return
//			}
//			if string(message) == "ping" {
//				err := ws.WriteMessage(websocket.TextMessage, []byte("pong"))
//				if err != nil {
//					global.GVA_LOG.Error("read message error", zap.Any("err", err))
//					return
//				}
//				continue
//			}
//			err = json.Unmarshal(message, &sc)
//			if err != nil {
//				global.GVA_LOG.Error("json covert faild", zap.Any("err", err))
//				return
//			}
//			fmt.Println("sc:", sc)
//			readChannel <- sc
//			// 如果status为false,则退出
//			if !sc.Status {
//				return
//			}
//		}
//	}()
//
//	// 后台写消息
//	go func() {
//		wg.Add(1)
//		defer wg.Done()
//		for {
//			select {
//			// 从sc中获取到outerchannel然后输出数据信息
//			case data := <-outerChannel:
//				err := ws.WriteMessage(websocket.TextMessage, data)
//				if err != nil {
//					global.GVA_LOG.Error("write message err", zap.Any("err", err))
//					return
//				}
//			case sc := <-writeChannel:
//				if !sc.Status {
//					str := "Close Done"
//					err := ws.WriteMessage(websocket.TextMessage, []byte(str))
//					if err != nil {
//						global.GVA_LOG.Error("write message err", zap.Any("err", err))
//						return
//					}
//					return
//				}
//
//			}
//
//		}
//	}()
//
//	// 循环操作管理logwater和pipe
//	for {
//		var sc model.SocketControl
//		var LogPipe log_watcher_pipe.LogWatcherPipe
//
//		// 阻塞,当前端传值过来时运行
//		sc = <-readChannel
//		writeChannel <- sc
//		fmt.Println("md5sc:", md5sc)
//		// 初始化logwater
//		if sc.Status {
//			// 初始化logwater
//			if ok := service.SelectLogWatcher(sc); ok {
//				service.CalcLogWatcherCount(sc, 1)
//				innerChannel, _ = service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
//			} else {
//				service.CreateLogWatcher(sc)
//				innerChannel, _ = service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.GetChannels()
//			}
//			// 初始化pipe
//			if LogPipe = sc.GetLogPipe(service.GlobalSocketInfo); LogPipe != nil {
//				err := LogPipe.InsertOuterChannel(md5sc, outerChannel)
//				if err != nil {
//					global.GVA_LOG.Error("log watch pipe error", zap.Any("err", err))
//				}
//			} else {
//				outerMap := make(map[string]chan<- []byte)
//				outerMap[md5sc] = outerChannel
//				LogPipe, err = log_watcher_pipe.NewLogWatcherPipe(innerChannel, outerMap)
//				if err != nil {
//					global.GVA_LOG.Error("log watch pipe error", zap.Any("err", err))
//				}
//				service.GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogPipe = LogPipe
//			}
//			LogPipe.Watching()
//
//		} else {
//			// 移除当前outerchannel
//			if service.SelectLogWatcher(sc) {
//				LogPipe = sc.GetLogPipe(service.GlobalSocketInfo)
//				LogPipe.Remove(md5sc)
//				// 移除logwater
//				service.CalcLogWatcherCount(sc, -1)
//			}
//			return
//		}
//	}
//}
