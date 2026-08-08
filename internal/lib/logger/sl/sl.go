package sl

import (
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		"error",
		slog.StringValue(err.Error()),
	}
}
