package generics

// MapKeys returns keys of a given map
func MapKeys[K comparable, V any](record map[K]V) []K {
	keys := make([]K, len(record))
	idx := 0
	for k := range record {
		keys[idx] = k
		idx++
	}
	return keys
}

// MapValues returns values of a given map
func MapValues[K comparable, V any](record map[K]V) []V {
	values := make([]V, len(record))
	idx := 0
	for _, v := range record {
		values[idx] = v
		idx++
	}
	return values
}

// MapHas checks whether a key is present in a given map
func MapHas[K comparable, V any](record map[K]V, key K) bool {
	_, has := record[key]
	return has
}

// MapContains checks whether a value is present in a given map
func MapContains[K comparable, V comparable](record map[K]V, value V) bool {
	for _, v := range record {
		if v == value {
			return true
		}
	}
	return false
}
