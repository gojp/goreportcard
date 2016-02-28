package handlers

type recentItem struct {
	Repo string `json:"repo"`
}

type recentHeap []recentItem

func (h recentHeap) Len() int           { return len(h) }
func (h recentHeap) Less(i, j int) bool { return true }
func (h recentHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *recentHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(recentItem))
}

func (h *recentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
