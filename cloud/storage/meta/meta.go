package meta

import "fmt"

type Store struct {
	images map[string]Image
}

func NewStore() *Store {
	return &Store{images: make(map[string]Image)}
}

func (s *Store) Add(img Image) {
	s.images[img.ID] = img
	fmt.Println("Meta: added image", img.ID)
}

func (s *Store) Get(id string) (Image, bool) {
	img, ok := s.images[id]
	return img, ok
}
