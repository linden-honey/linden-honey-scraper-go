package parser

import (
	"encoding/json"
	"testing"
)

func TestParseSong(t *testing.T) {
	html := `

	<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
	<html>
	<head>
	<meta name="generator" content="Воспалённая фантазия василыча">
	<META NAME="Robots" CONTENT="INDEX,FOLLOW">
	<META NAME="Document-state" CONTENT="Dynamic">
	<META NAME="Revizit-after" CONTENT="10 days">
	<META HTTP-EQUIV="Content-type" CONTENT="text/html;charset=windows-1251">
	<meta name="Keywords" lang="ru" content="го, г.о, г.о., гроб, гр.об., гр.об, гражданская оборона, группа гражданская оборона, группа ГО, Летов, Егор Летов, Игорь Летов, Е.Ф. Летов, Янка Дягилева, Егор и Опизденевшие, Коммунизм, Посев, панки, хой, анархия, андерграунд, русский прорыв, punk, go, gr.ob, grob, oborona, egor letov, Инструкция по выживанию, требьют, требьют гражданской обороны, требьют го, по плану, аккорды го, grob-records.go.ru, гроб-рекордз, новый альбом ГО, Звездопад">
	<title>Гражданская Оборона - официальный сайт группы | Всё идёт по плану</title>
	<style type="text/css" media="print">
	sup	{
	color:#CC0000;}
	a:link		{color: #000000; text-decoration: underline; font-weight:bold;}
	a:visited	{color: #000000; text-decoration: underline; font-weight:bold;}
	a:hover		{color: #000000; font-weight:bold; text-decoration: underline;}
	a:active	{color: #000000; text-decoration: underline; font-weight:bold;}
	.img {
	border:1px solid; border-color:#808080}
	</style>
	<style type="text/css" media="screen">
	sup	{
	color:#CC0000;}
	a:link		{color: #000000; text-decoration: underline; font-weight:bold;}
	a:visited	{color: #000000; text-decoration: underline; font-weight:bold;}
	a:hover		{color: #000000; font-weight:bold; text-decoration: underline;}
	a:active	{color: #000000; text-decoration: underline; font-weight:bold;}
	.img {
	border:1px solid; border-color:#808080}
	</style>
	</head>
	<script language="JavaScript">
	var path = window.opener.location.href;
	</script>
	<body>
	<hr size="2">
	<h2>Всё идёт по плану</h2><hr size="1">
	<p><strong>Автор:</strong> Е.Летов</p><p><strong>Альбом:</strong> Всё идёт по плану</p><p>Границы&nbsp;ключ&nbsp;переломлен&nbsp;пополам<br>А&nbsp;наш&nbsp;батюшка&nbsp;Ленин&nbsp;совсем&nbsp;усоп<br>Он&nbsp;разложился&nbsp;на&nbsp;плесень&nbsp;и&nbsp;на&nbsp;липовый&nbsp;мёд<br>А&nbsp;перестройка&nbsp;всё&nbsp;идёт&nbsp;и&nbsp;идёт&nbsp;по&nbsp;плану<br>А&nbsp;вся&nbsp;грязь&nbsp;превратилась&nbsp;в&nbsp;голый&nbsp;лёд&nbsp;&nbsp;<br>&nbsp;&nbsp;&nbsp;<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;и&nbsp;всё&nbsp;идёт&nbsp;по&nbsp;плану<br><br>А&nbsp;моя&nbsp;судьба&nbsp;захотела&nbsp;на&nbsp;покой<br>Я&nbsp;обещал&nbsp;ей&nbsp;не&nbsp;участвовать&nbsp;в&nbsp;военной&nbsp;игре<br>Но&nbsp;на&nbsp;фуражке&nbsp;на&nbsp;моей&nbsp;серп&nbsp;и&nbsp;молот&nbsp;и&nbsp;звезда<br>Как&nbsp;это&nbsp;трогательно&nbsp;—&nbsp;серп&nbsp;и&nbsp;молот&nbsp;и&nbsp;звезда<br>Лихой&nbsp;фонарь&nbsp;ожидания&nbsp;мотается<br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;и&nbsp;всё&nbsp;идёт&nbsp;по&nbsp;плану<br><br>А&nbsp;моей&nbsp;женой&nbsp;накормили&nbsp;толпу<br>Мировым&nbsp;кулаком&nbsp;растоптали&nbsp;ей&nbsp;грудь<br>Всенародной&nbsp;свободой&nbsp;растерзали&nbsp;ей&nbsp;плоть<br>Так&nbsp;закопайте&nbsp;её&nbsp;во&nbsp;Христе&nbsp;-<br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;ведь&nbsp;всё&nbsp;идёт&nbsp;по&nbsp;плану<br><br>Один&nbsp;лишь&nbsp;дедушка&nbsp;Ленин&nbsp;хороший&nbsp;был&nbsp;вождь<br>А&nbsp;все&nbsp;другие&nbsp;остальные&nbsp;—&nbsp;такое&nbsp;дерьмо<br>А&nbsp;все&nbsp;другие&nbsp;враги&nbsp;и&nbsp;такие&nbsp;дураки<br>Над&nbsp;родною&nbsp;над&nbsp;отчизной&nbsp;бесноватый&nbsp;снег&nbsp;шёл<br>Я&nbsp;купил&nbsp;журнал&nbsp;&laquo;Корея&raquo;&nbsp;—&nbsp;там&nbsp;тоже&nbsp;хорошо<br>Там&nbsp;товарищ&nbsp;Ким&nbsp;Ир&nbsp;Сен&nbsp;там&nbsp;тоже&nbsp;что&nbsp;у&nbsp;нас<br>Я&nbsp;уверен,что&nbsp;у&nbsp;них&nbsp;тоже&nbsp;самое&nbsp;—&nbsp;<br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;и&nbsp;всё&nbsp;идёт&nbsp;по&nbsp;плану<br><br>А&nbsp;при&nbsp;коммунизме&nbsp;всё&nbsp;будет&nbsp;заебись<br>Он&nbsp;наступит&nbsp;скоро&nbsp;—&nbsp;надо&nbsp;только&nbsp;подождать<br>Там&nbsp;всё&nbsp;будет&nbsp;бесплатно,там&nbsp;всё&nbsp;будет&nbsp;в&nbsp;кайф<br>Там&nbsp;наверное&nbsp;вощще&nbsp;не&nbsp;надо&nbsp;будет&nbsp;(умирать)<br>Я&nbsp;проснулся&nbsp;среди&nbsp;ночи&nbsp;и&nbsp;понял,&nbsp;что&nbsp;-<br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;ВСЁ&nbsp;ИДЁТ&nbsp;ПО&nbsp;ПЛАНУ</p><hr size="1">
	<small><script language="JavaScript">
	document.write("URL этой страницы: <a href=\"" + path + "\">" + path + "</a>");
	</script><br>
	Головной URL сайта: <a href="http://www.gr-oborona.ru">http://www.gr-oborona.ru</a><br>
	Все вопросы, замечания и предложения направляйте по электронной почте: <a href="mailto:info@mirgorod.ru">info@mirgorod.ru</a><br>
	<br>
	Копирование материалов с сайта в зоне WWW допускается только со ссылкой на источник.<br>
	Копирование материалов вне зоны WWW без разрешения запрещено. <br>
	Официальный сайт группы «Гражданская Оборона» 2020 г.</small>
	<hr size="2">
	</body>
	</html>
	`
	song := ParseSong(html)
	t.Log(song)
	bytes, _ := json.Marshal(song)
	t.Log(string(bytes))
}

func TestParsePreviews(t *testing.T) {
	html := `
	<ul id="abc_list">
	<li><a href="/texts/1056899068.html">Всё идёт по плану</a></li>
	<li><a href=""></a></li>
	<li><a href="">Unknown</a></li>
	<li><a href="/texts/1056901056.html">Всё как у людей</a></li>
	</ul>
	`
	previews := ParsePreviews(html)
	t.Log(previews)
	if len(previews) != 2 {
		t.Error("Parsing previews failed")
	}
}
