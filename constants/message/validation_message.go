package message

import "fmt"

const (
	Name = "name"
)

func RequiredMessage(label string) string {
	return fmt.Sprintf("%s is required", label)
}

func MaxLengthMessage(label string, n int) string {
	return fmt.Sprintf("%s must be at most %d characters", label, n)
}
