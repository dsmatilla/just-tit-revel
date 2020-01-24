package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dsmatilla/extremetube"
	"github.com/dsmatilla/keezmovies"
	"github.com/dsmatilla/pornhub"
	"github.com/dsmatilla/redtube"
	"github.com/dsmatilla/spankwire"
	"github.com/dsmatilla/tube8"
	"github.com/dsmatilla/xtube"
	"github.com/dsmatilla/youporn"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	search := c.Params.Query.Get("s")
	if len(search) > 0 {
		return c.Redirect(search + ".html")
	}

	values := map[string]interface{}{
		"PageTitle":    "Just Tit",
		"PageMetaDesc": "The most optimized adult video search engine",
	}
	c.ViewArgs = values
	return c.Render()
}

func (c App) Search() revel.Result {
	aux := strings.Replace(c.Request.URL.Path, ".html", "", -1)
	search := strings.Replace(aux, "/", "", -1)
	result := doSearch(search)
	values := map[string]interface{}{
		"PageTitle": fmt.Sprintf("Search results for %s", search),
		"Result": result,
		"PageDesc": fmt.Sprintf("Search results for %s", search),
		"Search": search,
	}
	c.ViewArgs = values
	return c.Render()
}

func (c App) Video() revel.Result {
	aux := strings.Replace(c.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	provider := str[1]
	videoID := str[2]

	BaseDomain := c.Request.URL.Host
	
	var redirect string
	switch provider {
	case "pornhub":
		redirect = "https://pornhub.com/view_video.php?viewkey=" + videoID + "&t=1&utm_source=just-tit.com&utm_medium=embed&utm_campaign=embed-logo-html5"
	case "redtube":
		redirect = "https://www.redtube.com/" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=embed-logo-html5"
	case "tube8":
		redirect = "https://www.tube8.com/video/title/" + videoID + "/?utm_source=just-tit.com&utm_medium=embed&utm_campaign=embed-logo-html5"
	case "youporn":
		redirect = "https://www.youporn.com/watch/" + videoID + "/title/?utm_source=just-tit.com&utm_medium=embed&utm_campaign=embed-logo-html5"
	case "xtube":
		redirect = "https://www.xtube.com/video-watch/watchin-xtube-" + videoID + "?t=0&utm_source=just-tit.com&utm_medium=embed&utm_campaign=embed-logo-html5"
	case "spankwire":
		redirect = "https://www.spankwire.com/title/video" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=embed-logo-html5"
	case "keezmovies":
		redirect = "https://www.keezmovies.com/video/title-" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=embed-logo-html5"
	case "extremetube":
		redirect = "https://www.extremetube.com/video/title-" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=embed-logo"
	}


	type TemplateData = map[string]interface{}
	
	replace := TemplateData{
		"ID":		videoID,
		"Domain":	BaseDomain,
		"Title": "",
	}

	switch provider {
	case "pornhub":
		video := pornhubGetVideoByID(videoID)
		embed := pornhubGetVideoEmbedCode(videoID).Embed.Code
		replace["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(embed)))
		replace["PageTitle"] = fmt.Sprintf("%s", video.Video.Title)
		replace["PageMetaDesc"] = fmt.Sprintf("%s", video.Video.Title)
		replace["Thumb"] = fmt.Sprintf("%s", video.Video.Thumb)
		replace["Url"] = fmt.Sprintf(BaseDomain+"/pornhub/%s.html", videoID)
		replace["Width"] = "580"
		replace["Height"] = "360"
		replace["PornhubVideo"] = video
	case "redtube":
		video := redtubeGetVideoByID(videoID)
		embed := redtubeGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace["Embed"] = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", html.UnescapeString(string(str))))
		replace["PageTitle"] = fmt.Sprintf("%s", video.Video.Title)
		replace["PageMetaDesc"] = fmt.Sprintf("%s", video.Video.Title)
		replace["Thumb"] = fmt.Sprintf("%s", video.Video.Thumb)
		replace["Url"] = fmt.Sprintf(BaseDomain+"/redtube/%s.html", videoID)
		replace["Width"] = "320"
		replace["Height"] = "180"
		replace["RedtubeVideo"] = video
	case "tube8":
		video := tube8GetVideoByID(videoID)
		embed := tube8GetVideoEmbedCode(videoID).EmbedCode.Code
		embed = strings.Replace(embed, "![CDATA[", "", -1)
		embed = strings.Replace(embed, "]]", "", -1)
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace["PageTitle"] = fmt.Sprintf("%s", video.Videos.Title)
		replace["PageMetaDesc"] = fmt.Sprintf("%s", video.Videos.Title)
		if len(video.Videos.Thumbs.Thumb) > 0 {
			replace["Thumb"] = fmt.Sprintf("%s", video.Videos.Thumbs.Thumb[0].Thumb)
		}
		replace["Url"] = fmt.Sprintf(BaseDomain+"/tube8/%s.html", videoID)
		replace["Width"] = "628"
		replace["Height"] = "362"
		replace["Tube8Video"] = video
	case "youporn":
		video := youpornGetVideoByID(videoID)
		embed := youpornGetVideoEmbedCode(videoID).Embed.Code
		replace["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(embed)))
		replace["PageTitle"] = fmt.Sprintf("%s", video.Video.Title)
		replace["PageMetaDesc"] = fmt.Sprintf("%s", video.Video.Title)
		replace["Thumb"] = fmt.Sprintf("%s", video.Video.Thumb)
		replace["Url"] = fmt.Sprintf(BaseDomain+"/youporn/%s.html", videoID)
		replace["Width"] = "628"
		replace["Height"] = "501"
		replace["YoupornVideo"] = video
	case "xtube":
		video := xtubeGetVideoByID(videoID)
		replace["Embed"] = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", video.EmbedCode))
		replace["PageTitle"] = fmt.Sprintf("%s", video.Title)
		replace["PageMetaDesc"] = fmt.Sprintf("%s", video.Description)
		replace["Thumb"] = fmt.Sprintf("%s", video.Thumb)
		replace["Url"] = fmt.Sprintf(BaseDomain+"/xtube/%s.html", videoID)
		replace["Width"] = "628"
		replace["Height"] = "501"
		replace["XtubeVideo"] = video
	case "spankwire":
		video := spankwireGetVideoByID(videoID)
		embed := spankwireGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace["PageTitle"] = fmt.Sprintf("%s", video.Video.Title)
		replace["PageMetaDesc"] = fmt.Sprintf("%s", video.Video.Title)
		replace["Thumb"] = fmt.Sprintf("%s", video.Video.Thumb)
		replace["Url"] = fmt.Sprintf(BaseDomain+"/spankwire/%s.html", videoID)
		replace["Width"] = "650"
		replace["Height"] = "550"
		replace["SpankwireVideo"] = video
	case "keezmovies":
		video := keezmoviesGetVideoByID(videoID)
		embed := keezmoviesGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace["PageTitle"] = fmt.Sprintf("%s", video.Video.Title)
		replace["PageMetaDesc"] = fmt.Sprintf("%s", video.Video.Title)
		replace["Thumb"] = fmt.Sprintf("%s", video.Video.Thumb)
		replace["Url"] = fmt.Sprintf(BaseDomain+"/keezmovies/%s.html", videoID)
		replace["Width"] = "650"
		replace["Height"] = "550"
		replace["KeezmoviesVideo"] = video
	case "extremetube":
		video := extremetubeGetVideoByID(videoID)
		embed := extremetubeGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace["PageTitle"] = fmt.Sprintf("%s", video.Video.Title)
		replace["PageMetaDesc"] = fmt.Sprintf("%s", video.Video.Title)
		replace["Thumb"] = fmt.Sprintf("%s", video.Video.Thumb)
		replace["Url"] = fmt.Sprintf(BaseDomain+"/extremetube/%s.html", videoID)
		replace["Width"] = "650"
		replace["Height"] = "550"
		replace["ExtremetubeVideo"] = video
	default:
		c.Response.Status = 301
		return c.Redirect("/")
	}

	replace["Result"] = doSearch(fmt.Sprintf("%s", replace["PageTitle"]))
	pagetitle := fmt.Sprintf("%v", replace["PageTitle"])
	if pagetitle == "" {
		c.Response.Status = 307
		return c.Redirect(redirect)
	}

	c.ViewArgs = replace
	if c.Params.Query.Get("tp") == "true" {
		return c.RenderTemplate("video/player.html")
	}
	return c.Render()
}

func (c App) ImageProxy() revel.Result {
	image := strings.Replace(c.Request.URL.Path, "/images/", "", -1)
	aux := strings.Split(image, ".")
	str, _ := base64.StdEncoding.DecodeString(aux[0])
	response, _ := http.Get(fmt.Sprintf("%s", str))
	return c.RenderBinary(response.Body, image, "inline", time.Now())
}

type searchResult struct {
	Pornhub pornhub.PornhubSearchResult
	Redtube redtube.RedtubeSearchResult
	Tube8   tube8.Tube8SearchResult
	Youporn youporn.YoupornSearchResult
	Flag    bool
}

var waitGroup sync.WaitGroup

func searchPornhub(search string, c chan pornhub.PornhubSearchResult) {
	defer waitGroup.Done()
	var result pornhub.PornhubSearchResult
	result = pornhub.SearchVideos(search)
	c <- result
	close(c)
}

func searchRedtube(search string, c chan redtube.RedtubeSearchResult) {
	defer waitGroup.Done()
	var result redtube.RedtubeSearchResult
	c <- result
	close(c)
}

func searchTube8(search string, c chan tube8.Tube8SearchResult) {
	defer waitGroup.Done()
	var result tube8.Tube8SearchResult
	result = tube8.SearchVideos(search)
	c <- result
	close(c)
}

func searchYouporn(search string, c chan youporn.YoupornSearchResult) {
	defer waitGroup.Done()
	var result youporn.YoupornSearchResult
	result = youporn.SearchVideos(search)
	c <- result
	close(c)
}

func doSearch(search string) searchResult {
	var cached searchResult
	if err := cache.Get(search, &cached); err == nil {
		return cached
	} else {
		log.Print("Cache NOT found")
		waitGroup.Add(4)

		PornhubChannel := make(chan pornhub.PornhubSearchResult)
		RedtubeChannel := make(chan redtube.RedtubeSearchResult)
		Tube8Channel := make(chan tube8.Tube8SearchResult)
		YoupornChannel := make(chan youporn.YoupornSearchResult)

		go searchPornhub(search, PornhubChannel)
		go searchRedtube(search, RedtubeChannel)
		go searchTube8(search, Tube8Channel)
		go searchYouporn(search, YoupornChannel)

		result := searchResult{<-PornhubChannel, <-RedtubeChannel, <-Tube8Channel, <-YoupornChannel, true}

		waitGroup.Wait()

		go cache.Set(search, result, 24 * time.Hour)
		return result
	}
}

func pornhubGetVideoByID(videoID string) pornhub.PornhubSingleVideo {
	cachedElement := getFromDB("pornhub-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result pornhub.PornhubSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][PORNHUB_GET]", err)
		}
		return result
	} else {
		video := pornhub.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("pornhub-video-"+videoID, string(jsonResult))
		return video
	}
}

func pornhubGetVideoEmbedCode(videoID string) pornhub.PornhubEmbedCode {
	cachedElement := getFromDB("pornhub-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result pornhub.PornhubEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][PORNHUB_EMBED]", err)
		}
		return result
	} else {
		embed := pornhub.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("pornhub-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func redtubeGetVideoByID(videoID string) redtube.RedtubeSingleVideo {
	cachedElement := getFromDB("redtube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result redtube.RedtubeSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][REDTUBE_GET]", err)
		}
		return result
	} else {
		video := redtube.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("redtube-video-"+videoID, string(jsonResult))
		return video
	}
}

