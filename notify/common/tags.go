package common

import(
	"github.com/sentiweb/monitor-lib/datastruct/sets"
)

// check if Tags in types.Notifier match with service tags
func CheckTags(tagsMap *sets.Set[string], tags []string) bool {
	if tagsMap.Empty() {
		return true
	}
	if len(tags) == 0 {
		// No tags for the service and the types.Notifier has some
		return false
	}
	return tagsMap.HasAny(tags)
}

// TagsToMap transforms a tag list into a hash map (empty value, only key is used)
func TagsToMap[ M *sets.Set[string]](tags []string) M {
	m := sets.New[string](len(tags))
	m.AddAll(tags)
	return m
}

func MapToTags(tagMap *sets.Set[string]) []string {
	return tagMap.Values()
}
