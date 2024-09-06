package database

// predefined Locations
var predefinedLocations = map[int32]string{
    0: "Argentina",
    1: "Brasil",
    2: "Paraguay",
    3: "Chile",
    4: "Uruguay",
}

// GetAllInterests returns all predefined interests
func GetAllLocations() []string {
    var interests []string
    for _, interest := range predefinedLocations {
        interests = append(interests, interest)
    }
    return interests
}

// IsValidLocatio checks if an interest ID is valid
func IsValidLocation(id int32) bool {
    _, exists := predefinedLocations[id]
    return exists
}

// GetLocationName returns the name of an interest given its ID
func GetLocationName(id int32) string {
    name, exists := predefinedLocations[id]
    if !exists {
        return ""
    }
    return name
}