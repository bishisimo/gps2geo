//@Time : 2020/9/28 上午11:32
//@Author : bishisimo
package utils

import (
	"github.com/bishisimo/errlog"
	"github.com/k0kubun/pp"
)

func Print(a ...interface{}) int {
	i, err := pp.Print(a...)
	if errlog.Debug(err) {
		return 0
	}
	return i
}
func Println(a ...interface{}) int {
	i, err := pp.Println(a...)
	if errlog.Debug(err) {
		return 0
	}
	return i
}
