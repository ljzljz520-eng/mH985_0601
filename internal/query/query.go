package query

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"trainingdesk/internal/model"
)

type Expression struct {
	Field    string
	Operator string
	Value    string
}

type Request struct {
	Expressions []Expression
	SortField   string
	Descending  bool
	Offset      int
	Limit       int
}

func Parse(input string) (Request, error) {
	request := Request{Expressions: make([]Expression, 0), Limit: 100}
	for _, token := range strings.Fields(input) {
		if strings.HasPrefix(token, "sort=") {
			field := strings.TrimPrefix(token, "sort=")
			if strings.HasPrefix(field, "-") {
				request.Descending = true
				field = strings.TrimPrefix(field, "-")
			}
			if !validField(field) {
				return Request{}, fmt.Errorf("invalid sort field %q", field)
			}
			request.SortField = field
			continue
		}
		if strings.HasPrefix(token, "offset=") {
			value, err := strconv.Atoi(strings.TrimPrefix(token, "offset="))
			if err != nil || value < 0 {
				return Request{}, errors.New("offset must be non-negative")
			}
			request.Offset = value
			continue
		}
		if strings.HasPrefix(token, "limit=") {
			value, err := strconv.Atoi(strings.TrimPrefix(token, "limit="))
			if err != nil || value < 1 || value > 500 {
				return Request{}, errors.New("limit must be between 1 and 500")
			}
			request.Limit = value
			continue
		}
		expression, err := parseExpression(token)
		if err != nil {
			return Request{}, err
		}
		request.Expressions = append(request.Expressions, expression)
	}
	return request, nil
}

func parseExpression(token string) (Expression, error) {
	for _, operator := range []string{"!=", ">=", "<=", "=", ">", "<", "~"} {
		if index := strings.Index(token, operator); index > 0 {
			field := token[:index]
			value := token[index+len(operator):]
			if !validField(field) || value == "" {
				return Expression{}, fmt.Errorf("invalid expression %q", token)
			}
			return Expression{Field: field, Operator: operator, Value: value}, nil
		}
	}
	return Expression{}, fmt.Errorf("expression must use an operator: %q", token)
}

func validField(field string) bool {
	switch field {
	case "id", "store", "title", "category", "status", "owner", "reviewer", "sort_key", "version":
		return true
	default:
		return false
	}
}

func (r Request) Match(record model.Record) bool {
	for _, expression := range r.Expressions {
		if !matchExpression(expression, record) {
			return false
		}
	}
	return true
}

func matchExpression(expression Expression, record model.Record) bool {
	actual := fieldValue(expression.Field, record)
	if expression.Operator == "~" {
		return strings.Contains(strings.ToLower(actual), strings.ToLower(expression.Value))
	}
	if expression.Operator == "=" {
		return actual == expression.Value
	}
	if expression.Operator == "!=" {
		return actual != expression.Value
	}
	if expression.Field == "sort_key" || expression.Field == "version" {
		left, leftErr := strconv.Atoi(actual)
		right, rightErr := strconv.Atoi(expression.Value)
		if leftErr != nil || rightErr != nil {
			return false
		}
		if expression.Operator == ">" {
			return left > right
		}
		if expression.Operator == "<" {
			return left < right
		}
		if expression.Operator == ">=" {
			return left >= right
		}
		return left <= right
	}
	return false
}

func fieldValue(field string, record model.Record) string {
	switch field {
	case "id":
		return record.ID
	case "store":
		return record.StoreID
	case "title":
		return record.Title
	case "category":
		return record.Category
	case "status":
		return string(record.Status)
	case "owner":
		return record.Owner
	case "reviewer":
		return record.Reviewer
	case "sort_key":
		return strconv.Itoa(record.SortKey)
	case "version":
		return strconv.Itoa(record.Version)
	default:
		return ""
	}
}
