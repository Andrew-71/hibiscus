package util

import "log/slog"

// HandleWrite "handles" error in output of ResponseWriter.Write.
// Much useful very wow.
func HandleWrite(_ int, err error) {
	if err != nil {
		slog.Error("error writing response", "error", err)
	}
}
