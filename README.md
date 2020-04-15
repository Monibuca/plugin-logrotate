# 简介

日志分割插件，带UI界面，可以实时查看日志输出，和日志查询
日志查询暂时只支持linux系统

# 插件名称

LogRotate

# 配置
```toml
[LogRotate]
Path = "log"
Size = 0
Days = 1
```
其中Path代表生成日志的目录
Size代表按大小分割，单位是字节，如果为0，则按时间分割
Days代表按时间分割，单位是天，即24小时