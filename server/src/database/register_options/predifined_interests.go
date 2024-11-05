package register_options

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

// GetInterestName returns the name of an interest given its ID or "" if the ID is invalid
func GetInterestName(id int) string {
	name, exists := predefinedInterests[id]
	if !exists {
		return ""
	}
	return name
}
