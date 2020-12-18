package models

// Person .. Is the Model for person type.
type Person struct {
	ID         int
	FirstName  string
	LastName   string
	Age        int64
	BloodGroup string
}

// Data .. Is the Model for Data.
type Data struct {
	Persons []Person
	Other   int
}
