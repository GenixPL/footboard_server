package utils

import "fmt"

func LogE(tag string, msg string, a ...interface{}) {
	fmt.Println("❌ [", tag, "]", msg)
}

func LogV(tag string, msg string, a ...interface{}) {
	fmt.Println("ℹ️ [", tag, "]", msg)
}
