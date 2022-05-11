package models

type file struct {
	Path   []interface{}
	Length int
}

//Share stores data detected bittorrent share
type Share struct {
	Name       string
	Size       int
	FileTree   string
	MagnetLink string
}

//type Share struct {
//	InfoHash string
//	Name     string
//	Files    []file
//	Length   int
//}
