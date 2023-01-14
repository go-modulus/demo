package domain

type Edge[N any] struct {
	Cursor string
	Node   N
}

type Page[N any] struct {
	Edges           []Edge[N]
	HasPreviousPage bool
	HasNextPage     bool
	StartCursor     *string
	EndCursor       *string
}
