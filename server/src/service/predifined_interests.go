package service

// Predefined interests
var predefinedInterests = map[int32]string{
    0: "programming",
    1: "movies",
    2: "reading",
    3: "traveling",
    4: "cooking",
}

// GetAllInterests returns all predefined interests
func GetAllInterests() []string {
    var interests []string
    for _, interest := range predefinedInterests {
        interests = append(interests, interest)
    }
    return interests
}

// IsValidInterest checks if an interest ID is valid
func IsValidInterest(id int32) bool {
    _, exists := predefinedInterests[id]
    return exists
}

// GetInterestName returns the name of an interest given its ID
func GetInterestName(id int32) string {
    name, exists := predefinedInterests[id]
    if !exists {
        return ""
    }
    return name
}