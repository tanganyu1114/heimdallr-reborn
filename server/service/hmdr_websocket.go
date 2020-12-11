package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var GlobalSocketInfo = make(map[uint]*model.SocketGroup)

func SelectLogWatcher(sc model.SocketControl) bool {
	_, ok := GlobalSocketInfo[sc.GroupId]
	if !ok {
		return ok
	}
	_, ok = GlobalSocketInfo[sc.GroupId].Host[sc.HostId]
	if !ok {
		return ok
	}
	_, ok = GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName]
	if !ok {
		return ok
	}
	_, ok = GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName]
	if !ok {
		return ok
	}
	return ok
}

func CreateLogWatcher(sc model.SocketControl) {
	// 判断客户端是否存在
	if ok := SelectBifrostGroup(sc.SearchConf); ok {
		logWatcher, err := BifrostGroups[sc.GroupId].Hosts[sc.HostId].Client.WatchLog(context.Background(), BifrostGroups[sc.GroupId].Hosts[sc.HostId].HmdrHost.Token, sc.SrvName, sc.LogName)
		if err != nil {
			global.GVA_LOG.Error("init logwatcher failed", zap.Any("err", err))
		}
		GlobalSocketInfo[sc.GroupId] = &model.SocketGroup{
			Host: make(map[uint]*model.SocketHost),
		}
		GlobalSocketInfo[sc.GroupId].Host[sc.HostId] = &model.SocketHost{
			SrvName: make(map[string]*model.SocketSrvName),
		}
		GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName] = &model.SocketSrvName{
			LogName: make(map[string]*model.SocketInfo),
		}
		GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName] = &model.SocketInfo{
			LogWatcher: logWatcher,
			Count:      1,
		}
	}
}

func CalcLogWatcherCount(sc model.SocketControl, count int) {
	GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].Count += count
	if GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].Count == 0 {
		DeleteLogWater(sc)
	}
}

func DeleteLogWater(sc model.SocketControl) {
	if ok := SelectLogWatcher(sc); ok {
		GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.Close()
		delete(GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName, sc.LogName)
	}
}

// websocket 封装

/*type Connection struct {
	wsConnect  *websocket.Conn
	innerChan  chan []byte
	outerChan  chan []byte
	closerChan chan byte

	mutex    sync.Mutex // 对closeChan关闭上锁
	isClosed bool       // 防止closeChan被关闭多次
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConnect:  wsConn,
		innerChan:  make(chan []byte, 1000),
		outerChan:  make(chan []byte, 1000),
		closerChan: make(chan byte, 1),
	}
	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()
	//
	go conn.logWatcherLoop()
	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {

	select {
	case data = <-conn.innerChan:
	case <-conn.closerChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {

	select {
	case conn.outerChan <- data:
	case <-conn.closerChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (conn *Connection) Close() {
	// 线程安全，可多次调用
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closerChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// 内部实现
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}
		//阻塞在这里，等待inChan有空闲位置
		select {
		case conn.innerChan <- data:
		case <-conn.closerChan: // closerChan 感知 conn断开
			goto ERR
		}

	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outerChan:
		case <-conn.closerChan:
			goto ERR
		}
		if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}

func (conn *Connection) logWatcherLoop() {
	var (
		sc   model.SocketControl
		data []byte
	)

	for {
		select {
		case data = <-conn.innerChan:
			if err := json.Unmarshal(data, &sc); err != nil {
				global.GVA_LOG.Error("json covert faild", zap.Any("err", err))
			}
			if sc.Status {
				if ok := SelectLogWatcher(sc); ok {
					CalcLogWatcherCount(sc, 1)
				} else {
					CreateLogWatcher(sc)
				}
			} else {
				CalcLogWatcherCount(sc, -1)
			}
		case <-conn.closerChan:
			goto ERR
		}
	}

ERR:
	conn.Close()

}
*/
