package ddsynfacade

func checkIfAllTagsExist(tagsArray []string, tagsToCheck []string) bool {

	// Create a map to store the elements of tagsArray for efficient lookup
	tagsArrayMap := make(map[string]bool)
	for _, item := range tagsArray {
		tagsArrayMap[item] = true
	}

	// Check if all items in arr1 exist in arr2Map
	allExist := true
	for _, item := range tagsToCheck {
		if _, ok := tagsArrayMap[item]; !ok {
			allExist = false
			break
		}
	}

	return allExist
}

func checkIfOneTagExist(tag string, tagsToCheck []string) bool {

	tagExists := false
	for _, item := range tagsToCheck {
		if item == tag {
			tagExists = true
			break
		}
	}
	return tagExists
}
