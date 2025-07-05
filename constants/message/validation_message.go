package message

import "fmt"

const (
	Name  = "name"
	Email = "email"
)

func RequiredMessage(label string) string {
	return fmt.Sprintf("%s is required", label)
}

func MaxLengthMessage(label string, n int) string {
	return fmt.Sprintf("%s must be at most %d characters", label, n)
}

func EmailMessage(label string) string {
	return fmt.Sprintf("%s must be a valid email address", label)
}

func InvalidFormatMessage(label string) string {
	return fmt.Sprintf("%s invalid format", label)
}

func NotFoundMessage(label string) string {
	return fmt.Sprintf("%s not found", label)
}
