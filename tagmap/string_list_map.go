package tagmap

import (
	"fmt"
)

type StringListMap map[string][]string

func (slm StringListMap) Keys() []string {
	keys := make([]string, len(slm))
	i := 0
	for key := range slm {
		keys[i] = key
		i++
	}
	return keys
}

func (slm StringListMap) GetOrCreate(key string) []string {
	tags, ok := slm[key]
	if !ok {
		tags = []string{}
		slm[key] = tags
	}
	return tags
}

func (slm StringListMap) HasValue(key, value string) bool {
	list, ok := slm[key]
	if !ok {
		return false
	}

	return contains(list, value)
}

func (slm StringListMap) Add(key string, value ...string) {
	list := slm.GetOrCreate(key)
	list = append(list, value...)
	slm[key] = list
}

func (slm StringListMap) AddUnique(key string, values ...string) {
	for _, value := range values {
		if !slm.HasValue(key, value) {
			slm.Add(key, value)
		}
	}
}

func (slm StringListMap) IndexOf(key, value string) int {
	list, ok := slm[key]
	if !ok {
		return -1
	}
	return indexOf(list, value)
}

func (slm StringListMap) Remove(key string, index int) error {
	list, ok := slm[key]
	if !ok {
		return nil
	}

	if index < 0 || index >= len(list) {
		return fmt.Errorf("index out of bounds: %d", index)
	}

	if index > -1 {
		list = remove(list, index)
		slm[key] = list
	}

	if len(list) == 0 {
		delete(slm, key)
	}

	return nil
}

func (slm StringListMap) RemoveFirst(key, value string) bool {
	i := slm.IndexOf(key, value)
	if i < 0 {
		return false
	}

	err := slm.Remove(key, i)
	if err != nil {
		panic(err)
	}

	return true
}

func indexOf(source []string, value string) int {
	for i, s := range source {
		if s == value {
			return i
		}
	}
	return -1
}

func remove(source []string, i int) []string {
	return append(source[:i], source[i+1:]...)
}

func contains(source []string, test string) bool {
	for _, s := range source {
		if s == test {
			return true
		}
	}
	return false
}

