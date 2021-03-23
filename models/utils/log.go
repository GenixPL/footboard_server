package utils

import "fmt"

func LogE(tag string, msg string, a ...interface{}) {
	fmt.Printf("❌ [%s] %s", tag, msg)
}

func LogV(tag string, msg string, a ...interface{}) {
	fmt.Printf("ℹ️ [%s] %s", tag, msg)
}
