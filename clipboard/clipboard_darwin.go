//go:build darwin

package clipboard

import (
	"context"
	"errors"
	"sync"
)

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
// #import <Cocoa/Cocoa.h>
// const char* getClipboardContent() {
//     NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
//     NSString *string = [pasteboard stringForType:NSPasteboardTypeString];
//     return [string UTF8String];
// }
import "C"

var clipboardLock sync.RWMutex

func getClipboardContent(_ context.Context) (string, error) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()

	cstr := C.getClipboardContent()
	if cstr == nil {
		return "", errors.New("failed to get clipboard content")
	}
	return C.GoString(cstr), nil
}

