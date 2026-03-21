package win

import "golang.org/x/sys/windows"

const swHide = 0

var (
	kernel32             = windows.NewLazySystemDLL("kernel32.dll")
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow").Call
	procShowWindow       = user32.NewProc("ShowWindow").Call
)

func Alert(msg string) {
	windows.MessageBox(0, utf16(msg), utf16("AMDAlert"), windows.MB_OK|windows.MB_ICONWARNING)
}

func HideConsole() {
	consoleWindow, _, _ := procGetConsoleWindow()
	if consoleWindow == 0 {
		return
	}

	procShowWindow(consoleWindow, swHide)
}

func utf16(s string) *uint16 {
	ptr, err := windows.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return ptr
}
