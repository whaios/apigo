package apidoc

import "strconv"

func StrToBool(val string) bool {
	if val == "是" {
		return true
	} else if val == "否" {
		return false
	}
	b, _ := strconv.ParseBool(val)
	return b
}
