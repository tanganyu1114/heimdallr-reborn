package logging

import (
	"gin-vue-admin/global"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type service struct {
	svc svcv1.Factory
}

func NewService(svc svcv1.Factory) svcv1.Factory {
	return &service{
		svc: svc,
	}
}

func (s *service) AgentInfos() svcv1.AgentInfoSrv {
	return newAgentInfos(s.svc)
}

func (s *service) Groups() svcv1.GroupSrv {
	return newGroups(s.svc)
}

func (s *service) Hosts() svcv1.HostSrv {
	return newHosts(s.svc)
}

func (s *service) WebServerConfigs() svcv1.WebServerConfigSrv {
	return newWebServerConfigs(s.svc)
}

func (s *service) WebServerLogWatchers() svcv1.WebServerLogWatcherSrv {
	// TODO: 新增日志服务
	// return newWebServerLogWatchers(s)
	return s.svc.WebServerLogWatchers()
}

func (s *service) WebServerStatistics() svcv1.WebServerStatisticsSrv {
	// TODO: 新增日志服务
	// return newWebServerStatistics(s)
	return s.svc.WebServerStatistics()
}

func log(level zapcore.Level, msg string, requestData, responseData any, err error) {
	var logfunc func(msg string, fields ...zapcore.Field)
	switch level {
	case zapcore.FatalLevel:
		logfunc = global.GVA_LOG.Fatal
	case zapcore.DPanicLevel:
		logfunc = global.GVA_LOG.DPanic
	case zapcore.PanicLevel:
		logfunc = global.GVA_LOG.Panic
	case zapcore.ErrorLevel:
		logfunc = global.GVA_LOG.Error
	case zapcore.WarnLevel:
		logfunc = global.GVA_LOG.Warn
	case zapcore.InfoLevel:
		logfunc = global.GVA_LOG.Info
	case zapcore.DebugLevel:
		logfunc = global.GVA_LOG.Debug
	default:
		logfunc = global.GVA_LOG.Info
	}
	logfunc(msg, zap.Any("requestData", requestData), zap.Any("responseData", responseData), zap.Any("error", err))
}
