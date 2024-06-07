package models

import "encoding/xml"

/*
	normal:
	<gamedata method="gameend">
		<favorite nr="3" musicnums="114?116?103?"/>
		<shoprank point="122" shopname="あかさ" pref="13"/>
		<session flag="0"/>
		<player no="0" cardid="eedd764615c48425" name="%41%42%43%44%45%46%47%48">
			<condition styles="32" hidden="0" recovery="0"/>
			<puzzle no="8" flags="16934" out="0"/>
			<judge perfect="1464" great="0" good="0" poor="0" miss="0"/>
			<play time="610" nr="5" clear="1" gamemode="1">
				<stage musicnum="114" seqmode="1" flags="4096" score="1055825" clear="1" skill="130" combo="194"/>
				<stage musicnum="116" seqmode="1" flags="4096" score="1393235" clear="1" skill="150" combo="201"/>
				<stage musicnum="103" seqmode="1" flags="4096" score="1264100" clear="1" skill="190" combo="215"/>
				<stage musicnum="91" seqmode="1" flags="4096" score="2208110" clear="1" skill="390" combo="323" extra="1"/>
				<stage musicnum="122" seqmode="1" flags="4096" score="6260550" clear="1" skill="510" combo="531" encore="1"/>
			</play>
		</player>
	</gamedata>

	course:
	<gamedata method="gameend">
		<shoprank point="293" shopname="あかさ" pref="13"/>
		<session flag="0"/>
		<player no="0" cardid="eedd764615c48425" name="%41%42%43%44%45%46%47%48">
			<condition styles="32" hidden="0" recovery="0"/>
			<puzzle no="8" flags="21158" out="0"/>
			<judge perfect="0" great="0" good="0" poor="0" miss="0"/>
			<play time="421" nr="0" clear="1" gamemode="3"></play>
		</player>
	</gamedata>
*/

type Request_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Favorite *struct {
		Count     int    `xml:"nr,attr"`
		MusicNums string `xml:"musicnums,attr"`
	} `xml:"favorite"`

	ShopRank *struct {
		Point    int    `xml:"point,attr"`
		ShopName string `xml:"shopname,attr"`
		Pref     int    `xml:"pref,attr"`
	} `xml:"shoprank"`

	Session struct {
		Flag int `xml:"flag,attr"`
	} `xml:"session"`

	Players []struct {
		Number int    `xml:"no,attr"`
		CardId string `xml:"cardid,attr"`
		Name   string `xml:"name,attr"`

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
			Time     int `xml:"time,attr"`
			Count    int `xml:"nr,attr"`
			Clear    int `xml:"clear,attr"`
			GameMode int `xml:"gamemode,attr"`

			Stages []struct {
				MusicNum int `xml:"musicnum,attr"`
				SeqMode  int `xml:"seqmode,attr"`
				Flags    int `xml:"flags,attr"`
				Score    int `xml:"score,attr"`
				Clear    int `xml:"clear,attr"`
				Skill    int `xml:"skill,attr"`
				Combo    int `xml:"combo,attr"`
				Encore   int `xml:"encore,attr"`
				Extra    int `xml:"extra,attr"`
			} `xml:"stage"`
		} `xml:"play"`
	} `xml:"player"`
}

type Response_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`

	Players []Response_GameData_GameEnd_Player `xml:"player"`

	ShopRank Response_GameData_GameEnd_ShopRank `xml:"shoprank"`
}

type Response_GameData_GameEnd_Player_Stage struct {
	All   int `xml:"all,attr"`
	Order int `xml:"order,attr"`
	Best  int `xml:"best,attr"`
}

type Response_GameData_GameEnd_Player struct {
	Number          int                                      `xml:"no,attr"`
	Skill           int                                      `xml:"skill,attr"`
	SkillOrder      int                                      `xml:"skill_order,attr"`
	SkillOrderCount int                                      `xml:"skill_order_nr,attr"`
	SkillAll        int                                      `xml:"skill_all,attr"`
	Stages          []Response_GameData_GameEnd_Player_Stage `xml:"stage"`
}

type Response_GameData_GameEnd_ShopRank struct {
	Hidden int `xml:"hidden,attr"`
}
