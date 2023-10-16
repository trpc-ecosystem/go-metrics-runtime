//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

package runtime

import (
	"runtime"
	"runtime/pprof"
	"time"

	"trpc.group/trpc-go/trpc-go/metrics"
)

func init() {
	// startup runtime monitor
	go func() {
		time.Sleep(time.Second * 3) // Waiting for the framework to finish starting up
		for {
			// report data at around the 30th second of each minute to avoid
			// the issue of 0 or 2 interval between minutes.
			time.Sleep((time.Duration)(90-time.Now().Second()) * time.Second)
			// if disable plugin, stop run metric
			if GetExtraConf().Disable {
				return
			}
			RuntimeMetrics()
		}
	}()
}

// getProfData get the metrics of goroutine, thread, CPU and etc.
func getProfData() {
	profiles := pprof.Profiles()
	for _, p := range profiles {
		switch p.Name() {
		case "goroutine":
			metrics.Gauge("trpc.GoroutineNum").Set(float64(p.Count()))
		case "threadcreate":
			metrics.Gauge("trpc.ThreadNum").Set(float64(p.Count()))
		default:
		}
	}
	metrics.Gauge("trpc.GOMAXPROCSNum").Set(float64(runtime.GOMAXPROCS(0)))
	metrics.Gauge("trpc.CPUCoreNum").Set(float64(runtime.NumCPU()))
}

// getMemStats get metrics about memory
func getMemStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	var pauseNs uint64
	var pause100us, pause500us, pause1ms, pause10ms, pause50ms, pause100ms, pause500ms, pause1s, pause1sp int
	for _, ns := range mem.PauseNs {
		pauseNs += ns
		if ns < 100e3 {
			pause100us++
		} else if ns < 500e3 {
			pause500us++
		} else if ns < 1e6 {
			pause1ms++
		} else if ns < 10e6 {
			pause10ms++
		} else if ns < 50e6 {
			pause50ms++
		} else if ns < 100e6 {
			pause100ms++
		} else if ns < 500e6 {
			pause500ms++
		} else if ns < 1e9 {
			pause1s++
		} else {
			pause1sp++
		}
	}
	pauseNs /= uint64(len(mem.PauseNs))
	metrics.Gauge("trpc.PauseNsLt100usTimes").Set(float64(pause100us))
	metrics.Gauge("trpc.PauseNs100_500usTimes").Set(float64(pause500us))
	metrics.Gauge("trpc.PauseNs500us_1msTimes").Set(float64(pause1ms))
	metrics.Gauge("trpc.PauseNs1_10msTimes").Set(float64(pause10ms))
	metrics.Gauge("trpc.PauseNs10_50msTimes").Set(float64(pause50ms))
	metrics.Gauge("trpc.PauseNs50_100msTimes").Set(float64(pause100ms))
	metrics.Gauge("trpc.PauseNs100_500msTimes").Set(float64(pause500ms))
	metrics.Gauge("trpc.PauseNs500ms_1sTimes").Set(float64(pause1s))
	metrics.Gauge("trpc.PauseNsBt1sTimes").Set(float64(pause1sp))

	metrics.Gauge("trpc.AllocMem_MB").Set(float64(mem.Alloc) / 1024 / 1024)
	metrics.Gauge("trpc.SysMem_MB").Set(float64(mem.Sys) / 1024 / 1024)
	metrics.Gauge("trpc.NextGCMem_MB").Set(float64(mem.NextGC) / 1024 / 1024)
	metrics.Gauge("trpc.PauseNs_us").Set(float64(pauseNs / 1000))
	metrics.Gauge("trpc.GCCPUFraction_ppb").Set(mem.GCCPUFraction * 1000)
}
