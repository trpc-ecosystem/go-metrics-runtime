[English](README.md) | 中文

# tRPC-Go runtime 监控

[![Go Reference](https://pkg.go.dev/badge/github.com/trpc-ecosystem/go-metrics-runtime.svg)](https://pkg.go.dev/github.com/trpc-ecosystem/go-metrics-runtime)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-metrics-runtime)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-metrics-runtime)
[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green.svg)](https://github.com/trpc-ecosystem/go-metrics-runtime/blob/main/LICENSE)
[![Releases](https://img.shields.io/github/release/trpc-ecosystem/go-metrics-runtime.svg?style=flat-square)](https://github.com/trpc-ecosystem/go-metrics-runtime/releases)
[![Tests](https://github.com/trpc-ecosystem/go-metrics-runtime/actions/workflows/prc.yml/badge.svg)](https://github.com/trpc-ecosystem/go-metrics-runtime/actions/workflows/prc.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-metrics-runtime/branch/main/graph/badge.svg)](https://app.codecov.io/gh/trpc-ecosystem/go-metrics-runtime/tree/main)

每分钟定时上报 runtime 关键监控信息

## 如何使用
业务服务 import 即可：

```golang
import _ "trpc.group/trpc-go/trpc-metrics-runtime"
```

# 说明
go 的 runtime 可认为是 go 运行所需要的基础设施，主要有内存分配/垃圾回收，协程调度，操作系统及 CPU 差异性封装，pprof 支持，内置类型和反射的实现等。

go runtime package 提供了一些函数来 inspect runtime 本身的状况 (runtime.ReadMemStats, runtime.GOMAXPROCS). 这个库会把这些数据做定时上报。

这里简单介绍这些监控的含义。

注意这些并不能完全解决 runtime 相关的问题。还可以使用 gctrace 或者 pprof 来排查。
## 指标说明

|         指标          |                                  含义                                   |                                                        异常判断                                                         |
| --------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| GoroutineNum          | 协程数                                                                  | 一般和请求量还有服务耗时有关，几千 qps 的服务有 1000,2000 的协程很正常。协程数增加后不会再减少。只要不是特别多（比如>1 万）, 就不用太关注      |
| ThreadNum             | go runtime 中的 m 数量（大致可认为是线程数，不包括 c 代码里启动的线程）            | 一般和机器核数有一些关系，正常来说 10-50 都比较正常。如果使用了阻塞的 cgo, 可能会比较多。结合负载来看，不超过 100                         |
| GOMAXPROCSNum         | go 这一层代码的最大并行度（非并发度）                                        | 一般等于 go 识别到的机器核数，在容器中需要使用 uber 的 automaxprocs 库来根据配额正确设置，否则会出现很多异常                              |
| CPUCoreNum            | go 认为的机器核数（其实是给进程设置的 cpu 亲和性，sched_getaffinity)            | 无                                                                                                                     |
| PauseNsLt100usTimes   | 近 256 次 gc 中停顿<100us 的次数                                              | 这个并不能判断 gc 是否频繁。只是近 256 次 gc 的数据。一般<500us 是最多的。偶尔有超过 1ms 也正常。如果出现比较多的>10ms, 就要关注了             |
| PauseNs100_500usTimes | 近 256 次 gc 中停顿 100-500us 的次数                                           | 同上                                                                                                                    |
| PauseNs500us_1msTimes | 近 256 次 gc 中停顿 500us-1ms 的次数                                           | 同上                                                                                                                    |
| PauseNs1_10msTimes    | 近 256 次 gc 中停顿 1ms-10ms 的次数                                            | 同上                                                                                                                    |
| PauseNs10_50msTimes   | 近 256 次 gc 中停顿 10ms-50ms 的次数                                           | 同上                                                                                                                    |
| PauseNs50_100msTimes  | 近 256 次 gc 中停顿 50ms-100ms 的次数                                          | 同上                                                                                                                    |
| PauseNs100_500msTimes | 近 256 次 gc 中停顿 100ms-500ms 的次数                                         | 同上                                                                                                                    |
| PauseNs500ms_1sTimes  | 近 256 次 gc 中停顿 500ms-1s 的次数                                            | 同上                                                                                                                    |
| PauseNsBt1sTimes      | 近 256 次 gc 中停顿>1s 的次数                                                 | 同上                                                                                                                    |
| AllocMem_MB           | 当前 gc 周期到现在分配的 gc heap 中的对象的字节数                              | 看分配的进度，一般意义不大                                                                                                  |
| SysMem_MB             | go 运行时认为 go 从系统中申请的内存，包含 go 的 heap, 栈，维护一些运行时结构等的内存 | 一般可认为 go 进程占用的虚拟内存量                                                                                            |
| NextGCMem_MB          | go 本次 gc 的目标 heap 值。| 会在 AllocMem_MB 大致小于 NextGCMem_MB 时，开始本次 gc. 大致在等于 NextGCMem_MB 结束 gc. 一般可认为是 gc 时占的内存                        |
| PauseNs_us            | 历史 gc 停顿总时间                                                         |                                                                                                                         |
| GCCPUFraction_ppb     | gc 总消耗占从进程启动到现在所有 cpu 时间的比例                                 | 没有太大意义。可能值很小，因为被平均了。正常服务都是 0%,1%, 如果有 5%这样子，那就要考虑是否有问题                                      |
| MaxFdNum              | 给进程设置的允许最大 fd 数量（仅 Unix 平台）                                     |                                                                                                                        |
| CurrentFdNum          | 当前进程打开的 fd 数量（仅 Unix 平台）                                          |                                                                                                                        |
| PidNum                | 当前机器/容器中的进程数（仅 Unix 平台）                                        |                                                                                                                        |
| TcpNum                | 机器已使用（已分配+待关闭等）的所有协议套接字描述符总量（仅 Unix 平台）             | 一般在容器中会比较大（比如几万）. 对于大多数业务来说意义不大。在需要巨大量（比如几十万？) 连接的服务上可供参考                              |
| DiskFree              | 机器/根目录的磁盘可用（仅 Unix 平台，单位 GB）                                | 无 |
| DiskUsed              | 机器/根目录的磁盘已用（仅 Unix 平台，单位 GB）                                | 无 |
| DiskUsedFraction      | 机器/根目录的磁盘使用率（仅 Unix 平台）                                      | 无 |
