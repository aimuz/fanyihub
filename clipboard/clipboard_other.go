//go:build !darwin

package clipboard

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func getClipboardContent(ctx context.Context) (string, error) {
	clipboardText, err := runtime.ClipboardGetText(ctx)
	if err != nil {
		return "", err
	}
	return clipboardText, nil
}
