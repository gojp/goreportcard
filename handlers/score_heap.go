package handlers

type scoreItem struct {
	Repo  string  `json:"repo"`
	Score float64 `json:"score"`
	Files int     `json:"files"`
}

// An scoreHeap is a min-heap of ints.
type scoreHeap []scoreItem

func (h scoreHeap) Len() int           { return len(h) }
func (h scoreHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h scoreHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *scoreHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(scoreItem))
}

func (h *scoreHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
