// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

// go:build !windows && !openbsd

package runtime

import (
	"syscall"

	"trpc.group/trpc-go/trpc-go/metrics"
)

// getDiskUsage get disk usage metrics
func getDiskUsage(path string) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	diskAll := float64(fs.Blocks * uint64(fs.Bsize))
	diskFree := float64(fs.Bfree * uint64(fs.Bsize))
	diskUsed := diskAll - diskFree
	diskUsedFraction := diskUsed / diskAll

	metrics.Gauge("trpc.DiskFree(GB)").Set(diskFree / float64(GB))
	metrics.Gauge("trpc.DiskUsed(GB)").Set(diskUsed / float64(GB))
	metrics.Gauge("trpc.DiskUsedFraction(%)").Set(diskUsedFraction)

	// The following metrics names have the same meaning as above,
	// Since Prometheus and ZhiYan-Platform do not support parentheses and percent signs
	// in a metrics name, the following duplicate metrics names have been added.
	metrics.Gauge("trpc.DiskFree_GB").Set(diskFree / float64(GB))
	metrics.Gauge("trpc.DiskUsed_GB").Set(diskUsed / float64(GB))
	metrics.Gauge("trpc.DiskUsedFraction").Set(diskUsedFraction)
}
