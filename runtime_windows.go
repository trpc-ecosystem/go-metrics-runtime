// go:build windows

package runtime

// RuntimeMetrics runtime monitor report detail metrics of runtime every minutes
func RuntimeMetrics() {
	getProfData()
	getMemStats()
}
