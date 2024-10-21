package autoclicker

import (
	"HETClicker/logger"
	"fmt"
	"golang.org/x/sys/windows"
	"math/rand"
	"syscall"
	"time"
	"unsafe"
)

// Structure for the main object of an autoclicker
type TAutoclickerObj struct {
	procGetKeyState          *windows.LazyProc
	procMouseEvent           *windows.LazyProc
	procGetForegroundWindow  *windows.LazyProc
	procGetWindowTextW       *windows.LazyProc
	procGetWindowTextLengthW *windows.LazyProc

	changeTimeoutKeybind int
	changeJitterKeybind  int
	changeKeybindKeybind int
	pressKeybind         int

	procName  string
	isDead    bool
	timeoutMS int
	jitterMS  int
}

var goroutinesCount int = 0

/*
-	 Function Name: Initialise_Autoclicker
-	 Description: Initialise an instance of the autoclicker
*/
func Initialise_Autoclicker(processName string) TAutoclickerObj {
	userDll := windows.NewLazyDLL("user32.dll")
	return TAutoclickerObj{
		procGetKeyState:          userDll.NewProc("GetKeyState"),
		procMouseEvent:           userDll.NewProc("mouse_event"),
		procGetForegroundWindow:  userDll.NewProc("GetForegroundWindow"),
		procGetWindowTextW:       userDll.NewProc("GetWindowTextW"),
		procGetWindowTextLengthW: userDll.NewProc("GetWindowTextLengthW"),
		changeTimeoutKeybind:     0x43, // C
		changeKeybindKeybind:     0x56, // V
		changeJitterKeybind:      0x4A, // J
		pressKeybind:             0x45, // E
		procName:                 processName,
		isDead:                   false,
		timeoutMS:                2,
		jitterMS:                 200,
	}
}

/*
-	 Function Name: IsAllAutoclickerDone
-	 Description: Returns a boolean to say if all coroutines are dead
*/
func IsAllAutoclickerDone() bool {
	return goroutinesCount == 0
}

/*
-	 Function Name: Start_Autoclicker
-	 Description: Starts the autoclicker process in a goroutine
*/
func (autoclickerInst *TAutoclickerObj) Start_Autoclicker() {
	goroutinesCount++
	go func() {
		defer func() {
			goroutinesCount--
			logger.QuickLog(logger.TC_WARN, "Killed autoclicker for targetted window", map[string]interface{}{"window": autoclickerInst.procName})
		}()

		logger.QuickLog(logger.TC_INFO, "Started the autoclicker for window", map[string]interface{}{"window": autoclickerInst.procName})
		logger.QuickLog(logger.TC_INFO, "Default keybinds", map[string]interface{}{"enable/disable": KEY_MAPPINGS[autoclickerInst.pressKeybind], "change_timeout": KEY_MAPPINGS[autoclickerInst.changeTimeoutKeybind], "change_jitter": KEY_MAPPINGS[autoclickerInst.changeJitterKeybind], "change_keybind": KEY_MAPPINGS[autoclickerInst.changeKeybindKeybind], "show_parameters": "F5", "kill": "ESC"})
		for {
			// Gets foreground window handle
			hwndForeground, _, _ := autoclickerInst.procGetForegroundWindow.Call()

			// Gets the length of the title of the window
			sizeOfTitle, _, _ := autoclickerInst.procGetWindowTextLengthW.Call(uintptr(hwndForeground))

			// Gets the title of the title of the foreground window
			buffer := make([]uint16, sizeOfTitle+1)
			autoclickerInst.procGetWindowTextW.Call(hwndForeground, uintptr(unsafe.Pointer(&buffer[0])), sizeOfTitle+1)
			foregroundWindowTitle := syscall.UTF16ToString(buffer)

			// Is key for autoclicker pressed
			isValidKeyPressed := IsKeyPressed(autoclickerInst.pressKeybind, autoclickerInst.procGetKeyState)
			if (foregroundWindowTitle == autoclickerInst.procName || autoclickerInst.procName == "") && isValidKeyPressed {
				autoclickerInst.procMouseEvent.Call(0x0002, 0, 0, 0)
				autoclickerInst.procMouseEvent.Call(0x0004, 0, 0, 0)

				jitterAmount := rand.Float32() * float32(autoclickerInst.jitterMS)
				timeoutTime := float32(autoclickerInst.timeoutMS) + jitterAmount

				time.Sleep(time.Duration(timeoutTime) * time.Millisecond)

				logger.QuickLog(logger.TC_INFO, "Fired a mouse event for targetted window", map[string]interface{}{"window": autoclickerInst.procName, "cooldown": autoclickerInst.timeoutMS, "jitter": jitterAmount})
			}

			if autoclickerInst.isDead || IsKeyPressed(0x1B, autoclickerInst.procGetKeyState) {
				return // To stop the goroutine
			}

			if IsKeyPressed(0x74, autoclickerInst.procGetKeyState) {
				logger.QuickLog(logger.TC_INFO, "New keybinds", map[string]interface{}{"enable/disable": KEY_MAPPINGS[autoclickerInst.pressKeybind], "change_timeout": KEY_MAPPINGS[autoclickerInst.changeTimeoutKeybind], "change_jitter": KEY_MAPPINGS[autoclickerInst.changeJitterKeybind], "change_keybind": KEY_MAPPINGS[autoclickerInst.changeKeybindKeybind], "show_parameters": "F5", "kill": "ESC"})
				logger.QuickLog(logger.TC_INFO, "Parameters", map[string]interface{}{"Timeout": autoclickerInst.timeoutMS, "jitter": autoclickerInst.jitterMS})
				time.Sleep(1 * time.Second)
			}

			if IsKeyPressed(autoclickerInst.changeKeybindKeybind, autoclickerInst.procGetKeyState) {
				time.Sleep(200 * time.Millisecond)
				autoclickerInst.pressKeybind = GetOneKeyPressed(autoclickerInst.procGetKeyState)
				logger.QuickLog(logger.TC_INFO, "New keybinds", map[string]interface{}{"enable/disable": KEY_MAPPINGS[autoclickerInst.pressKeybind], "change_timeout": KEY_MAPPINGS[autoclickerInst.changeTimeoutKeybind], "change_jitter": KEY_MAPPINGS[autoclickerInst.changeJitterKeybind], "change_keybind": KEY_MAPPINGS[autoclickerInst.changeKeybindKeybind], "show_parameters": "F5", "kill": "ESC"})
			}

			if IsKeyPressed(autoclickerInst.changeTimeoutKeybind, autoclickerInst.procGetKeyState) {
				var tempVar int
				fmt.Printf("New timeout (in MS): ")
				fmt.Scanf("%d", &tempVar)
				autoclickerInst.timeoutMS = tempVar
			}

			if IsKeyPressed(autoclickerInst.changeJitterKeybind, autoclickerInst.procGetKeyState) {
				var tempVar int
				fmt.Printf("New jitter (in MS): ")
				fmt.Scanf("%d", &tempVar)
				autoclickerInst.jitterMS = tempVar
			}
		}
	}()
}

/*
-	 Function Name: Kill_Autoclicker
-	 Description: Kills the autoclicker instance
*/
func (autoclickerInst *TAutoclickerObj) Kill_Autoclicker() {
	autoclickerInst.isDead = true
}
