package parser

import (
"testing"
)

func TestParseLimit(t *testing.T) {
q, err := Run("SELECT * FROM /tmp LIMIT 10")
if err != nil {
t.Fatalf("Expected no error, got %v", err)
}
if q.Limit != 10 {
t.Fatalf("Expected limit 10, got %d", q.Limit)
}
}
