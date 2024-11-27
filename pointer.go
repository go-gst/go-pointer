package pointer

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

var (
	mutex  sync.RWMutex
	store  = map[unsafe.Pointer]interface{}{}
	stacks = map[unsafe.Pointer][]uintptr{}
)

func Save(v interface{}) unsafe.Pointer {
	if v == nil {
		return nil
	}

	// Generate real fake C pointer.
	// This pointer will not store any data, but will be used for indexing purposes.
	// This way we are keeping a real pointer and don't get casting errors from uintptr->Pointer conversions
	var ptr unsafe.Pointer = C.malloc(C.size_t(1))
	if ptr == nil {
		panic("can't allocate 'cgo-pointer hack index pointer': ptr == nil")
	}

	pc := make([]uintptr, 10)
	n_pc := runtime.Callers(1, pc)

	mutex.Lock()
	store[ptr] = v
	stacks[ptr] = pc[:n_pc]
	mutex.Unlock()

	return ptr
}

func Restore(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		panic("gopointer Restore received nil pointer")
	}

	mutex.RLock()
	v = store[ptr]
	mutex.RUnlock()
	return
}

func Unref(ptr unsafe.Pointer) {
	if ptr == nil {
		panic("gopointer Unref received nil pointer")
	}

	mutex.Lock()

	if _, ok := store[ptr]; !ok {
		panic("received invalid go-pointer key in Unref")
	}

	delete(stacks, ptr)
	delete(store, ptr)
	mutex.Unlock()

	C.free(ptr)
}

func DumpStoredPointers() {
	mutex.Lock()

	if len(store) == 0 {
		fmt.Println("No stored gopointers left")
	}

	for p, t := range store {
		pc := stacks[p]

		frames := runtime.CallersFrames(pc)

		fmt.Printf("%p: %T\n", p, t)

		for {
			frame, more := frames.Next()
			fmt.Printf("- %s\n", frame.Function)

			// Check whether there are more frames to process after this one.
			if !more {
				break
			}
		}
	}

	mutex.Unlock()
}
