package cache

import (
	"maps"
	"math"
	"reflect"
	"unsafe"
)

// https://github.com/golang/go/blob/go1.25.0/src/internal/runtime/maps/map.go
// если кол-во элементов 8 и меньше, то вообще 0 страниц
const elemsPerTable = 896

func (c *inMemoryCache[K, V]) Realloc() {
	// Получаем указатель на map через reflect
	mapValue := reflect.ValueOf(c.cache)
	mapPointer := unsafe.Pointer(mapValue.Pointer())

	// Структура Map (может меняться в разных версиях Go)
	type Map struct {
		// The number of filled slots (i.e. the number of elements in all
		// tables). Excludes deleted slots.
		// Must be first (known by the compiler, for len() builtin).
		used uint64

		// seed is the hash seed, computed as a unique random number per map.
		seed uintptr

		// The directory of tables.
		//
		// Normally dirPtr points to an array of table pointers
		//
		// dirPtr *[dirLen]*table
		//
		// The length (dirLen) of this array is `1 << globalDepth`. Multiple
		// entries may point to the same table. See top-level comment for more
		// details.
		//
		// Small map optimization: if the map always contained
		// abi.MapGroupSlots or fewer entries, it fits entirely in a
		// single group. In that case dirPtr points directly to a single group.
		//
		// dirPtr *group
		//
		// In this case, dirLen is 0. used counts the number of used slots in
		// the group. Note that small maps never have deleted slots (as there
		// is no probe sequence to maintain).
		dirPtr unsafe.Pointer
		dirLen int

		// The number of bits to use in table directory lookups.
		globalDepth uint8

		// The number of bits to shift out of the hash for directory lookups.
		// On 64-bit systems, this is 64 - globalDepth.
		globalShift uint8

		// writing is a flag that is toggled (XOR 1) while the map is being
		// written. Normally it is set to 1 when writing, but if there are
		// multiple concurrent writers, then toggling increases the probability
		// that both sides will detect the race.
		writing uint8

		// tombstonePossible is false if we know that no table in this map
		// contains a tombstone.
		tombstonePossible bool

		// clearSeq is a sequence counter of calls to Clear. It is used to
		// detect map clears during iteration.
		clearSeq uint64
	}

	mapPtr := (*Map)(mapPointer)
	tableCount := mapPtr.dirLen
	elemCount := mapPtr.used // или просто len(c.cache)
	if int(math.Log2(float64(tableCount*elemsPerTable/int(elemCount)))) > 1 {
		mapCopy := make(map[K]V, elemCount)
		maps.Copy(mapCopy, c.cache)
		c.cache = mapCopy
	}
}
