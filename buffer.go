package parseme

type parseBuffer struct {
	value []byte
}

func (b *parseBuffer) append(char byte) {
	b.value = append(b.value, char)
}

func (b *parseBuffer) appendAll(chars *[]byte) {
	b.value = append(b.value, (*chars)...)
}

func (b *parseBuffer) get() *[]byte {
	return &b.value
}

func (b *parseBuffer) clear() {
	b.value = b.value[:0]
}
