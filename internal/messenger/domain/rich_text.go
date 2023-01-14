package domain

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

const (
	PlainRichTextType = "plain"
)

type RichText struct {
	Text  sql.NullString
	Parts RichTextParts
}

type RichTextParts []RichTextPart

func (p *RichTextParts) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	return json.Unmarshal(value.([]byte), p)
}

func (p RichTextParts) MarshalJSON() ([]byte, error) {
	rawParts := make([]map[string]any, len(p))

	for i, part := range p {
		rawPart := map[string]any{
			"type": part.Type(),
			"text": part.Text(),
		}
		rawParts[i] = rawPart
	}

	return json.Marshal(rawParts)
}

func (p *RichTextParts) UnmarshalJSON(input []byte) error {
	var rawParts []map[string]any
	if err := json.Unmarshal(input, &rawParts); err != nil {
		return err
	}

	for _, rawPart := range rawParts {
		switch rawPart["type"] {
		case PlainRichTextType:
			plainText := PlainRichText{
				text: rawPart["text"].(string),
			}
			*p = append(*p, plainText)
		default:
			return fmt.Errorf("unknown rich text part type: %s", rawPart["type"])
		}
	}

	return nil
}

type RichTextPart interface {
	Text() string
	Type() string
}

type PlainRichText struct {
	text string
}

func (p PlainRichText) Text() string { return p.text }

func (p PlainRichText) Type() string { return PlainRichTextType }

type RichTextParser struct {
}

func NewRichTextParser() *RichTextParser {
	return &RichTextParser{}
}

func (p *RichTextParser) Parse(text string) (RichText, error) {
	if text == "" {
		return RichText{}, nil
	}

	return RichText{
		Text:  sql.NullString{String: text, Valid: true},
		Parts: RichTextParts{PlainRichText{text: text}},
	}, nil
}
