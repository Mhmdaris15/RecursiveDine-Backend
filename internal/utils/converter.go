package utils

import (
	"fmt"
)

func UintToString(u uint) string {
	return fmt.Sprintf("%d", u)
}

func StringToUint(s string) uint {
	// Simple conversion - in production, use proper error handling
	switch s {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	default:
		return 0
	}
}
