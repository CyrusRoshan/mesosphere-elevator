package elevator

import "github.com/cyrusroshan/mesosphere-elevator/utils"

// Note that for max 16 elevators, size and time complexity of array copying, and most other operations in this pagkage are negligible, and the time required for a hashing function is relatively high compared to iterating through an array to find an elevator id

type Elevator struct {
	ID           int
	CurrentFloor int
	GoingUp      bool
	IsMoving     bool
	QueuedFloors []int
}

type elevatorControl struct {
	elevators        []*Elevator
	queuedUpFloors   []int
	queuedDownFloors []int
}

// Creates a new elevator control
func NewControl(totalElevators int) elevatorControl {
	var e elevatorControl
	e.elevators = make([]*Elevator, totalElevators)

	for i, _ := range e.elevators {
		e.elevators[i] = &Elevator{}
		e.elevators[i].ID = i
	}

	return e
}

// Returns elevator status that cannot be used to directly modify elevators
func (e *elevatorControl) Status() []Elevator {
	totalElevators := len(e.elevators)

	statusCopy := make([]Elevator, totalElevators)
	for i, _ := range e.elevators {
		statusCopy[i] = *e.elevators[i]
	}

	return statusCopy
}

// Updates an elevator with its CurrentFloor, DestinationFloor, and QueuedFloors
func (e *elevatorControl) UpdateElevator(elevatorID int, elevator Elevator) {
	*e.elevators[elevatorID] = elevator
	elevator.ID = elevatorID
}

// Requests pickup from a floor. There is no floor that the request is being made from in the example, so I assumed that this request could be given twice for a person: once to get the elevator to stop at their floor and once to enter their destination location.
func (e *elevatorControl) RequestPickup(destinationFloor int, goingUp bool) {
	e.queueFloor(destinationFloor, goingUp, false)
}

// Private queueFloor, used so that queueing a previously queued floor does not append to the elevator control queued floors, which are used when all elevators are busy and none are going in the direction that the queueued floors want to go in
func (e *elevatorControl) queueFloor(destinationFloor int, goingUp bool, previouslyQueued bool) (shouldBeRemoved bool) {
	var closestElevator *Elevator

	for _, elevator := range e.elevators {
		if !elevator.IsMoving {
			closestElevator = elevator
			continue
		}

		if elevator.GoingUp != goingUp {
			continue
		}

		if (goingUp && elevator.CurrentFloor > destinationFloor) || (!goingUp && elevator.CurrentFloor < destinationFloor) {
			continue
		}

		if closestElevator == nil {
			closestElevator = elevator
			continue
		}

		if (goingUp && elevator.CurrentFloor > closestElevator.CurrentFloor) || (!goingUp && elevator.CurrentFloor < closestElevator.CurrentFloor) {
			closestElevator = elevator
		}
	}

	if closestElevator == nil {
		if !previouslyQueued {
			if goingUp {
				e.queuedUpFloors = append(e.queuedUpFloors, destinationFloor)
			} else {
				e.queuedDownFloors = append(e.queuedDownFloors, destinationFloor)
			}
		}
		return false
	}

	if !utils.IntArrayContains(closestElevator.QueuedFloors, destinationFloor) {
		closestElevator.QueuedFloors = append(closestElevator.QueuedFloors, destinationFloor)
		if !closestElevator.IsMoving {
			closestElevator.GoingUp = destinationFloor > closestElevator.CurrentFloor
			closestElevator.IsMoving = true
		}
	}

	return true
}

type direction struct {
	GoingUp        bool
	DirectionQueue *[]int
}

// Check queued floors and update the elevator-specific queues as needed
func (e *elevatorControl) checkQueuedFloors() {
	directions := []direction{
		direction{
			true,
			&e.queuedUpFloors,
		},
		direction{
			false,
			&e.queuedDownFloors,
		},
	}

	for _, direction := range directions {
		for _, elevator := range e.elevators {
			if (!elevator.IsMoving || elevator.GoingUp == direction.GoingUp) && len(*direction.DirectionQueue) > 0 {
				qFloor := (*direction.DirectionQueue)[len(*direction.DirectionQueue)-1]
				qSuccess := e.queueFloor(qFloor, true, true)
				if qSuccess {
					*direction.DirectionQueue = (*direction.DirectionQueue)[:len(*direction.DirectionQueue)-1]
				}
			}
		}
	}
}

// Check queued floors, move elevators, update elevator-specific queues
func (e *elevatorControl) Step() {
	e.checkQueuedFloors()

	for _, elevator := range e.elevators {
		if !elevator.IsMoving {
			continue
		}

		newQueuedFloors, contains := utils.IntArrayRemoveIfContains(elevator.QueuedFloors, elevator.CurrentFloor)

		if contains {
			elevator.QueuedFloors = newQueuedFloors

			if len(elevator.QueuedFloors) == 0 {
				elevator.IsMoving = false
			}
		}

		if elevator.GoingUp {
			elevator.CurrentFloor++
		} else {
			elevator.CurrentFloor--
		}
	}
}
