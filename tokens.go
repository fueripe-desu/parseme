package parseme

type preTokenType int

const (
	tagStart = iota
	slashTagStart
	slashTagEnd
	tagName
	tagEnd
	content
	propertyName
	propertyValue
)

type preToken struct {
	tokenType preTokenType
	value     string
}
