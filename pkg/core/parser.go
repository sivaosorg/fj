package core

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sivaosorg/fj/pkg/errors"
)

// Parser implements the JSONParser interface with enhanced error handling,
// event publishing, and configuration management.
type Parser struct {
	eventPublisher EventPublisher
	validator     ValidationEngine
}

// NewParser creates a new parser instance with default configuration
func NewParser() JSONParser {
	return &Parser{}
}

// NewParserWithConfig creates a new parser with custom configuration  
func NewParserWithConfig() JSONParser {
	return &Parser{}
}

// SetEventPublisher sets the event publisher for this parser
func (p *Parser) SetEventPublisher(publisher EventPublisher) {
	p.eventPublisher = publisher
}

// SetValidator sets the validation engine for this parser
func (p *Parser) SetValidator(validator ValidationEngine) {
	p.validator = validator
}

// Parse parses a JSON string and returns a Context representing the parsed data
func (p *Parser) Parse(json string) Context {
	startTime := time.Now()
	
	// Publish parse start event if event publishing is enabled
	if p.eventPublisher != nil {
		p.eventPublisher.PublishParseStart(json)
	}
	
	// Validate input if validation is enabled
	if p.validator != nil {
		if err := p.validator.ValidateJSON(json); err != nil {
			duration := time.Since(startTime)
			fjError := errors.WrapError(errors.JSONParseError, "JSON validation failed", err).
				WithInput(json).
				WithContext("validation_enabled", true)
			
			if p.eventPublisher != nil {
				p.eventPublisher.PublishParseError(json, fjError, duration)
			}
			
			return NewContextWithError(fjError)
		}
	}
	
	// Perform the actual parsing using existing logic
	// For now, we'll delegate to the existing Parse function from the root package
	result := p.parseInternal(json)
	
	duration := time.Since(startTime)
	
	// Publish completion event
	if p.eventPublisher != nil {
		if result.HasError() {
			p.eventPublisher.PublishParseError(json, result.Error(), duration)
		} else {
			p.eventPublisher.PublishParseComplete(result, duration)
		}
	}
	
	return result
}

// ParseBytes parses JSON from a byte slice
func (p *Parser) ParseBytes(data []byte) Context {
	if p.validator != nil {
		if err := p.validator.ValidateJSONBytes(data); err != nil {
			fjError := errors.WrapError(errors.JSONParseError, "JSON bytes validation failed", err).
				WithInput(string(data)).
				WithContext("validation_enabled", true)
			return NewContextWithError(fjError)
		}
	}
	
	return p.Parse(string(data))
}

// ParseReader parses JSON from an io.Reader
func (p *Parser) ParseReader(reader io.Reader) Context {
	if reader == nil {
		return NewContextWithError(errors.NewIOError("reader cannot be nil", nil))
	}
	
	data, err := io.ReadAll(reader)
	if err != nil {
		return NewContextWithError(errors.NewIOError("failed to read from reader", err))
	}
	
	return p.ParseBytes(data)
}

// ParseFile parses JSON from a file path
func (p *Parser) ParseFile(filepath string) Context {
	if filepath == "" {
		return NewContextWithError(errors.NewIOError("file path cannot be empty", nil))
	}
	
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewContextWithError(errors.NewIOError("file not found: "+filepath, err).
				WithContext("file_path", filepath))
		}
		if os.IsPermission(err) {
			return NewContextWithError(errors.NewIOError("permission denied: "+filepath, err).
				WithContext("file_path", filepath))
		}
		return NewContextWithError(errors.NewIOError("failed to open file: "+filepath, err).
			WithContext("file_path", filepath))
	}
	defer file.Close()
	
	return p.ParseReader(file)
}

// IsValid checks if the provided JSON string is valid
func (p *Parser) IsValid(json string) bool {
	if p.validator != nil {
		return p.validator.ValidateJSON(json) == nil
	}
	
	// Fallback to basic parsing check
	result := p.parseInternal(json)
	return !result.HasError()
}

// IsValidBytes checks if the provided JSON bytes are valid
func (p *Parser) IsValidBytes(data []byte) bool {
	if p.validator != nil {
		return p.validator.ValidateJSONBytes(data) == nil
	}
	
	return p.IsValid(string(data))
}

// parseInternal performs the actual parsing logic
// This is where we'll integrate with the existing parsing code from fj.go
func (p *Parser) parseInternal(json string) Context {
	// For now, this is a placeholder that creates a basic context
	// In the actual implementation, this would call the existing Parse function
	// from the root package and convert the result to our enhanced Context interface
	
	// Basic validation - check for empty input
	if json == "" {
		return NewContextWithError(errors.NewJSONParseError("empty JSON input", 0))
	}
	
	// Trim whitespace
	json = strings.TrimSpace(json)
	
	if json == "" {
		return NewContextWithError(errors.NewJSONParseError("JSON contains only whitespace", 0))
	}
	
	// Basic type detection
	switch json[0] {
	case '{':
		return NewContextFromExisting(6, json, p) // JSON object
	case '[':
		return NewContextFromExisting(6, json, p) // JSON array
	case '"':
		// Extract string value (simplified)
		if len(json) >= 2 && json[len(json)-1] == '"' {
			return NewContext(4, json, json[1:len(json)-1], 0, p) // String
		}
		return NewContextWithError(errors.NewJSONParseError("unterminated string", 0))
	case 't':
		if json == "true" {
			return NewContext(5, json, "true", 1, p) // True
		}
		return NewContextWithError(errors.NewJSONParseError("invalid boolean value", 0))
	case 'f':
		if json == "false" {
			return NewContext(1, json, "false", 0, p) // False
		}
		return NewContextWithError(errors.NewJSONParseError("invalid boolean value", 0))
	case 'n':
		if json == "null" {
			return NewContext(0, json, "", 0, p) // Null
		}
		return NewContextWithError(errors.NewJSONParseError("invalid null value", 0))
	default:
		// Try to parse as number
		if val, err := parseNumber(json); err == nil {
			return NewContext(2, json, "", val, p) // Number
		}
		return NewContextWithError(errors.NewJSONParseError("unexpected token: "+string(json[0]), 0))
	}
}

// parseNumber attempts to parse a string as a number
func parseNumber(s string) (float64, error) {
	// This is a simplified number parser
	// The actual implementation would use the more sophisticated number parsing
	// from the existing codebase
	return strconv.ParseFloat(s, 64)
}