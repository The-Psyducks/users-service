package database

import "fmt"

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

func GetInterestsByIds(ids []int32) ([]string, error) {
    var interests []string
    for _, id := range ids {
        name, err := GetInterestName(id)
        if err != nil {
            return nil, err
        }
        interests = append(interests, name)
    }
    return interests, nil
}

// IsValidInterest checks if an interest ID is valid
func IsValidInterest(id int32) bool {
    _, exists := predefinedInterests[id]
    return exists
}

// GetInterestName returns the name of an interest given its ID
func GetInterestName(id int32) (string, error) {
    name, exists := predefinedInterests[id]
    if !exists {
        return "", fmt.Errorf("interest with ID %d not found", id)
    }
    return name, nil
}