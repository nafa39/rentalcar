package repo

// CarRepository defines the methods for interacting with the Car entity.
type CarRepository interface {
	RentCar(carID int64, userID int64) error // Rent a car by updating its status
}
