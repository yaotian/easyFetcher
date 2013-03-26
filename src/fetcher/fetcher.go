package fetcher

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"util"
)

//For cookie
type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	//util.Trace(cookies)
	p.jar[u.Host] = append(p.jar[u.Host], cookies...)
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

type Fetcher struct {
	Client *http.Client
}

func (self *Fetcher) init() {

	//创建client
	transport := &http.Transport{
		Dial: func(netw string, addr string) (net.Conn, error) {
			// we want to wait a maximum of 1.75 seconds...
			// since we're specifying a 1 second connect timeout and deadline 
			// (read/write timeout) is specified in absolute time we want to 
			// calculate that time first (before connecting)
			deadline := time.Now().Add(5 * time.Second)
			c, err := net.DialTimeout(netw, addr, 5*time.Second)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
	}
	client := &http.Client{Transport: transport}

	myjar := &myjar{}
	myjar.jar = make(map[string][]*http.Cookie)
	client.Jar = myjar

	self.Client = client
}

func (self *Fetcher) checkClientReady(retry int) bool {
	for i := 1; i < retry; i++ {
		util.Trace("check Client ready retry:", strconv.Itoa(i))
		if self.Client != nil {
			return true
		} else {
			self.init()
		}
	}
	return false
}

//open a page
func (self *Fetcher) openPage(url string) (string, error) {
	if !self.checkClientReady(3) {
		return "", nil
	}

	web, err_web := self.Client.Get(url)
	if err_web != nil {
		util.Report_Error(err_web)
		return "", err_web
	}

	defer web.Body.Close()
	web_page, err_web_page := ioutil.ReadAll(web.Body)

	if err_web_page != nil {
		util.Report_Error(err_web_page)
		return "", err_web_page
	}

	fmt.Println("done the openPage")

	return string(web_page), nil
}

func Process_Url(fetcher *Fetcher, url string) {
	fmt.Println("start to process url")
	page, _ := fetcher.openPage(url)
	if page != "" {
		process_page(page)
	}
}

func process_page(page string) {
	fmt.Println(page)
	contents, err := util.Use_Reg_Get_Many_Strings(page, "<tr logr=.*?</tr>")
	if err == nil {
		for _, content := range contents {
			process_content(content)
		}
	} else {
		util.Report_Error(err)
	}
}

func process_content(content string) {
	fmt.Println(content)
}
