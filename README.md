English | [中文](README.zh_CN.md)

# tRPC-Go runtime monitoring

[![Go Reference](https://pkg.go.dev/badge/github.com/trpc-ecosystem/go-metrics-runtime.svg)](https://pkg.go.dev/github.com/trpc-ecosystem/go-metrics-runtime)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-metrics-runtime)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-metrics-runtime)
[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green.svg)](https://github.com/trpc-ecosystem/go-metrics-runtime/blob/main/LICENSE)
[![Releases](https://img.shields.io/github/release/trpc-ecosystem/go-metrics-runtime.svg?style=flat-square)](https://github.com/trpc-ecosystem/go-metrics-runtime/releases)
[![Tests](https://github.com/trpc-ecosystem/go-metrics-runtime/actions/workflows/prc.yml/badge.svg)](https://github.com/trpc-ecosystem/go-metrics-runtime/actions/workflows/prc.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-metrics-runtime/branch/main/graph/badge.svg)](https://app.codecov.io/gh/trpc-ecosystem/go-metrics-runtime/tree/main)

Reporting runtime metrics every minutes


## How to use
import the package into your service code:

```golang
import _ "trpc.group/trpc-go/trpc-metrics-runtime"
```

# Description
The runtime of Go can be considered as the infrastructure required for Go to run, primarily consisting of memory allocation, garbage collection, corouting scheduling, encapsulation of differences between operating systems and CPUs, pprof support, implementation of built-in types and reflection, etc.

Go runtime package provides functions to inspect the state of the runtime itself (such as runtime.ReadMemStats, runtime.GOMAXPROCS). This library will report this data periodically.

Here is a brief introduction to the meanings of these metrics.

Note that these may not completely solve runtime-related issues. You can also use gctrace or pprof to troubleshoot.

## Metric Description

|         Metric        |                                Description                             |                                                     Abnormal Judgment                                                  |
| --------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| GoroutineNum          | Number of Goroutine                                                    | This is generally related to the request count and service consumption time. It is normal for a service with a QPS of several thousand to have 1000 or 2000 goroutines. The number of goroutines will not decrease after it increases. As long as it is not too many (such as over 10,000), there is no need to pay too much attention to it |
| ThreadNum             | Number of M in Go runtime(It's roughly the number of threads, exclude the threads start in C code) | Generally, the number of cores in a machine is somewhat related to it, and it is normal to have 10-50 cores. If blocking cgo is used, there may be more cores involved. Based on the load, it should not exceed 100 |
| GOMAXPROCSNum         | The maximum parallelism (not concurrency) of the code in Go            | Generally, it is equal to the number of machine cores recognized by Go. In containers, the Uber's automaxprocs library should be used to correctly set quotas, otherwise many exceptions may occur |
| CPUCoreNum            | The number of logical CPUs usable by the current Go process (Actually, it is setting CPU affinity for processes, see sched_getaffinity) | None |
| PauseNsLt100usTimes   | The number of times the pause time <100us in the past 256 GC cycles | This cannot determine whether GC is frequent. It only shows that within the recent 256 GC data, the majority is less than 500us, and occasionally exceeding 1ms is normal. If there are a lot of data exceeding 10ms, attention should be paid. |
| PauseNs100_500usTimes | The number of times the pause time between 100-500us in the past 256 GC cycles | Same as above                                                                                                  |
| PauseNs500us_1msTimes | The number of times the pause time between 500us-1ms in the past 256 GC cycles | Same as above                                                                                                  |
| PauseNs1_10msTimes    | The number of times the pause time between 1-10ms in the past 256 GC cycles    | Same as above                                                                                                  |
| PauseNs10_50msTimes   | The number of times the pause time between 10-50ms in the past 256 GC cycles   | Same as above                                                                                                  |
| PauseNs50_100msTimes  | The number of times the pause time between 50-100ms in the past 256 GC cycles  | Same as above                                                                                                  |
| PauseNs100_500msTimes | The number of times the pause time between 100-500ms in the past 256 GC cycles | Same as above                                                                                                  |
| PauseNs500ms_1sTimes  | The number of times the pause time between 500ms-1s in the past 256 GC cycles  | Same as above                                                                                                  |
| PauseNsBt1sTimes      | The number of times the pause time >1s in the past 256 GC cycles               | Same as above                                                                                                  |
| AllocMem_MB           | The number of bytes of objects allocated in the GC heap from the current GC cycle until now | This depends on the progress of memory allocation, and in general, it is not significant          |
| SysMem_MB             | The memory allocated by Go from system, including heap, stack, and memory used by some structures in runtime | It can be considered as the amount of virtual memory used by the Go process      |
| NextGCMem_MB          | The target heap size of the next GC cycle | When AllocMem_MB is approximately less than NextGCMem_MB, the next GC process will begin. The GC process is expected to end when it reaches approximately equal to NextGCMem_MB. Generally speaking, this refers to the memory occupied during GC |
| PauseNs_us            | Total time of historical GC pause                                      |                                                                                                                        |
| GCCPUFraction_ppb     | The fraction of this program's available CPU time used by the GC since the program started | It is not significant. It may have a very small value as it's averaged out. Normal service rates are 0% or 1%. If it's around 5%, then it should be considered if there's an issue |
| MaxFdNum              | Maximum number of allowed file descriptors for a process (Unix only)   | |
| CurrentFdNum          | The number of file descriptors currently opened by the process (Unix only) | |
| PidNum                | The number of processes currently running in the machine/container (Unix only) | |
| TcpNum                | Total number of protocol socket descriptors in use (including allocated and pending closure) by the system (Unix only)  | Generally, the number of connections in a container is relatively large (such as tens of thousands), which is not significant for most businesses. However, it can be a reference for services that require a massive amount of connections (such as hundreds of thousands) |
| DiskFree              | Disk available in machine/root directory (Unix only, in GB)            | None |
| DiskUsed              | Disk usage of machine/root directory (Unix only, in GB)                | None |
| DiskUsedFraction      | Disk usage fraction of machine/root directory (Unix only)              | None |
