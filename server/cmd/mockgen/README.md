# Mock代码批量生成工具

## 📖 概述

这是一个批量扫描并执行 `//go:generate mockgen` 命令的工具，用于自动生成所有模块的 Mock 测试代码。

## 🎯 功能特点

- **自动扫描**：扫描项目中所有包含 `//go:generate mockgen` 注解的 Go 文件
- **批量执行**：一次性执行所有找到的 mockgen 命令
- **智能过滤**：只处理 mockgen 相关命令，忽略其他 generate 指令
- **详细输出**：显示每个命令的执行状态和结果
- **统计汇总**：提供成功/失败统计信息

## 🚀 快速开始

### 方式一：使用批处理脚本（推荐）

在项目根目录双击运行：

```batch
generate-mocks.bat
```

或在命令行执行：

```powershell
.\generate-mocks.bat
```

### 方式二：直接运行 Go 程序

```bash
cd E:\GO_Project\src\eir
go run cmd/mockgen/main.go
```

### 方式三：使用 go generate（传统方式）

针对单个包生成：

```bash
# Server Service 层
cd internal/apiserver/service/v1
go generate .

# Server Store 层
cd internal/apiserver/store/v1
go generate .

# CMDBClient
cd internal/pkg/cmdbclient
go generate .

# Client Endpoint
cd pkg/client/apiserver/v1/endpoint
go generate .

# Client Service
cd pkg/client/apiserver/v1/service
go generate .

# Client Transport
cd pkg/client/apiserver/v1/transport
go generate .

# Client Common
cd pkg/client/common
go generate .

# Task Scheduler Client
cd pkg/client/task_scheduler/v1/service
go generate .
```

## 📋 支持的 Mock 包

工具会自动扫描以下文件中的 `//go:generate` 注解：

| 文件名 | 路径 | Mock 内容 |
|--------|------|----------|
| `service.go` | `internal/apiserver/service/v1/` | Factory, InfraServiceFactory, ApplicationDefSrv, ApplicationSrv, PolicySrv, TroubleshootingPolicySrv, InfraDefSrv, InfraSrv, TaskSrv, TroubleshootingTaskSrv, HandlerTaskSrv, ThresholdDetectionTaskSrv, TaskSchedulerSrv, TroubleshootingTaskPublisherSrv, CustomResourceDefSrv, CustomResourceSrv, **CustomObjCLASSrv** |
| `store.go` | `internal/apiserver/store/v1/` | Factory, InfraFactory, ApplicationDefStore, ApplicationStore, CustomResourceDefStore, CustomResourceStore, PolicyStore, TroubleshootingPolicyStore, InfraDefStore, InfraStore, TaskStore, TaskSchedulerStore, TroubleshootingTaskPublisherStore, TroubleshootingTaskStore, HandlerTaskStore, ThresholdDetectionTaskStore, **CustomObjCLASStore** |
| `cmdbclient.go` | `internal/pkg/cmdbclient/` | CMDB, ClassificationService, ObjectDefinitionService, InstanceService, AttributeService, AttributeGroupService |
| `endpoint.go` | `pkg/client/apiserver/v1/endpoint/` | Factory, ApplicationDefEndpoints, ApplicationEndpoints, CustomResourceDefEndpoints, CustomResourceEndpoints, PolicyEndpoints, TroubleshootingPolicyEndpoints, TaskEndpoints, TroubleshootingTaskEndpoints, HandlerTaskEndpoints, ThresholdDetectionTaskEndpoints, TaskSchedulerEndpoints, TroubleshootingTaskPublisherEndpoints, InfraFactory, InfraDefEndpoints, InfraEndpoints, **CustomObjCLASEndpoints** |
| `service.go` | `pkg/client/apiserver/v1/service/` | Factory, ApplicationDefService, ApplicationService, CustomResourceDefService, CustomResourceService, PolicyService, TroubleshootingPolicyService, TaskService, TroubleshootingTaskService, HandlerTaskService, ThresholdDetectionTaskService, TaskSchedulerService, TroubleshootingTaskPublisherService, InfraFactory, InfraDefService, InfraService, **CustomObjCLASService** |
| `transport.go` | `pkg/client/apiserver/v1/transport/` | Factory, ApplicationDefTransport, ApplicationTransport, CustomResourceDefTransport, CustomResourceTransport, PolicyTransport, TroubleshootingPolicyTransport, TaskTransport, TroubleshootingTaskTransport, HandlerTaskTransport, ThresholdDetectionTaskTransport, TaskSchedulerTransport, TroubleshootingTaskPublisherTransport, InfraFactory, InfraDefTransport, InfraTransport, **CustomObjCLASTransport** |
| `common.go` | `pkg/client/common/` | Client |
| `service.go` | `pkg/client/task_scheduler/v1/service/` | Factory, TroubleshootingTaskSchedulerService |

## 📊 输出示例

