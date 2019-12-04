package customSocket

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const PasswordLength  = 256
type password [PasswordLength]byte

func init()  {
	//更新随机种子
	rand.Seed(time.Now().Unix())
}

//base64编码将密码转换为字符串
func (password *password) String() string {
	return base64.StdEncoding.EncodeToString(password[:])
}

func parsePassword(passwordString string) (*password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(passwordString))
	if err != nil || len(bs) != PasswordLength {
		err = fmt.Errorf("密码长度不符合，当前的长度为 %d", len(bs))
		return nil, err
	}
	password := password{}
	copy(password[:], bs)
	bs = nil
	return &password, nil
}

func RandPassword() string {
	//随机生成由0-255组成的byte数组
	intArr := rand.Perm(PasswordLength)
	password := &password{}
	for i, v := range intArr {
		password[i] = byte(v)
		if i == v {
			// 确保不会出现一个byte位重复
			return RandPassword()
		}
	}
	return password.String()
}