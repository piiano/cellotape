package utils

import (
	"strconv"
	"testing"
)

func TestFind(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	lookFor := 3
	el, found := Find(intSlice, func(el int) bool { return el == lookFor })
	if !found {
		t.Errorf("expect to second output to be true for found %d element in %+v", el, intSlice)
	}
	if el != lookFor {
		t.Errorf("expect to first output to be the found element %d, returned value is %d", lookFor, el)
	}
}
func TestFindForMissingElement(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	lookFor := 6
	el, found := Find(intSlice, func(el int) bool { return el == lookFor })
	if found {
		t.Error("expect to second output to be false when element not found")
	}
	if el != 0 {
		t.Error("expect to first output to be zero value for the type when element not found")
	}
}

func TestIndexOf(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 3, 5}
	lookForIndex := 2
	anotherOccurrenceAtIndex := 4
	lookFor := intSlice[lookForIndex]
	foundIndex := IndexOf(intSlice, func(el int) bool { return el == lookFor })
	if intSlice[lookForIndex] != intSlice[anotherOccurrenceAtIndex] {
		t.Error("expect to perform test when the slice has multiple occurrences of the element")
	}
	if foundIndex != lookForIndex {
		t.Errorf("expect index to be the index of found element (%d)", lookForIndex)
	}
}

func TestIndexOfForMissingElement(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	lookFor := 6
	foundIndex := IndexOf(intSlice, func(el int) bool { return el == lookFor })
	if foundIndex != -1 {
		t.Errorf("expect index to be -1 when no element found (%d)", foundIndex)
	}
}

func TestLastIndexOf(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 3, 5}
	lookForIndex := 4
	anotherOccurrenceAtIndex := 2
	lookFor := intSlice[lookForIndex]
	foundIndex := LastIndexOf(intSlice, func(el int) bool { return el == lookFor })
	if intSlice[lookForIndex] != intSlice[anotherOccurrenceAtIndex] {
		t.Error("expect to perform test when the slice has multiple occurrences of the element")
	}
	if foundIndex != lookForIndex {
		t.Errorf("expect index to be the index of found element (%d)", lookForIndex)
	}
}

func TestLastIndexOfForMissingElement(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	lookFor := 6
	foundIndex := LastIndexOf(intSlice, func(el int) bool { return el == lookFor })
	if foundIndex != -1 {
		t.Errorf("expect index to be -1 when no element found (%d)", foundIndex)
	}
}

func TestFilter(t *testing.T) {
	intSlice := []int{2, 5, 4, 1, 3}
	greaterThan := 3
	filteredSlice := Filter(intSlice, func(el int) bool { return el > greaterThan })
	expected := []int{5, 4}
	for i, el := range filteredSlice {
		if el <= greaterThan {
			t.Error("expect filtered slice to include only matched elements")
		}
		if el != expected[i] {
			t.Error("expect filtered slice to preserve original order")
		}
	}
}
func TestFilterWithNoFilteredElements(t *testing.T) {
	intSlice := []int{2, 5, 4, 1, 3}
	filteredSlice := Filter(intSlice, func(el int) bool { return true })
	for i, el := range filteredSlice {
		if el != intSlice[i] {
			t.Error("expect filtered slice to be identical to the original order")
		}
	}
}

func TestFilterWithAllFiltered(t *testing.T) {
	intSlice := []int{2, 5, 4, 1, 3}
	filteredSlice := Filter(intSlice, func(el int) bool { return false })
	if len(filteredSlice) != 0 {
		t.Error("expect filtered slice to be empty")
	}
}

func TestMap(t *testing.T) {
	stringSlice := []string{"2", "5", "4", "1", "3"}
	intSlice := []int{2, 5, 4, 1, 3}
	mappedSlice := Map(stringSlice, func(el string) int {
		parseInt, _ := strconv.ParseInt(el, 10, 0)
		return int(parseInt)
	})
	for i, el := range mappedSlice {
		if el != intSlice[i] {
			t.Errorf("expect %v to equal %v", el, intSlice[i])
		}
	}
}

func TestConcatSlices(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	seg1 := intSlice[:2]
	seg2 := intSlice[2:5]
	seg3 := intSlice[5:]
	concatenatedSegments := ConcatSlices(seg1, seg2, seg3)
	for i, el := range concatenatedSegments {
		if el != intSlice[i] {
			t.Errorf("expect %v to equal %v", el, intSlice[i])
		}
	}
}
func TestConcatSlicesWithNoValue(t *testing.T) {
	emptySlice := ConcatSlices[int]()
	if len(emptySlice) != 0 {
		t.Errorf("expect concat with no arguments to return an empty slice")
	}
}
