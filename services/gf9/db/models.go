// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

type Gf9dm8CardProfile struct {
	GameType int64
	Cardid   string
	Name     string
	Recovery int64
	Styles   int64
	Hidden   int64
	Skill    int64
	Expired  int64
}

type Gf9dm8Demomusic struct {
	GameType int64
	Musicid  int64
}

type Gf9dm8Favorite struct {
	GameType int64
	Musicid  int64
	Count    int64
}

type Gf9dm8Message struct {
	ID       int64
	GameType int64
	Enabled  int64
	Message  string
}

type Gf9dm8Puzzle struct {
	GameType int64
	Cardid   string
	Number   int64
	Flags    int64
	Out      int64
}

type Gf9dm8Score struct {
	ID       int64
	Cardid   string
	GameType int64
	MusicNum int64
	SeqMode  int64
	Score    int64
	Flags    int64
	Clear    int64
	Skill    int64
	Combo    int64
	Encore   int64
	Extra    int64
}

type Gf9dm8Shop struct {
	GameType int64
	Pref     int64
	Name     string
	Points   int64
}

type Gf9dm8ShopMachine struct {
	GameType int64
	Pcbid    string
	Pref     int64
	Name     string
}
