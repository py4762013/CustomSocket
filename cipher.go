package customSocket

type cipher struct {
	encodePassword *password
	decodePassword *password
}

//加密原数据
func (cipher *cipher) encode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.encodePassword[v]
	}
}

//解密加密后的数据
func (cipher *cipher) decode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.decodePassword[v]
	}
}

func NewCipher(encodePassword *password) *cipher {
	decodePassword := &password{}
	for i, v := range encodePassword{
		encodePassword[i] = v
		decodePassword[v] = byte(i)
	}
	return &cipher{
		encodePassword: encodePassword,
		decodePassword: decodePassword,
	}
}

