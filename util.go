package main

import (
	"fmt"
	"strings"
)

func hasImageSuffix(fname string) bool {
	fname = strings.ToLower(fname)
	if strings.HasSuffix(fname, ".jpg") ||
		strings.HasSuffix(fname, ".jpeg") ||
		strings.HasSuffix(fname, ".png") {
		return true
	}
	return false
}

func hasVideoSuffix(fname string) bool {
	fname = strings.ToLower(fname)
	if strings.HasSuffix(fname, ".mkv") ||
		strings.HasSuffix(fname, ".mov") ||
		strings.HasSuffix(fname, ".mp4") {
		return true
	}
	return false
}

func formatSize(size int64) string {
	v := float64(size)

	units := []string{"bytes", "KB", "MB", "GB", "TB", "PB"}

	for i, u := range units {
		if v < 1024 && i == 0 {
			return fmt.Sprintf("%d %s", size, u)
		}
		if v < 10 {
			return fmt.Sprintf("%.2f %s", v, u)
		}
		if v < 100 {
			return fmt.Sprintf("%.1f %s", v, u)
		}
		if v < 1000 {
			return fmt.Sprintf("%.0f %s", v, u)
		}
		v /= 1024
	}
	return fmt.Sprintf("%.0f %s", v*1024, units[len(units)-1])
}