func redtubeGetVideoEmbedCode(videoID string) redtube.RedtubeEmbedCode {
	cachedElement := getFromDB("redtube-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result redtube.RedtubeEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][REDTUBE_EMBED]", err)
		}
		return result
	} else {
		embed := redtube.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("redtube-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func tube8GetVideoByID(videoID string) tube8.Tube8SingleVideo {
	cachedElement := getFromDB("tube8-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result tube8.Tube8SingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][TUBE8_GET]", err)
		}
		return result
	} else {
		video := tube8.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("tube8-video-"+videoID, string(jsonResult))
		return video
	}
}

func tube8GetVideoEmbedCode(videoID string) tube8.Tube8EmbedCode {
	cachedElement := getFromDB("tube8-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result tube8.Tube8EmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][TUBE8_EMBED]", err)
		}
		return result
	} else {
		embed := tube8.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("tube8-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func youpornGetVideoByID(videoID string) youporn.YoupornSingleVideo {
	cachedElement := getFromDB("youporn-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result youporn.YoupornSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][YOUPORN_GET]", err)
		}
		return result
	} else {
		video := youporn.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("youporn-video-"+videoID, string(jsonResult))
		return video
	}
}

func youpornGetVideoEmbedCode(videoID string) youporn.YoupornEmbedCode {
	cachedElement := getFromDB("youporn-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result youporn.YoupornEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][YOUPORN_EMBED]", err)
		}
		return result
	} else {
		embed := youporn.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("youporn-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func xtubeGetVideoByID(videoID string) xtube.XtubeVideo {
	cachedElement := getFromDB("xtube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result xtube.XtubeVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][XTUBE_GET]", err)
		}
		return result
	} else {
		video := xtube.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("xtube-video-"+videoID, string(jsonResult))
		return video
	}
}

