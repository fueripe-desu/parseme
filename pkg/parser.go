package parseme

import (
	"unicode"
)

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

func parseBytes(bytes *[]byte) *[]preToken {
	buffer := &parseBuffer{}
	tokens := make([]preToken, 0)

	// Indicates if char is inside a tag
	isTag := false

	// Indicates if next element is a property value
	isProperty := false

	for i := 0; i < len(*bytes); i++ {
		current := (*bytes)[i]

		if unicode.IsControl(rune(current)) {
			continue
		}

		if isTag {
			// If inside a tag
			lastIndex := len(tokens) - 1
			last := tokens[lastIndex]

			if current == byte('/') {
				// Replace last opening token with slash opening
				if last.tokenType == tagStart {
					token := &preToken{tokenType: slashTagStart, value: "</"}
					tokens = append(tokens[:lastIndex], *token)
				} else {
					// Handle '/' in the middle of the token
					// TODO: Add warning in error pool
				}
				continue
			}

			if current == byte(' ') {
				continue
			}

			if current == byte('=') {
				isProperty = true
				continue
			}

			if current == byte('>') {
				isTag = false
				token := &preToken{tokenType: tagEnd, value: ">"}
				tokens = append(tokens, *token)
				continue
			}

			if unicode.IsGraphic(rune(current)) {
				i = parseTagString(i, bytes, buffer)
				value := buffer.get()
				var token *preToken

				if last.tokenType == tagStart || last.tokenType == slashTagStart {
					token = &preToken{tokenType: tagName, value: string(*value)}
				} else if last.tokenType == tagName {
					token = &preToken{tokenType: propertyName, value: string(*value)}
				} else if last.tokenType == propertyName {
					if isProperty {
						token = &preToken{tokenType: propertyValue, value: string(*value)}
						isProperty = false
					} else {
						token = &preToken{tokenType: propertyValue, value: string(*value)}
					}
				}
				tokens = append(tokens, *token)
				buffer.clear()
				continue
			}
		} else {
			// If outside a tag

			if current == byte('<') {
				isTag = true
				token := &preToken{tokenType: tagStart, value: "<"}
				tokens = append(tokens, *token)
			} else if current == byte(' ') {
				continue
			} else {
				i = parseContent(i, bytes, buffer)
				value := buffer.get()
				token := &preToken{tokenType: content, value: string(*value)}
				tokens = append(tokens, *token)
				buffer.clear()
			}
		}
	}

	return &tokens
}

func removeControl(data *[]byte) *[]byte {
	result := make([]byte, 0, len(*data))
	for _, b := range *data {
		if !unicode.IsControl(rune(b)) {
			result = append(result, b)
		}
	}
	return &result
}

func parseContent(start int, bytes *[]byte, buffer *parseBuffer) int {
	i := start
	for ; i < len(*bytes); i++ {
		current := (*bytes)[i]

		if current == byte('<') {
			break
		}

		buffer.append(current)
	}

	return i - 1
}

func parseTagString(start int, bytes *[]byte, buffer *parseBuffer) int {
	i := start
	for ; i < len(*bytes); i++ {
		current := (*bytes)[i]

		if current == byte(' ') {
			break
		}

		if current == byte('>') {
			break
		}

		if current == byte('=') {
			break
		}

		buffer.append(current)
	}

	return i - 1
}
