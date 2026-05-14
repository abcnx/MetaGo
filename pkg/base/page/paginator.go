// Package page provides base pagination utilities.
package page

// Paginator is a base paginator for internal use.
type Paginator struct {
	CurrentPage int
	PerPage     int
}

// Offset returns the offset for database queries.
func (p *Paginator) Offset() int {
	return (p.CurrentPage - 1) * p.PerPage
}
