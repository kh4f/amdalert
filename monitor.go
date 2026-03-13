package main

import ("log"; "unsafe"; "time"; w "golang.org/x/sys/windows")

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
	dll = w.NewLazySystemDLL("atiadlxx.dll")
	ADL_Main_Control_Create = dll.NewProc("ADL_Main_Control_Create").Call
	ADL_Main_Control_Destroy = dll.NewProc("ADL_Main_Control_Destroy").Call
	ADL_Adapter_NumberOfAdapters_Get = dll.NewProc("ADL_Adapter_NumberOfAdapters_Get").Call
	ADL_Overdrive5_Temperature_Get = dll.NewProc("ADL_Overdrive5_Temperature_Get").Call
	ADL_Overdrive5_FanSpeed_Get = dll.NewProc("ADL_Overdrive5_FanSpeed_Get").Call
)

func utf16(s string) *uint16 {
    ptr, err := w.UTF16PtrFromString(s)
    if err != nil { panic(err) }; return ptr
}

func alert(msg string) {
	w.MessageBox(0, utf16(msg), utf16("HotAMD"), w.MB_OK|w.MB_ICONWARNING)
}

func getAdapters() int32 {
	var adapters int32
	ADL_Adapter_NumberOfAdapters_Get(uintptr(unsafe.Pointer(&adapters)))
	if adapters == 0 { log.Fatal("GPU adapters not found") }
	return adapters
}

func readGPU(adapter uintptr) (temperature int32, fan int32) {
	temp := ADLTemperature{Size: int32(unsafe.Sizeof(ADLTemperature{}))}
	fanSpeed := ADLFanSpeedValue{Size: int32(unsafe.Sizeof(ADLFanSpeedValue{}))}

	ADL_Overdrive5_Temperature_Get(0, adapter, uintptr(unsafe.Pointer(&temp)))
	ADL_Overdrive5_FanSpeed_Get(0, adapter, uintptr(unsafe.Pointer(&fanSpeed)))

	return temp.Temperature / 1000, fanSpeed.FanSpeed
}

func initADL() {
	adlMalloc := w.NewCallback(func(size int32) uintptr {
		ptr, _ := w.LocalAlloc(0, uint32(size))
		return uintptr(ptr)
	})
	ADL_Main_Control_Create(adlMalloc, 1)
}

func destroyADL() {
	ADL_Main_Control_Destroy()
}

func start() {
	for {
		getAdapters()
		temp, rpm := readGPU(0)
		if temp > 40 && rpm == 0 {
			alert("GPU temperature > 40°C and fan is not spinning")
		} else if temp > 60 {
			alert("GPU temperature > 60°C")
		}
		time.Sleep(10 * time.Second)
	}
}