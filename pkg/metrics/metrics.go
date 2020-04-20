package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ranzhendong/irishman/pkg/datastruck"
	"github.com/ranzhendong/irishman/pkg/kvnuts"
	"github.com/shirou/gopsutil/mem"
	"time"
)

var (
	c datastruck.Config

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

	upstreamListCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Irishman_upstream_list_counts",
			Help: "the up list counts in nutsDB",
		}, []string{"counts"})

	statusUpListCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Irishman_up_server_counts",
			Help: "the up server counts in nutsDB",
		}, []string{"counts"})

	statusDownListCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Irishman_down_server_counts",
			Help: "the down server counts in nutsDB",
		}, []string{"counts"})
)

//setMetrics: values set
func setMetrics() {
	var (
		v   *mem.VirtualMemoryStat
		err error
	)

	//host usedMemory
	if v, err = mem.VirtualMemory(); err == nil {
		usedPercent := v.UsedPercent
		diskPercent.WithLabelValues("usedMemory").Set(usedPercent)
	}

	//up List Counts
	countUpIpPorts := 0
	countDownIpPorts := 0
	if UpstreamList, err := kvnuts.SMem(c.NutsDB.Tag.UpstreamList, c.NutsDB.Tag.UpstreamList); err == nil {
		upstreamListCounts.WithLabelValues("counts").Set(float64(len(UpstreamList)))

		for _, v := range UpstreamList {
			UpIpPorts, _ := kvnuts.SMem(c.NutsDB.Tag.Up, v)
			DownIpPorts, _ := kvnuts.SMem(c.NutsDB.Tag.Down, v)
			countUpIpPorts += len(UpIpPorts)
			countDownIpPorts += len(DownIpPorts)
		}
		statusUpListCounts.WithLabelValues("counts").Set(float64(countUpIpPorts))
		statusDownListCounts.WithLabelValues("counts").Set(float64(countDownIpPorts))
	}
}

//IrishManMetrics: entry function
func IrishManMetrics(interval int) {
	_ = c.Config()
	prometheus.MustRegister(diskPercent)
	prometheus.MustRegister(upstreamListCounts)
	prometheus.MustRegister(statusUpListCounts)
	prometheus.MustRegister(statusDownListCounts)
	go func() {
		for {
			time.Sleep(time.Duration(interval) * time.Millisecond)
			setMetrics()
		}
	}()
}
