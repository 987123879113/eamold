// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

type Gf8dm7puvCardProfile struct {
	GameType int64
	Cardid   string
	Name     string
	Color    int64
	Recovery int64
	Styles   int64
	Hidden   int64
	Expired  int64
}

type Gf8dm7puvDemomusic struct {
	GameType int64
	Musicid  int64
}

type Gf8dm7puvFavorite struct {
	GameType int64
	Musicid  int64
	Count    int64
}

type Gf8dm7puvMessage struct {
	ID       int64
	GameType int64
	Enabled  int64
	Message  string
}

type Gf8dm7puvPuzzle struct {
	GameType int64
	Cardid   string
	Number   int64
	Flags    int64
	Out      int64
}

type Gf8dm7puvScore struct {
	ID       int64
	Cardid   string
	GameType int64
	Musicid  int64
	Musicnum int64
	Seq      int64
	Score    int64
	Flags    int64
	Clear    int64
	Skill    int64
	Combo    int64
	Encore   int64
	Extra    int64
}

type Gf8dm7puvShop struct {
	GameType int64
	Pref     int64
	Name     string
	Points   int64
}

type Gf8dm7puvShopMachine struct {
	GameType int64
	Pcbid    string
	Pref     int64
	Name     string
}
