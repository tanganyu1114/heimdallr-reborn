package middleware

import (
	"bytes"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	"github.com/tanganyu1114/heimdallr-reborn/model"
	"github.com/tanganyu1114/heimdallr-reborn/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/service"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var userId int
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				global.GVA_LOG.Error("read body from request error:", zap.Any("err", err))
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}
		if claims, ok := c.Get("claims"); ok {
			waitUse := claims.(*request.CustomClaims)
			userId = int(waitUse.ID)
		} else {
			id, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
			if err != nil {
				userId = 0
			}
			userId = id
		}
		record := model.SysOperationRecord{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
			UserID: userId,
		}
		// 检查是否为SDK登录
		if isSDKLogin, ok := c.Get("is_sdk_login"); ok && isSDKLogin.(bool) {
			record.IsSDKLogin = true
		}
		// 存在某些未知错误 TODO
		//values := c.Request.Header.Values("content-type")
		//if len(values) >0 && strings.Contains(values[0], "boundary") {
		//	record.Body = "file"
		//}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Now().Sub(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency

		// 检查是否为二进制响应（如 Excel 文件下载）
		contentType := c.Writer.Header().Get("Content-Type")
		if contentType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" ||
			contentType == "application/octet-stream" ||
			contentType == "application/pdf" ||
			strings.HasPrefix(contentType, "image/") {
			// 二进制文件不记录响应内容，避免数据库错误
			record.Resp = "[Binary File - Not Logged]"
		} else {
			record.Resp = writer.body.String()
		}

		if err := service.CreateSysOperationRecord(record); err != nil {
			global.GVA_LOG.Error("create operation record error:", zap.Any("err", err))
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
