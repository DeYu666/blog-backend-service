package shuangpin

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	nums := 10

	words, _ := GetWords([]string{"IT"}, nums)

	words = append(words, "xaml文件1267890")

	for _, word := range words {
		fmt.Println(word)
		pY, _ := GetPingYin(word)
		fmt.Println(pY)

		sp := PinYinTurnShuangPin(pY)

		sp = CheckResult(word, pY, sp)
		fmt.Println(sp)
		fmt.Println("-----------------")
	}
}

func TestGetWords(t *testing.T) {
	nums := 70
	words, err := GetWords([]string{}, nums)
	fmt.Println(words)
	fmt.Println(err)
}

func TestGetPingYinByApi(t *testing.T) {
	word := "xml文件1267890"

	result, _ := GetPingYin(word)
	fmt.Println(result)
}

func TestTurnShuangPin(t *testing.T) {
	word := "xml文件1267890"

	result, _ := GetPingYin(word)

	fmt.Println(result)

	for _, r := range result {
		fmt.Println(r)
		fmt.Println(turnShuangPin(r))
	}

}
