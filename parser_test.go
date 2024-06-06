package parseme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseBytes(t *testing.T) {
	testcasesFile := []struct {
		name     string
		value    string
		expected []preToken
	}{
		{
			"regular tag",
			"<html></html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"leading spaces outside tag",
			"      <html></html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"trailing spaces outside tag",
			"<html></html>     ",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"nested tags",
			"<html><p></p></html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"leading spaces in nested tag",
			"<html>     <p></p></html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"trailing spaces in nested tag",
			"<html><p></p>      </html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"nested tags with content",
			"<html><p>This is a paragraph</p></html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: content, value: "This is a paragraph"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"tag with content",
			"<p>This is a paragraph</p>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: content, value: "This is a paragraph"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"with control chars",
			"<html>\n\t<p>This is a paragraph</p>\n</html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: content, value: "This is a paragraph"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"leading spaces between tag and content",
			"<html><p>        This is a paragraph</p></html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: content, value: "This is a paragraph"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"trailing spaces between tag and content",
			"<html><p>This is a paragraph</p>        </html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: content, value: "This is a paragraph"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"spaces between words in content",
			"<html><p>This is a      paragraph</p></html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: content, value: "This is a      paragraph"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "p"},
				{tokenType: tagEnd, value: ">"},
				{tokenType: slashTagStart, value: "</"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"content outside of tag",
			"This is content outside of tag",
			[]preToken{
				{tokenType: content, value: "This is content outside of tag"},
			},
		},
		{
			"non conflicting unescaped special characters",
			"This is a ; &",
			[]preToken{
				{tokenType: content, value: "This is a ; &"},
			},
		},
		{
			"unicode content",
			"This is a unicode character: ♥",
			[]preToken{
				{tokenType: content, value: "This is a unicode character: ♥"},
			},
		},
		{
			"special characters in tag name",
			"<ht&l>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "ht&l"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"leading spaces in tag name",
			"<      html>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"trailing spaces in tag name",
			"<html      >",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"tag with value property",
			"<html lang=\"en-US\">",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: propertyName, value: "lang"},
				{tokenType: propertyValue, value: "\"en-US\""},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"tag with bool property",
			"<html activated>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: propertyName, value: "activated"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"bool property followed by value property",
			"<html activated lang=\"en-US\">",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: propertyName, value: "activated"},
				{tokenType: propertyName, value: "lang"},
				{tokenType: propertyValue, value: "\"en-US\""},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"value property followed by bool property",
			"<html lang=\"en-US\" activated>",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "html"},
				{tokenType: propertyName, value: "lang"},
				{tokenType: propertyValue, value: "\"en-US\""},
				{tokenType: propertyName, value: "activated"},
				{tokenType: tagEnd, value: ">"},
			},
		},
		{
			"void tag",
			"<img src=\"some.random.source\" />",
			[]preToken{
				{tokenType: tagStart, value: "<"},
				{tokenType: tagName, value: "img"},
				{tokenType: propertyName, value: "src"},
				{tokenType: propertyValue, value: "\"some.random.source\""},
				{tokenType: slashTagEnd, value: "/>"},
			},
		},
	}

	for _, tc := range testcasesFile {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			bytes := []byte(tc.value)
			tokens := parseBytes(&bytes)
			assert.Equal(*tokens, tc.expected)
		})
	}

}

func Test_removeControl(t *testing.T) {
	testcases := []struct {
		name     string
		value    string
		expected string
	}{
		{
			"string with control chars",
			"This is some \n\t random text.\n",
			"This is some  random text.",
		},
		{
			"string without control char",
			"This is some random text.",
			"This is some random text.",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			bytes := []byte(tc.value)
			newBytes := removeControl(&bytes)
			expected := []byte(tc.expected)
			assert.Equal(*newBytes, expected)
		})
	}
}

func Test_parseContent(t *testing.T) {
	testcases := []struct {
		name     string
		start    int
		value    string
		expected string
	}{
		{
			"normal text",
			0,
			"Some random text",
			"Some random text",
		},
		{
			"start greater than zero",
			2,
			"Some random text",
			"me random text",
		},
		{
			"value with stop char",
			0,
			"Some random< text",
			"Some random",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			buffer := &parseBuffer{}

			t.Cleanup(func() {
				buffer.clear()
			})

			byteValue := []byte(tc.value)
			parseContent(tc.start, &byteValue, buffer)
			byteExpected := []byte(tc.expected)
			assert.Equal(*(buffer.get()), byteExpected)
		})
	}
}

func Test_parseTagString(t *testing.T) {
	testcases := []struct {
		name     string
		start    int
		value    string
		expected string
	}{
		{
			"normal text",
			0,
			"Some",
			"Some",
		},
		{
			"start greater than zero",
			2,
			"Some",
			"me",
		},
		{
			"value with space stop char",
			0,
			"Some random",
			"Some",
		},
		{
			"value with tag end stop char",
			0,
			"Some>random",
			"Some",
		},
		{
			"value with equal stop char",
			0,
			"Some=random",
			"Some",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			buffer := &parseBuffer{}

			t.Cleanup(func() {
				buffer.clear()
			})

			byteValue := []byte(tc.value)
			parseTagString(tc.start, &byteValue, buffer)
			byteExpected := []byte(tc.expected)
			assert.Equal(*(buffer.get()), byteExpected)
		})
	}
}
