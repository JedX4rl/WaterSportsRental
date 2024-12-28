package errors

func JsonError(message string) string {
	return `{"message": "` + message + `"}`
}
