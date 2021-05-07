package similar_searcher

type SimilarSearcher struct {
}

func NewSimilarSearcher() *SimilarSearcher {
	return &SimilarSearcher{}
}

func (ss *SimilarSearcher) Search(q string) (string, error) {
	// TBD
	return "", nil
}
