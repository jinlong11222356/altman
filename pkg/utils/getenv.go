package utils

import (
	"os"
)

func Getenv(var_name string) string {
	var_content := os.Getenv(var_name)
	if var_content != "" {
		return var_content
	}
	var_file_path := os.Getenv(var_name + "_FILE")
	if var_file_path == "" {
		return ""
	}
	var_file_content, err := os.ReadFile(var_file_path)
	if err != nil {
		return string(var_file_content)
	}
	return ""
}