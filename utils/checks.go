package utils

import (
	"github.com/attreios/holmes"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/host"
)

func CheckMem(size uint64) {
	h := holmes.New("debug", "UTMStack")
	v, _ := mem.VirtualMemory()
	if v.Total < size {
		h.FatalError("Your system does not have the minimal memory required: %v MB", size)
	}
}

func CheckCPU(cores int) {
	h := holmes.New("debug", "UTMStack")
	c, _ := cpu.Counts(true)
	if c < cores {
		h.FatalError("Your system does not have the minimal CPU required: %v Cores", cores)
	}
}

func CheckDistro(distro string){
	h := holmes.New("debug", "UTMStack")
	info, _ := host.Info()
	if info.Platform != distro {
		h.FatalError("Your Linux distribution (%s) is not %s", info.Platform, distro)
	}
}