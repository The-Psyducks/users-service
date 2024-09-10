package register_options

// predefined Locations
var predefinedLocations = map[int]string{
	0: "Argentina",
	1: "Brasil",
	2: "Paraguay",
	3: "Chile",
	4: "Uruguay",
}

// GetAllInterests returns all predefined interests
func GetAllLocationsAndIds() map[int]string {
    return predefinedLocations
}
// IsValidLocatio checks if an interest ID is valid
func IsValidLocation(id int) bool {
	_, exists := predefinedLocations[id]
	return exists
}

// GetLocationName returns the name of an interest given its ID
func GetLocationName(id int) string {
	name, exists := predefinedLocations[id]
	if !exists {
		return ""
	}
	return name
}
