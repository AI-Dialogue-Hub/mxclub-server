package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceToMap(t *testing.T) {
	// Test case 1: Simple integer slice
	intSlice := []int{1, 2, 2, 3, 3, 3}
	intMap := SliceToMap(intSlice, func(ele int) int {
		return ele
	})

	assert.Equal(t, 3, intMap.Size())
	assert.ElementsMatch(t, []int{1}, intMap.MustGet(1))
	assert.ElementsMatch(t, []int{2, 2}, intMap.MustGet(2))
	assert.ElementsMatch(t, []int{3, 3, 3}, intMap.MustGet(3))

	// Test case 2: Struct slice
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 20},
		{"Bob", 30},
		{"Alice", 25},
	}
	personMap := SliceToMap(people, func(p Person) string {
		return p.Name
	})

	assert.Equal(t, 2, personMap.Size())
	assert.ElementsMatch(t, []Person{{"Alice", 20}, {"Alice", 25}}, personMap.MustGet("Alice"))
	assert.ElementsMatch(t, []Person{{"Bob", 30}}, personMap.MustGet("Bob"))
}
