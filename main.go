package main

import (
	"HETClicker/autoclicker"
	"time"
)

func main() {
	// If the windows name is empty, it interacts with all windows
	r1 := autoclicker.Initialise_Autoclicker("")
	r1.Start_Autoclicker()
	for autoclicker.IsAllAutoclickerDone() == false {
		time.Sleep(1 * time.Second)
	}
}
