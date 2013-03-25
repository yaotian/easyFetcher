package util

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strings"
	"unicode/utf8"
)

func Trace(content ...interface{}) {
	if show_log := os.Getenv("SHOW_LOG"); show_log != "" {
		log.Println(content, "\n")
	}
	return
}

func Show_log(content ...interface{}) {
	log.Println(content)
	return
}

func Report_Error(content ...interface{}) {
	if show_log := os.Getenv("SHOW_LOG"); show_log != "" {
		log.Println("!!!!!!!!!!=Start=!!!!!!!!!!!!!!!!!!!!!")

		//记录行号等
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Println("Func Name=" + runtime.FuncForPC(funcName).Name())
			log.Printf("file: %s    line=%d\n", file, line)
		}
		log.Println(content)
		log.Println("!!!!!!!!!!=End=!!!!!!!!!!!!!!!!!!!!!")
	}
	return
}

func CheckIsIn(all []int, for_check int) bool {
	fmt.Print(all)
	fmt.Print(for_check)
	for i := 0; i < len(all); i++ {
		if all[i] == for_check {
			fmt.Print("return true")
			return true
		}
	}
	return false
}

func CheckStringLongThan(raw_str string, length int) bool {
	return utf8.RuneCountInString(raw_str) > length
}

func CheckHasBadWords(raw_str string, badWords []string) bool {
	for _, badword := range badWords {
		if strings.Contains(raw_str, badword) {
			return true
		}
	}
	return false

}

func Remove_Html_Tag(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")

	re, _ = regexp.Compile("&\\S+;")
	src = re.ReplaceAllString(src, "")

	//去除连续的换行符
	//re, _ = regexp.Compile("\\s{2,}")
	//src = re.ReplaceAllString(src, " ")

	return (strings.TrimSpace(src))
}

func Get_A_Randome_String(raw_str_list []string) string {

	length := len(raw_str_list)
	if length > 0 {
		i := rand.Intn(length)
		return raw_str_list[i]
	}
	return ""
}
