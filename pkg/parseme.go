package parseme

type HtmlParser struct {
	filepath string
}

func NewHtmlParser(filepath string) *HtmlParser {
	return &HtmlParser{filepath: filepath}
}

func (p *HtmlParser) Parse() (*[]preToken, error) {
	bytes, readErr := fetchFileContents(p.filepath)

	if readErr != nil {
		return nil, readErr
	}

	preTokens := parseBytes(bytes)

	return preTokens, nil
}
