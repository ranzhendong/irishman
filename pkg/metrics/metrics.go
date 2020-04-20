package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"
	"time"
)

var (
	/*
	   	# HELP memory_percent memory use percent
	   # TYPE memory_percent gauge
	   memory_percent{percent="usedMemory"} 55
	*/
	diskPercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Irishman_host_memory_percent",
			Help: "the host memory use percent",
		}, []string{"percent"})

	upListCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Irishman_up_list_counts",
			Help: "the up list counts in nutsDB",
		}, []string{"counts"})

	downListCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Irishman_down_list_counts",
			Help: "the down list counts in nutsDB",
		}, []string{"counts"})
)

//setMetrics: values set
func setMetrics() {
	var (
		v   *mem.VirtualMemoryStat
		err error
	)

	//host usedMemory
	if v, err = mem.VirtualMemory(); err != nil {
	}
	usedPercent := v.UsedPercent
	diskPercent.WithLabelValues("usedMemory").Set(usedPercent)

	//up List Counts
}

//IrishManMetrics: entry function
func IrishManMetrics(interval int) {
	prometheus.MustRegister(diskPercent)
	prometheus.MustRegister(upListCounts)
	prometheus.MustRegister(downListCounts)
	go func() {
		for {
			time.Sleep(time.Duration(interval) * time.Millisecond)
			setMetrics()
		}
	}()
}
