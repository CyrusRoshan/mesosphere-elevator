package main

import (
	"fmt"

	"github.com/cyrusroshan/mesosphere-elevator/elevator"
)

// example usage of the elevator control system
func main() {
	control := elevator.NewControl(3) // create an elevator system with 3 elevators
	fmt.Println(control.Status())

	control.RequestPickup(5, true)

	fmt.Println(control.Status())
	control.Step()
	fmt.Println(control.Status())

	control.RequestPickup(2, true)

	fmt.Println(control.Status())
	control.Step()
	fmt.Println(control.Status())

	control.RequestPickup(1, true)

	fmt.Println(control.Status())
	control.Step()
	fmt.Println(control.Status())

	control.RequestPickup(2, true)

	fmt.Println(control.Status())
	control.Step()
	fmt.Println(control.Status())
}
