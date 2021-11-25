package github

import (
	"fmt"
	"strings"
)

type Event struct {
	title, file, level, message   string
	col, endColumn, line, endLine uint64
}

func NewEvent(level, message string) *Event {
	return &Event{
		level:   level,
		message: message,
	}
}

func (e *Event) WithTitle(title string) *Event {
	e.title = title

	return e
}

func (e *Event) String() string {
	var labels []string
	if e.title != "" {
		labels = append(labels, fmt.Sprintf("title=%s", e.title))
	}
	if e.file != "" {
		labels = append(labels, fmt.Sprintf("file=%s", e.file))
	}
	if e.col != 0 {
		labels = append(labels, fmt.Sprintf("col=%d", e.col))
	}
	if e.endColumn != 0 {
		labels = append(labels, fmt.Sprintf("endColumn=%d", e.endColumn))
	}
	if e.line != 0 {
		labels = append(labels, fmt.Sprintf("line=%d", e.line))
	}
	if e.endLine != 0 {
		labels = append(labels, fmt.Sprintf("endLine=%d", e.endLine))
	}

	return fmt.Sprintf("::%s %s::%s", e.level, strings.Join(labels, ","), e.message)
}
