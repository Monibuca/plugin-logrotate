# LogRotate插件

可以实时查看日志输出，和日志查询
日志查询暂时只支持linux系统

## 插件地址
https://github.com/Monibuca/plugin-logrotate

# 插件引入
```go
import (
    _ "m7s.live/plugin/logrotate/v4"
)
```
## 默认配置
```yaml
logrotate:
 path: ./logs # 生成日志的目录
 size: 0 # 每个日志文件的大小，单位字节，0表示不限制
 days: 1 # 按时间分割，单位是天，即24小时
 formatter : 2006-01-02T15 # 日志文件名格式化，按照go layout格式化，默认按照小时
```
## API接口

### GET `logrotate/api/tail`
监听日志输出，该请求是一个SSE（server-sent Event）

### GET `logrotate/api/find` 
查找日志，目前只支持linux系统（使用grep），参数为 `?query=xxx`，比如 `/logrotate/api/find?query=monibuca`
### GET `logrotate/api/list` 
日志列表，列出所有日志文件
### GET `logrotate/api/open` 
打开日志，入参是 `?file=xxx`， `xxx`为文件名，文件名来自日志列表中的文件，比如 `logrotate/api/open?file=2024-01-23T09.log` 
### GET `logrotate/api/download` 
下载某个日志，入参是`?file=xxx`，参数同上打开日志
