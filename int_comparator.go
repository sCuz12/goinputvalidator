package validator

type Comparator interface {
	IsValid(int,int) bool 
}


type GreaterThanComparator struct{}


type LessThanComparator struct{}

type GreaterEqualThanComparator struct {}
 
type LessThanEqualComparator struct {}




func (*GreaterThanComparator) IsValid( input int,  ceiling int) bool {
	return input > ceiling
}


func  (*LessThanComparator) IsValid(input int , floor int) bool {
	 return input < floor
}

func (*GreaterEqualThanComparator) IsValid(input int , ceiling int ) bool {
	return input >= ceiling 
}

func (*LessThanEqualComparator) IsValid(input int , floor int) bool {
	return input <= floor
}


