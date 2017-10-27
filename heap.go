package compressor

type node struct {
	char      byte
	frequency int
	internal  bool
	left      *node
	right     *node
}

// Implmentation of minimum heap on the node's frequency value
type minHeap []*node

func (h minHeap) Len() int {
	return len(h)
}

func (h minHeap) Less(i, j int) bool {
	return h[i].frequency < h[j].frequency
}

func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(*node))
}

func (h *minHeap) Pop() interface{} {
	currentHeap := *h
	length := len(currentHeap)

	// update the heap
	*h = currentHeap[0 : length-1]

	val := currentHeap[length-1]

	return val
}
