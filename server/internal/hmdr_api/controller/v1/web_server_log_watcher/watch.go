package web_server_log_watcher

import (
	"encoding/json"
	"fmt"
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w *WebServerLogWatcherController) Watch(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("升级websocket失败", zap.Any("err", err))

		return
	}

	wg := new(sync.WaitGroup)
	defer wg.Wait()
	defer func() {
		err := ws.Close()
		fmt.Println("关闭ws")
		if err != nil {
			global.GVA_LOG.Error("关闭websocket错误", zap.Any("err", err))
		}
	}()

	wsReqC := make(chan []byte)
	wg.Add(1)
	go func(ch chan<- []byte) {
		defer wg.Done()
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				global.GVA_LOG.Info("读取websocket信息失败", zap.Any("err", err))
				ch <- []byte("stop")

				return
			}
			ch <- msg
		}
	}(wsReqC)

	var r metav1.WebServerLogOptions
	wsReq := <-wsReqC
	err = json.Unmarshal(wsReq, &r)
	if err != nil {
		global.GVA_LOG.Error("解析websocket请求失败", zap.Any("err", err))

		return
	}

	output, cancel, err := w.svc.WebServerLogWatchers().Watch(c, r)
	if err != nil {
		global.GVA_LOG.Error("获取Web服务器日志监看失败", zap.Any("err", err))

		return
	}
	defer cancel()

	for {
		select {
		case data := <-output:
			if data == nil {
				global.GVA_LOG.Info("log watcher 中断", zap.Any("meta", r))
				return
			}
			fmt.Println(string(data))
			err := ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				global.GVA_LOG.Warn("反馈websocket信息失败", zap.Any("err", err))

				return
			}
		case msg := <-wsReqC:
			if strings.EqualFold(string(msg), "stop") {
				global.GVA_LOG.Info("websocket关闭", zap.Any("meta", r))

				return
			}
		case <-c.Done():
			global.GVA_LOG.Info("websocket中断", zap.Any("meta", r))

			return
		}
	}
}
