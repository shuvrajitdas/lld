package parkinglot

func Run() {
	parkingLot := GetParkingLotInstance()

	parkingLot.Addlevel(NewLevel(1, 5))
	parkingLot.Addlevel(NewLevel(2, 5))

	car := NewCar("ABC123")
	truck := NewTruck("XYZ789")
	motorcycle := NewMotorcycle("M1234")

	parkingLot.ParkVehicle(car)
	parkingLot.ParkVehicle(truck)
	parkingLot.ParkVehicle(motorcycle)

	parkingLot.DisplayAvailability()

	parkingLot.UnparkVehicle(motorcycle)

	parkingLot.DisplayAvailability()
}
