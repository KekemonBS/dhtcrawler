package models

//Share stores data detected bittorrent share
type Share struct {
	Name       string
	Size       int
	FileTree   string
	MagnetLink string
}

//SharesPage stores resulting page of Shares
type SharesPage struct {
	Total   int
	Results []Share
}
