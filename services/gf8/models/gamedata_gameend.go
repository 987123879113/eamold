package models

import "encoding/xml"

/*
	beginner:
	<gamedata method="gameend"><favorite nr="3" musicid="733?733?713?"/><shoprank point="28" shopname="スヌフ" pref="13"/></gamedata>

	normal:
	<gamedata method="gameend">
		<favorite nr="3" musicid="500?601?538?"/>
		<shoprank point="146" shopname="スヌフ" pref="13"/>
		<rank stage_nr="5" musicid="500?601?538?717?703?">
			<player no="0" seq="1?1?1?1?1?" score="1811750?1989710?2358575?3509825?4831700?"/>
		</rank>
	</gamedata>

	course:
	<gamedata method="gameend"><shoprank point="237" shopname="スヌフ" pref="13"/></gamedata>
*/

type Request_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Favorite *struct {
		Count    int    `xml:"nr,attr"`
		MusicIds string `xml:"musicid,attr"`
	} `xml:"favorite"`

	ShopRank *struct {
		Point    int    `xml:"point,attr"`
		ShopName string `xml:"shopname,attr"`
		Pref     int    `xml:"pref,attr"`
	} `xml:"shoprank"`

	Rank *struct {
		StageCount int    `xml:"stage_nr,attr"`
		MusicIDs   string `xml:"musicid,attr"`

		Player []struct {
			Number int    `xml:"no,attr"`
			Seqs   string `xml:"seq,attr"`
			Scores string `xml:"score,attr"`
		} `xml:"player"`
	} `xml:"rank"`
}

type Response_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Rank Response_GameData_GameEnd_RankData `xml:"rank"`
}

type Response_GameData_GameEnd_RankData_Rank struct {
	Number int    `xml:"no,attr"`
	Order  string `xml:"order,attr"`
}

type Response_GameData_GameEnd_RankData struct {
	All string `xml:"all,attr"`

	Rank []Response_GameData_GameEnd_RankData_Rank `xml:"rank"`
}
