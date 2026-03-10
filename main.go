package main

import ("fmt"; "log"; "unsafe"; "time"; w "golang.org/x/sys/windows")

type ADLTemperature struct {
	Size int32
	Temperature int32
}

var (
	dll = w.NewLazySystemDLL("atiadlxx.dll")
	ADL_Main_Control_Create = dll.NewProc("ADL_Main_Control_Create").Call
	ADL_Main_Control_Destroy = dll.NewProc("ADL_Main_Control_Destroy").Call
	ADL_Adapter_NumberOfAdapters_Get = dll.NewProc("ADL_Adapter_NumberOfAdapters_Get").Call
	ADL_Overdrive5_Temperature_Get = dll.NewProc("ADL_Overdrive5_Temperature_Get").Call
)

func main() {
	adlMalloc := w.NewCallback(func(size int32) uintptr {
		ptr, _ := w.LocalAlloc(0, uint32(size))
		return uintptr(ptr)
	})
	ADL_Main_Control_Create(adlMalloc, 1)
	defer ADL_Main_Control_Destroy()

	var adapters int32
	ADL_Adapter_NumberOfAdapters_Get(uintptr(unsafe.Pointer(&adapters)))
	if adapters == 0 { log.Fatal("GPU adapters not found") }
	fmt.Println("Adapters found:", adapters)

	temp := ADLTemperature{Size: int32(unsafe.Sizeof(ADLTemperature{}))}
	for {
		r, _, _ := ADL_Overdrive5_Temperature_Get(0, 0, uintptr(unsafe.Pointer(&temp)))
		if r == 0 {
			fmt.Printf("GPU Temperature: %d°C\n", temp.Temperature/1000)
		} else {
			fmt.Println("Temperature read failed")
		}

		time.Sleep(10 * time.Second)
	}
}