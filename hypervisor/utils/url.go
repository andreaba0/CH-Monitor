package utils

import "strings"

func JoinUri(args ...string) string {
	var res []string = []string{}
	for _, path := range args {
		res = append(res, strings.Trim(path, "/"))
	}
	return strings.Join(res, "/")
}
