package v1

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/response"
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

// 解析请求信息
func GetLogsInfo(c *gin.Context) {
	var sf model.SearchConf
	_ = c.ShouldBind(&sf)
	key := service.GetLogsInfo(sf)
	if key == "" {
		global.GVA_LOG.Error("获取失败!")
		response.FailWithMessage("获取失败", c)
		return
	} else {
		response.OkWithDetailed(gin.H{"status": "ok", "key": key}, "获取成功", c)
	}
}

//webSocket请求ping 返回pong
func WebSocket(c *gin.Context) {

	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("升级websocket失败", zap.Any("err", err))
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

/*func GetBiforstLogs()  {

	defer bifrostClient.Close()
	dataChan := make(chan []byte, 1)
	signal := make(chan int, 1)
	timeout := time.After(time.Second * 20)
	go func() {
		err := bifrostClient.WatchLog(context.Background(), token, SvrName, "access.log", dataChan, signal)
		t.Log(err)
	}()

	defer func() { signal <- 9 }()
	for {
		select {
		case data := <-dataChan:
			t.Logf(string(data))
		case <-timeout:
			t.Log("test end")
			return
		}
	}
}*/
