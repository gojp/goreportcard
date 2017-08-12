package handlers

type scoreItem struct {
	Repo  string  `json:"repo"`
	Score float64 `json:"score"`
	Files int     `json:"files"`
}

// An ScoreHeap is a min-heap of ints.
type ScoreHeap []scoreItem

func (h ScoreHeap) Len() int           { return len(h) }
func (h ScoreHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h ScoreHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push onto the heap
func (h *ScoreHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(scoreItem))
}

// Pop item off of the heap
func (h *ScoreHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
