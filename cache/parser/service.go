package parser

// NewParserService to set up the entity Cache service
func NewParserService() IParser {
	return &SParser{}
}
