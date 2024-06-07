package models

import "encoding/xml"

/*
	<gamedata method="rankset">
		<rank stage_nr="5" musicid="601?431?609?717?703?">
			<player seq="1?1?1?2?2?" score="1989710?2280230?2048210?5194925?7130900?" flags="0?0?0?0?0?" name="%42%44%47%4a"/>
		</rank>
	</gamedata>
*/

type Request_GameData_RankSet struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Rank struct {
		StageCount int    `xml:"stage_nr,attr"`
		MusicIDs   string `xml:"musicid,attr"`

		Player []struct {
			Seqs   string `xml:"seq,attr"`
			Scores string `xml:"score,attr"`
			Flags  string `xml:"flags,attr"`
			Name   string `xml:"name,attr"`
		} `xml:"player"`
	} `xml:"rank"`
}

type Response_GameData_RankSet struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`
}
