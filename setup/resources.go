package setup

import (
	"encoder/app"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var resourcesInterval = time.Second * 3
var resourcesDeleteInterval = time.Minute * 1
var netSent uint64 = 0
var netRecv uint64 = 0

func Resources() {
	go func() {
		// prevent resource arrays being to big
		deleteFromArrayHistory := (resourcesDeleteInterval.Seconds() / resourcesInterval.Seconds()) + 1
		for {
			time.Sleep(resourcesDeleteInterval)

			if len(app.ResourcesHistory.Cpu) > app.MaxResourcesHistory {
				app.ResourcesHistory.Cpu = app.ResourcesHistory.Cpu[int(deleteFromArrayHistory):]
			}

			if len(app.ResourcesHistory.Mem) > app.MaxResourcesHistory {
				app.ResourcesHistory.Mem = app.ResourcesHistory.Mem[int(deleteFromArrayHistory):]
			}

			if len(app.ResourcesHistory.NetIn) > app.MaxResourcesHistory {
				app.ResourcesHistory.NetIn = app.ResourcesHistory.NetIn[int(deleteFromArrayHistory):]
			}

			if len(app.ResourcesHistory.NetOut) > app.MaxResourcesHistory {
				app.ResourcesHistory.NetOut = app.ResourcesHistory.NetOut[int(deleteFromArrayHistory):]
			}
		}
	}()
	go func() {
		for {
			v, _ := mem.VirtualMemory()
			c, _ := cpu.Percent(time.Second*2, false)
			n, _ := net.IOCounters(false)

			printCpu := c[0]
			printRam := v.UsedPercent

			var printNetSent uint64 = 0
			if netSent == 0 {
				netSent = n[0].BytesSent
			} else {
				printNetSent = n[0].BytesSent - netSent
				netSent = n[0].BytesSent
			}

			var printNetRecv uint64 = 0
			if netRecv == 0 {
				netRecv = n[0].BytesRecv
			} else {
				printNetRecv = n[0].BytesRecv - netRecv
				netRecv = n[0].BytesRecv
			}

			app.ResourcesHistory.Cpu = append(app.ResourcesHistory.Cpu, printCpu)
			app.ResourcesHistory.Mem = append(app.ResourcesHistory.Mem, printRam)
			app.ResourcesHistory.NetOut = append(app.ResourcesHistory.NetOut, printNetSent/uint64(resourcesInterval.Seconds()))
			app.ResourcesHistory.NetIn = append(app.ResourcesHistory.NetIn, printNetRecv/uint64(resourcesInterval.Seconds()))

			time.Sleep(resourcesInterval)
		}
	}()
}
