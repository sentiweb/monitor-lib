package common

// check if Tags in types.Notifier match with service tags
func CheckTags(tagsMap map[string]struct{}, tags []string) bool {
	if len(tagsMap) == 0 {
		return true
	}
	if len(tags) == 0 {
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

// TagsToMap transforms a tag list into a hash map (empty value, only key is used)
func TagsToMap(tags []string) map[string]struct{} {
	m := make(map[string]struct{}, len(tags))
	for _, s := range tags {
		m[s] = struct{}{}
	}
	return m
}

func MapToTags(tagMap map[string]struct{}) []string {
	m := make([]string, 0, len(tagMap))
	for k := range tagMap {
		m = append(m, k)
	}
	return m
}
