//+build !windows

package monkey

import (
	"golang.org/x/sys/unix"
)

func mprotectCrossPage(addr uintptr, length int, prot int) {
	pageSize := unix.Getpagesize()
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := rawMemoryAccess(p, pageSize)
		err := unix.Mprotect(page, prot)
		if err != nil {
			panic(err)
		}
	}
}

// this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func copyToLocation(location uintptr, data []byte) {
	f := rawMemoryAccess(location, len(data))

	mprotectCrossPage(location, len(data), unix.PROT_READ|unix.PROT_WRITE|unix.PROT_EXEC)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), unix.PROT_READ|unix.PROT_EXEC)
}