```
🔍 Scanning for go:generate mockgen commands...
============================================================
✅ Found 8 mockgen commands

[1/8] Generating mock for service.go
  📍 File: internal/apiserver/service/v1/service.go:3
  📂 Dir:  internal/apiserver/service/v1
  🔧 Cmd:  mockgen -self_package=gitee.com/ClessLi/eir/internal/apiserver/service/v1 -destination=mock_service.go -package=v1 gitee.com/ClessLi/eir/internal/apiserver/service/v1 Factory,InfraServiceFactory,...
    ✅ Success

[2/8] Generating mock for store.go
  📍 File: internal/apiserver/store/v1/store.go:3
  📂 Dir:  internal/apiserver/store/v1
  🔧 Cmd:  mockgen -self_package=gitee.com/ClessLi/eir/internal/apiserver/store/v1 -destination=mock_store.go -package=v1 gitee.com/ClessLi/eir/internal/apiserver/store/v1 Factory,InfraFactory,...
    ✅ Success

...

============================================================
📊 Generation Summary:
  Total:    8
  Success:  8
  Failed:   0

✅ All mocks generated successfully!
```

## 🔧 工作原理

1. **扫描阶段**：
   - 遍历项目根目录下的所有 `.go` 文件
   - 匹配文件名：`service.go`, `store.go`, `cmdbclient.go`, `endpoint.go`, `transport.go`, `common.go`
   - 读取文件内容，查找 `//go:generate mockgen` 开头的行

2. **解析阶段**：
   - 提取完整的 mockgen 命令
   - 记录源文件路径、行号、工作目录

3. **执行阶段**：
   - 切换到源文件所在目录
   - 执行 mockgen 命令
   - 捕获输出并显示结果

4. **汇总阶段**：
   - 统计成功/失败数量
   - 显示汇总报告

## ⚠️ 注意事项

### 依赖安装

确保已安装 mockgen 工具：

```bash
go install go.uber.org/mock/mockgen@latest
```

验证安装：

```bash
mockgen -version
```

### 更新接口列表

当新增接口类型时，必须手动更新对应文件中的 `//go:generate` 注解。例如：

```go
//go:generate mockgen -self_package=gitee.com/ClessLi/eir/internal/apiserver/service/v1 -destination=mock_service.go -package=v1 gitee.com/ClessLi/eir/internal/apiserver/service/v1 Factory,...,NewInterfaceName
```

### 生成失败处理

如果某个命令执行失败：
1. 检查错误信息
2. 确认源文件是否有编译错误
3. 确认 mockgen 版本兼容性
4. 尝试手动执行该命令查看详细错误

### 与 go generate 的区别

| 特性 | 本工具 | go generate |
|------|--------|-------------|
| 批量执行 | ✅ 一次性执行所有 | ❌ 需要指定具体包 |
| 智能过滤 | ✅ 仅执行 mockgen | ❌ 执行所有 generate |
| 统计报告 | ✅ 详细汇总 | ❌ 无汇总 |
| 跨平台 | ✅ Windows批处理 | ✅ 原生支持 |

## 🐛 故障排除

### 问题 1：找不到 mockgen 命令

**症状**：
```
exec: "mockgen": executable file not found in %PATH%
```

**解决方案**：
```bash
go install go.uber.org/mock/mockgen@latest
```

### 问题 2：生成的 Mock代码有编译错误

**可能原因**：
- 源接口定义有语法错误
- 接口依赖的类型未正确导入
- mockgen 版本过旧

**解决方案**：
1. 检查源文件是否能正常编译
2. 更新 mockgen：`go get -u go.uber.org/mock/mockgen`
3. 查看具体错误信息

### 问题 3：部分 Mock 未生成

**症状**：只生成了部分 Mock 文件

**原因**：某些文件不在扫描范围内

**解决方案**：
- 检查文件名是否在扫描列表中
- 如需添加新文件类型，修改 `main.go` 中的 `patterns` 数组

## 📝 最佳实践

1. **提交前生成**：在提交代码前运行一次，确保 Mock 文件是最新的
2. **接口变更后立即生成**：修改接口定义后立即重新生成 Mock
3. **CI/CD集成**：在 CI流水线中加入 Mock生成步骤，确保一致性
4. **版本控制**：将生成的 Mock 文件纳入版本控制

## 🔗 相关文件

- 主程序：`cmd/mockgen/main.go`
- 批处理脚本：`generate-mocks.bat`
- 各模块 Mock 文件位置：
  - `internal/apiserver/service/v1/mock_service.go`
  - `internal/apiserver/store/v1/mock_store.go`
  - `internal/pkg/cmdbclient/mock_cmdbclient.go`
  - `pkg/client/apiserver/v1/endpoint/mock_endpoint.go`
  - `pkg/client/apiserver/v1/service/mock_service.go`
  - `pkg/client/apiserver/v1/transport/mock_transport.go`
  - `pkg/client/common/mock_common.go`
  - `pkg/client/task_scheduler/v1/service/mock_service.go`
