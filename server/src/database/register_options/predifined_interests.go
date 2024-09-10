package register_options

import "fmt"

// Predefined interests
var predefinedInterests = map[int]string{
	0: "programming",
	1: "movies",
	2: "reading",
	3: "traveling",
	4: "cooking",
}

// GetAllInterestsAndIds returns all predefined interests
func GetAllInterestsAndIds() map[int]string {
    return predefinedInterests
}

func GetInterestsByIds(ids []int) ([]string, error) {
	var interests []string
	for _, id := range ids {
		name := GetInterestName(id)
		if name == "" {
			return nil, fmt.Errorf("interest with ID %d not found", id)
		}
		interests = append(interests, name)
	}
	return interests, nil
}

// IsValidInterest checks if an interest ID is valid
func IsValidInterest(id int) bool {
	_, exists := predefinedInterests[id]
	return exists
}

// GetInterestName returns the name of an interest given its ID or "" if the ID is invalid
func GetInterestName(id int) string {
	name, exists := predefinedInterests[id]
	if !exists {
		return ""
	}
	return name
}
