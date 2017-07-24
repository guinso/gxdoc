package util

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

//ToSnakeCase is to generate string name from camel case (asdQwe) to snake case (asd_qwe)
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

//MakeMD5 is to generate MD5 check summary from string (32 charactor long)
func MakeMD5(value string) string {
	md5 := md5.New()
	md5.Write([]byte(value))

	return hex.EncodeToString(md5.Sum(nil))
}
