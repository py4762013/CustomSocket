package customSocket

import (
	"crypto/rand"
	"reflect"
	"testing"
)

const (
	MB = 1024 * 1024
)

func TestCipher(t *testing.T)  {
	password := RandPassword()
	t.Log(password)
	p, _ := parsePassword(password)
	cipher := NewCipher(p)
	// 原数据
	org := make([]byte, PasswordLength)
	for i := 0; i < PasswordLength; i++ {
		org[i] = byte(i)
	}
	// 复制一份原数据到tmp
	tmp := make([]byte, PasswordLength)
	copy(tmp, org)
	t.Log(tmp)
	//加密tmp
	cipher.encode(tmp)
	t.Log(tmp)
	//解密tmp
	cipher.decode(tmp)
	t.Log(tmp)
	if !reflect.DeepEqual(org, tmp) {
		t.Error("编码解码后数据不一致")
	}
}

func BenchmarkEncode(b *testing.B)  {
	password := RandPassword()
	p, _ := parsePassword(password)
	cipher := NewCipher(p)
	bs := make([]byte, MB)
	b.ResetTimer()
	rand.Read(bs)
	cipher.encode(bs)
}

func BenchmarkDecode(b *testing.B)  {
	password := RandPassword()
	p, _ := parsePassword(password)
	cipher := NewCipher(p)
	bs := make([]byte, MB)
	b.ResetTimer()
	rand.Read(bs)
	cipher.decode(bs)
}