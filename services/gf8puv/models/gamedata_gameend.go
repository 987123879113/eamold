package models

import "encoding/xml"

/*
	normal:
	<gamedata method="gameend">
		<favorite nr="2">
			<music id="750"/>
			<music id="104"/>
		</favorite>
		<shoprank point="3" shopname="ああかた" pref="13"/>
		<player no="0" cardid="8f0d9f3774e78f74" name="%41%42%43" color="0">
			<condition styles="0" hidden="0" recovery="0"/>
			<puzzle no="3" flags="8208" out="0"/>
			<judge perfect="91" great="0" good="0" poor="0" miss="105"/>
			<play time="162" nr="2" mode="1">
				<stage musicid="750" musicnum="113" seq="1" flags="32" score="223610" clear="1" skill="0" combo="91"/>
				<stage musicid="104" musicnum="7" seq="3" flags="0"/>
			</play>
		</player>
	</gamedata>

	course:
	<gamedata method="gameend">
		<shoprank point="488" shopname="ああかた" pref="13"/>
		<player no="0" cardid="8f0d9f3774e78f74" name="%41%42%43" color="0">
			<condition styles="0" hidden="1" recovery="1"/>
			<puzzle no="3" flags="47386" out="0"/>
			<judge perfect="0" great="0" good="0" poor="0" miss="0"/>
			<play time="469" nr="0" clear="1" mode="3"></play>
		</player>
	</gamedata>

	without card:
	<gamedata method="gameend">
		<favorite nr="3">
			<music id="626"/>
			<music id="750"/>
			<music id="737"/>
		</favorite>
		<shoprank point="126" shopname="ああかた" pref="13"/>
		<player no="0">
			<play nr="5" clear="1" mode="1">
				<stage musicid="626" musicnum="76" seq="1" flags="0" score="1298135" clear="1"/>
				<stage musicid="750" musicnum="113" seq="1" flags="0" score="1027505" clear="1"/>
				<stage musicid="737" musicnum="111" seq="1" flags="0" score="1869910" clear="1"/>
				<stage musicid="717" musicnum="92" seq="1" flags="0" score="3509825" clear="1"/>
				<stage musicid="703" musicnum="82" seq="1" flags="0" score="4831700" clear="1"/>
			</play>
		</player>
	</gamedata>
*/

type Request_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Favorite *struct {
		Count int `xml:"nr,attr"`
		Music []struct {
			Id int `xml:"id,attr"`
		} `xml:"music"`
	} `xml:"favorite"`

	ShopRank *struct {
		Point    int    `xml:"point,attr"`
		ShopName string `xml:"shopname,attr"`
		Pref     int    `xml:"pref,attr"`
	} `xml:"shoprank"`

	Players []struct {
		Number int    `xml:"no,attr"`
		CardId string `xml:"cardid,attr"`
		Name   string `xml:"name,attr"`
		Color  int    `xml:"color,attr"`

		Condition struct {
			Styles   int `xml:"styles,attr"`
			Hidden   int `xml:"hidden,attr"`
			Recovery int `xml:"recovery,attr"`
		} `xml:"condition"`

		Puzzle struct {
			Number int `xml:"number,attr"`
			Flags  int `xml:"flags,attr"`
			Out    int `xml:"out,attr"`
		} `xml:"puzzle"`

		Judge struct {
			Perfect int `xml:"perfect,attr"`
			Great   int `xml:"great,attr"`
			Good    int `xml:"good,attr"`
			Poor    int `xml:"poor,attr"`
			Miss    int `xml:"miss,attr"`
		} `xml:"judge"`

		Play struct {
			Time  int `xml:"time,attr"`
			Count int `xml:"nr,attr"`
			Clear int `xml:"clear,attr"`
			Mode  int `xml:"mode,attr"`

			Stages []struct {
				MusicId  int `xml:"musicid,attr"`
				MusicNum int `xml:"musicnum,attr"`
				Seq      int `xml:"seq,attr"`
				Flags    int `xml:"flags,attr"`
				Encore   int `xml:"encore,attr"`
				Extra    int `xml:"extra,attr"`
				Score    int `xml:"score,attr"`
				Clear    int `xml:"clear,attr"`
				Skill    int `xml:"skill,attr"`
				Combo    int `xml:"combo,attr"`
			} `xml:"stage"`
		} `xml:"play"`
	} `xml:"player"`
}

type Response_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Players []Response_GameData_GameEnd_Player `xml:"player"`

	ShopRank Response_GameData_GameEnd_ShopRank `xml:"shoprank"`
}

type Response_GameData_GameEnd_Player_Stage struct {
	All   int `xml:"all,attr"`
	Order int `xml:"order,attr"`
	Best  int `xml:"best,attr"`
}

type Response_GameData_GameEnd_Player struct {
	Number    int `xml:"no,attr"`
	Status    int `xml:"status,attr"`
	SkillPrev int `xml:"skill_prev,attr"`
	Skill     int `xml:"skill,attr"`

	Stages []Response_GameData_GameEnd_Player_Stage `xml:"stage"` // Does this go in the player response?
}

type Response_GameData_GameEnd_ShopRank struct {
	Point int `xml:"point,attr"`
	Order int `xml:"order,attr"`
}
