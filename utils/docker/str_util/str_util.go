package str_util

import (
	"regexp"
	"strings"
)

type PositionMap struct {
	Start int
	End   int
}

func FindRegex(val, pattern string) string {

	re := regexp.MustCompile(pattern)
	loc := re.FindString(val)

	return loc
}

func ReplaceIgnoreCaseKeepCaseWithWrapper(upperLower bool, input, old, new, prefix, suffix string) string {
	// 构建一个忽略大小写的正则表达式
	var re *regexp.Regexp
	if upperLower {
		re = regexp.MustCompile("(?i)" + regexp.QuoteMeta(old))
	} else {
		re = regexp.MustCompile(regexp.QuoteMeta(old))
	}

	if !re.MatchString(input) {
		return ""
	}

	// 使用正则表达式替换，保留原始字符串的大小写，并添加前后缀
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// 根据匹配到的原始字符串调整替换后的大小写
		newStr := new
		for i := 0; i < len(match) && i < len(new); i++ {
			if strings.ToUpper(string(match[i])) == string(match[i]) {
				newStr = newStr[:i] + strings.ToUpper(string(new[i])) + newStr[i+1:]
			} else {
				newStr = newStr[:i] + strings.ToLower(string(new[i])) + newStr[i+1:]
			}
		}
		// 返回带有前缀和后缀的替换字符串
		return prefix + newStr + suffix
	})

	return result
}
