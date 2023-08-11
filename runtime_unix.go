// go:build !windows

package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"trpc.group/trpc-go/trpc-go/metrics"
)

var pid = os.Getpid()

// GB Unit conversion for disk capacity monitoring.
const GB = 1 * 1024 * 1024 * 1024

// RuntimeMetrics runtime monitor report detail metrics of runtime every minutes
func RuntimeMetrics() {
	getProfData()
	getMemStats()

	// The following metrics are only available on non-Windows platforms
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	getFDs(ctx)
	getPidCount(ctx)
	getTcpSocket()
	getDiskUsage("/")
}

// getFDs get metrics of fd
func getFDs(ctx context.Context) {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err == nil {
		metrics.Gauge("trpc.MaxFdNum").Set(float64(limit.Cur))
	}

	out, err := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("ls /proc/%d/fd | wc -l", pid)).Output()
	if err != nil {
		return
	}
	num, err := strconv.Atoi(strings.Trim(string(out), " \n\t"))
	if err != nil {
		return
	}
	metrics.Gauge("trpc.CurrentFdNum").Set(float64(num))
}

// getPidCount get metrics of Pid
func getPidCount(ctx context.Context) {
	shell := fmt.Sprintf("ps -eLF|wc -l")
	out, err := exec.CommandContext(ctx, "bash", "-c", shell).Output()
	if err != nil {
		return
	}
	pidNum, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return
	}
	metrics.Gauge("trpc.PidNum").Set(pidNum)
}

// getTcpSocket get metrics of TCP
func getTcpSocket() {
	/// proc/net/sockstat
	st, err := os.Open("/proc/net/sockstat")
	if err != nil {
		return
	}
	data := make([]byte, 50)
	c, err := st.Read(data)
	if err != nil || c == 0 {
		return
	}
	stats := string(data[:func() int {
		for i, s := range data {
			if s == '\n' {
				return i
			}
		}
		return 0
	}()])
	sum, err := strconv.ParseFloat(strings.Split(stats, " ")[2], 64)
	if err != nil {
		return
	}
	metrics.Gauge("trpc.TcpNum").Set(sum)
}
