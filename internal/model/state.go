package model

// ViewMode represents the current view mode
type ViewMode int

const (
	NormalView ViewMode = iota
	ZoomView
	FilterView
)
