package utils

import (
	"strconv"
)

func ToString(i int) string {
	return strconv.Itoa(i)
}

func ToInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func StripLastInsertID(id int64, err error) int {
	return int(id)
}
