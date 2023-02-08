package shuangpin

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/mozillazg/go-pinyin"
)

var title = []string{"animal", "caijing", "car", "chengyu", "food", "IT", "law", "lishimingren", "medical", "poem"}

func GetWords(names []string, nums int) (result []string, err error) {
	if len(names) == 0 {
		names = title
	}

	num := (nums / len(names)) + 1

	for _, name := range names {
		words, err2 := getWordForDir(name, num)
		if err2 == nil {
			result = append(result, words...)
		}
	}

	return
}

func getWordForDir(name string, nums int) (result []string, err error) {
	rand.Seed(time.Now().Unix())

	var f []byte
	f, err = ioutil.ReadFile("./data/THUOCL_" + name + ".txt")
	if err != nil {
		fmt.Println("read fail", err)
		return
	}

	words := strings.Split(string(f), "\n")

	for i := 0; i < nums; i++ {
		r := rand.Intn(len(words))
		w := words[r]
		result = append(result, strings.Split(w, "\t")[0])
	}

	return
}

// GetPingYin 传入字符串，获取每个字符串的拼音，如果字符串中存在 英文字母 或 数字，则返回它本身。
func GetPingYin(words string) ([]string, error) {

	tempPinYin := pinyin.LazyPinyin(words, pinyin.NewArgs()) // 如果存在 英文字母 则会直接忽略

	var result []string

	i := 0

	for _, word := range words {

		if ('a' <= word && word <= 'z') || ('A' <= word && word <= 'Z') || ('0' <= word && word <= '9') {
			result = append(result, string(word))
		} else {
			if i < len(tempPinYin) {
				result = append(result, tempPinYin[i])
				i++
			}
		}
	}

	return result, nil
}

func PinYinTurnShuangPin(wordsPinYin []string) []string {

	var result []string

	for _, word := range wordsPinYin {
		result = append(result, turnShuangPin(word))
	}

	return result
}

func CheckResult(words string, wordsPinYin, wordsShuangPin []string) []string {

	for i, word := range words {
		if 'a' <= word && word <= 'z' {
			wordsShuangPin[i] = wordsPinYin[i]
		}
	}

	return wordsShuangPin
}

// turnShuangPin 将拼音py 转化成双拼
func turnShuangPin(py string) string {

	dict := initDict()

	sm, ym := false, false
	rSm, rYm := "", ""

	if len(py) == 6 {
		// 2 + 4
		rSm, sm = reqDict2(py[:2])
		rYm, ym = reqDict4(py[2:])
		if sm && ym {
			return dict[rSm] + dict[rYm]
		}
	}

	if len(py) == 5 {
		// 2 + 3
		sm, ym = false, false
		rSm, sm = reqDict2(py[:2])
		rYm, ym = reqDict3(py[2:])
		if sm && ym {
			return dict[rSm] + dict[rYm]
		}

		// 1 + 4
		sm, ym = false, false
		rSm, sm = reqDict1(py[:1])
		rYm, ym = reqDict4(py[1:])
		if sm && ym {
			return rSm + dict[rYm]
		}
	}

	if len(py) == 4 {
		// 1 + 3
		sm, ym = false, false
		rSm, sm = reqDict1(py[:1])
		rYm, ym = reqDict3(py[1:])
		if sm && ym {
			return rSm + dict[rYm]
		}

		// 2 + 2
		sm, ym = false, false
		rSm, sm = reqDict2(py[:2])
		rYm, ym = reqDict2(py[2:])
		if sm && ym {
			return dict[rSm] + dict[rYm]
		}
	}

	if len(py) == 3 {
		// 1 + 2
		sm, ym = false, false
		rSm, sm = reqDict1(py[:1])
		rYm, ym = reqDict2(py[1:])
		if sm && ym {
			return rSm + dict[rYm]
		}

		// 2 + 1
		sm, ym = false, false
		rSm, sm = reqDict2(py[:2])
		rYm, ym = reqDict1(py[2:])
		if sm && ym {
			return dict[rSm] + rYm
		}
	}

	if len(py) == 2 {
		// 1 + 1
		sm, ym = false, false
		rSm, sm = reqDict1(py[:1])
		rYm, ym = reqDict1(py[1:])
		if sm && ym {
			return rSm + rYm
		}
	}

	// 特殊处理
	if py == "a" {
		return "aa"
	}
	if py == "ang" {
		return "ah"
	}
	if py == "e" {
		return "ee"
	}
	if py == "eng" {
		return "eg"
	}
	if py == "o" {
		return "oo"
	}

	// 其余返回本身
	return py
}

//dict := ["iu", "ei", "uan", "ue", "ve", "un", "uo", "ie", "iong", "ong", "ai", "en", "eng", "ang", "an", "uai", "ing", "uang", "iang", "ou", "ua", "ia", "ao", "ui", "in", "iao", "ian", "sh", "ch", "zh", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "a", "s", "d", "f", "g", "h", "j", "k", "l", "z", "x", "c", "v", "b", "n", "m", "ai", "iu"]

func reqDict1(char string) (string, bool) {
	dict1 := []string{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "a", "s", "d", "f", "g", "h", "j", "k", "l", "z", "x", "c", "v", "b", "n", "m"}

	if len(char) != 1 {
		return "", false
	}

	flag := false
	result := ""

	for _, c := range dict1 {
		if c == char {
			result = c
			flag = true
			break
		}
	}
	return result, flag
}

func reqDict2(char string) (string, bool) {
	dict2 := []string{"iu", "ei", "ue", "ve", "un", "uo", "ie", "ai", "en", "an", "ou", "ua", "ia", "ao", "ui", "in", "sh", "ch", "zh", "ai", "iu"}

	if len(char) != 2 {
		return "", false
	}

	flag := false
	result := ""

	for _, c := range dict2 {
		if c == char {
			result = c
			flag = true
			break
		}
	}
	return result, flag
}

func reqDict3(char string) (string, bool) {
	dict3 := []string{"uan", "ong", "eng", "ang", "uai", "ing", "iao", "ian"}

	if len(char) != 3 {
		return "", false
	}

	flag := false
	result := ""

	for _, c := range dict3 {
		if c == char {
			result = c
			flag = true
			break
		}
	}
	return result, flag
}

func reqDict4(char string) (string, bool) {
	dict4 := []string{"iong", "uang", "iang"}

	if len(char) != 4 {
		return "", false
	}

	flag := false
	result := ""

	for _, c := range dict4 {
		if c == char {
			result = c
			flag = true
			break
		}
	}
	return result, flag
}

func initDict() map[string]string {
	dict := make(map[string]string)
	dict["iu"] = "q"
	dict["ei"] = "w"
	dict["uan"] = "r"
	dict["ue"] = "t"
	dict["ve"] = "t"
	dict["un"] = "y"
	dict["sh"] = "u"
	dict["ch"] = "i"
	dict["uo"] = "o"
	dict["ie"] = "p"
	dict["iong"] = "s"
	dict["ong"] = "s"
	dict["ai"] = "d"
	dict["en"] = "f"
	dict["eng"] = "g"
	dict["ang"] = "h"
	dict["an"] = "j"
	dict["uai"] = "k"
	dict["ing"] = "k"
	dict["uang"] = "l"
	dict["iang"] = "l"
	dict["ou"] = "z"
	dict["ia"] = "x"
	dict["ua"] = "x"
	dict["ao"] = "c"
	dict["zh"] = "v"
	dict["ui"] = "v"
	dict["in"] = "b"
	dict["iao"] = "n"
	dict["ian"] = "m"
	return dict
}
