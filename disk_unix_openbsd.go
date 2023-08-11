// The "syscall.Statfs_t" structure in OpenBSD is different from other Unix systems, and needs to be handled separately here.
// go:build openbsd

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
	diskAll := float64(fs.F_blocks * uint64(fs.F_bsize))
	diskFree := float64(fs.F_bfree * uint64(fs.F_bsize))
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
