package customSocket

import (
	"sort"
	"testing"
)

func (password *password) Len() int {
	return PasswordLength
}

func (password *password) Less(i, j int) bool {
	return password[i] < password[j]
}

func (password *password) Swap(i, j int) {
	password[i], password[j] = password[j], password[i]
}

func TestRandPassword(t *testing.T) {
	password := RandPassword()
	t.Log(password)
	bsPassword, err := parsePassword(password)
	if err != nil{
		t.Error(err)
	}
	sort.Sort(bsPassword)
	for i := 0; i < PasswordLength; i++ {
		if bsPassword[i] != byte(i) {
			t.Error("不能出现重复的byte位,必须由0-255组成")
		}
	}

}