func spankwireGetVideoByID(videoID string) spankwire.SpankwireSingleVideo {
	cachedElement := getFromDB("spankwire-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result spankwire.SpankwireSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][SPANKWIRE_GET]", err)
		}
		return result
	} else {
		video := spankwire.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("spankwire-video-"+videoID, string(jsonResult))
		return video
	}
}

func spankwireGetVideoEmbedCode(videoID string) spankwire.SpankwireEmbedCode {
	cachedElement := getFromDB("spankwire-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result spankwire.SpankwireEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][SPANKWIRE_EMBED]", err)
		}
		return result
	} else {
		embed := spankwire.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("spankwire-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func keezmoviesGetVideoByID(videoID string) keezmovies.KeezmoviesSingleVideo {
	cachedElement := getFromDB("keezmovies-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result keezmovies.KeezmoviesSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][KEEZMOVIES_GET]", err)
		}
		return result
	} else {
		video := keezmovies.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("keezmovies-video-"+videoID, string(jsonResult))
		return video
	}
}

func keezmoviesGetVideoEmbedCode(videoID string) keezmovies.KeezmoviesEmbedCode {
	cachedElement := getFromDB("keezmovies-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result keezmovies.KeezmoviesEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][KEEZMOVIES_EMBED]", err)
		}
		return result
	} else {
		embed := keezmovies.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("keezmovies-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func extremetubeGetVideoByID(videoID string) extremetube.ExtremetubeSingleVideo {
	cachedElement := getFromDB("extremetube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result extremetube.ExtremetubeSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][EXTREMETUBE_GET]", err)
		}
		return result
	} else {
		video := extremetube.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("extremetube-video-"+videoID, string(jsonResult))
		return video
	}
}

func extremetubeGetVideoEmbedCode(videoID string) extremetube.ExtremetubeEmbedCode {
	cachedElement := getFromDB("extremetube-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result extremetube.ExtremetubeEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][EXTREMETUBE_EMBED]", err)
		}
		return result
	} else {
		embed := extremetube.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("extremetube-embed-"+videoID, string(jsonResult))
		return embed
	}
}

type JustTitCache struct {
	ID        string `json:"id"`
	Result    string `json:"result"`
	Timestamp int64  `json:"timestamp"`
}

func getFromDB(ID string) JustTitCache {
	var result string
	_ = cache.Get(ID, &result)
	var cached JustTitCache
	json.Unmarshal([]byte(result), &cached)
	return cached
}

func putToDB(ID string, Result string) {
	go cache.Set(ID, Result, time.Hour * 24)
}