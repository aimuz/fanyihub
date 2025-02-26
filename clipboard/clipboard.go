package clipboard

import "context"

func GetText(ctx context.Context) (string, error) {
	return getClipboardContent(ctx)
}

