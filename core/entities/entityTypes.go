package entities

// EntityType is just an alias type to hide the actual enum value
type EntityType string

// EntityTypesEnum is a helper type to enforce enum values. I've
// noticed that this is not how enums are usually implemented in
// Go, so I could be very wrong here
type EntityTypesEnum struct {
	TABLE  EntityType
	DIGEST EntityType
	ACTION EntityType
}

// EntityTypes is the actual global variable that will be used
// everywhere for getting Entity types
var EntityTypes = &EntityTypesEnum{
	TABLE:  "TABLE",
	DIGEST: "DIGEST",
	ACTION: "ACTION",
}
