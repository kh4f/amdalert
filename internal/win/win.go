package win

import "golang.org/x/sys/windows"

type Notifier interface {
	Alert(msg string)
}

type MessageBoxNotifier struct{}

func (MessageBoxNotifier) Alert(msg string) {
	windows.MessageBox(0, utf16(msg), utf16("AMDAlert"), windows.MB_OK|windows.MB_ICONWARNING)
}

func utf16(s string) *uint16 {
	ptr, err := windows.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return ptr
}

func EventName(name string) *uint16 {
	return utf16(name)
}
