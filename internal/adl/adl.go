package adl

import (
	"fmt"; "unsafe"; "time"
	"golang.org/x/sys/windows"
	"hotamd/internal/config"
	"hotamd/internal/win"
)

type ADLTemperature struct {
	Size int32
	Temperature int32
}

type ADLFanSpeedValue struct {
	Size  int32
	SpeedType int32
	FanSpeed int32
	Flags int32
}

var (
	dll = windows.NewLazySystemDLL("atiadlxx.dll")
	ADL_Main_Control_Create = dll.NewProc("ADL_Main_Control_Create").Call
	ADL_Main_Control_Destroy = dll.NewProc("ADL_Main_Control_Destroy").Call
	ADL_Adapter_NumberOfAdapters_Get = dll.NewProc("ADL_Adapter_NumberOfAdapters_Get").Call
	ADL_Overdrive5_Temperature_Get = dll.NewProc("ADL_Overdrive5_Temperature_Get").Call
	ADL_Overdrive5_FanSpeed_Get = dll.NewProc("ADL_Overdrive5_FanSpeed_Get").Call
)

func ReadGPU() (temperature int32, fan int32) {
	temp := ADLTemperature{Size: int32(unsafe.Sizeof(ADLTemperature{}))}
	fanSpeed := ADLFanSpeedValue{Size: int32(unsafe.Sizeof(ADLFanSpeedValue{}))}

	ADL_Overdrive5_Temperature_Get(0, 0, uintptr(unsafe.Pointer(&temp)))
	ADL_Overdrive5_FanSpeed_Get(0, 0, uintptr(unsafe.Pointer(&fanSpeed)))

	return temp.Temperature / 1000, fanSpeed.FanSpeed
}

func Init() {
	adlMalloc := windows.NewCallback(func(size int32) uintptr {
		ptr, _ := windows.LocalAlloc(0, uint32(size))
		return uintptr(ptr)
	})
	ADL_Main_Control_Create(adlMalloc, 1)
}

func Destroy() {
	ADL_Main_Control_Destroy()
}

func Start() {
	for {
		config.ReloadConfigIfChanged()

		temp, rpm := ReadGPU()
		if int(temp) > config.Config.MaxFanOffTemp && rpm == 0 {
			win.Alert("GPU temperature > " + fmt.Sprintf("%v", config.Config.MaxFanOffTemp) + "°C and fan is not spinning")
		} else if int(temp) > config.Config.MaxTemp {
			win.Alert("GPU temperature > " + fmt.Sprintf("%v", config.Config.MaxTemp) + "°C")
		}

		time.Sleep(10 * time.Second)
	}
}