package Protocol

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func CheckXss(targetURL string, body []byte) (bool, string) {
	foundXSS := false
	ply := ""
	xssPayloads := []string{
		`<sCrIpt>alert("Golin")</SCriPt>`,
		`<img src=x onerror=alert("Golin")>`,
		`<div onmouseover="alert('Golin')">`,
		`<sCrIpt\x09>javascript:alert('Golin')</SCriPt>`,
		`<IMG SRC="javascript:alert('Golin');">`,
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
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

			respDoc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
			if err != nil {
				continue
			}

			xssDetected := false
			respDoc.Find("script, img, div").Each(func(i int, s *goquery.Selection) {
				if s.Is("script") && strings.Contains(s.Text(), "GYS") {
					xssDetected = true
				} else {
					for _, attr := range []string{"onerror", "onmouseover", "src"} {
						attrValue, exists := s.Attr(attr)
						if exists && strings.Contains(attrValue, "GYS") {
							xssDetected = true
							break
						}
					}
				}
				if xssDetected {
					foundXSS = true
					ply = payload
					return
				}
			})

			if xssDetected {
				break
			}
		}
	})

	if foundXSS {
		return true, ply
	}
	return false, ""
}
