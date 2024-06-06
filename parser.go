package parseme

import (
	"unicode"
)

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
					buffer.clear()
					buffer.append(byte('/'))
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
				value := *(buffer.get())

				var token *preToken

				if len(value) == 1 && value[0] == byte('/') {
					token = &preToken{tokenType: slashTagEnd, value: "/>"}
				} else {
					token = &preToken{tokenType: tagEnd, value: ">"}
				}
				tokens = append(tokens, *token)
				continue
			}

			if unicode.IsGraphic(rune(current)) {
				i = parseTagString(i, bytes, buffer)
				value := buffer.get()
				var token *preToken

				if last.tokenType == tagStart || last.tokenType == slashTagStart {
					token = &preToken{tokenType: tagName, value: string(*value)}
				} else if last.tokenType == tagName || last.tokenType == propertyValue {
					token = &preToken{tokenType: propertyName, value: string(*value)}
				} else if last.tokenType == propertyName {
					if isProperty {
						token = &preToken{tokenType: propertyValue, value: string(*value)}
						isProperty = false
					} else {
						token = &preToken{tokenType: propertyName, value: string(*value)}
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
