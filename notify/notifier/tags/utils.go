package tags

// check if Tags in types.Notifier match with service tags
func CheckTags(tagsMap map[string]interface{}, tags []string) bool {
	if tagsMap == nil || len(tagsMap) == 0 {
		return true
	}
	if tags == nil {
		// No tags for the service and the types.Notifier has some
		return false
	}
	for _, tag := range tags {
		_, ok := tagsMap[tag]
		if ok {
			return true
		}
	}
	return false
}
