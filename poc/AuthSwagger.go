package poc

import (
	"io"
	"net/http"
	"strings"
)

func AuthSwagger(url string) {

	paths := []string{
		"/swagger/ui/index",
		"/swagger-ui.html",
		"/api/swagger-ui.html",
		"/service/swagger-ui.html",
		"/web/swagger-ui.html",
		"/swagger/swagger-ui.html",
		"/actuator/swagger-ui.html",
		"/libs/swagger-ui.html",
		"/template/swagger-ui.html",
		"/api_docs",
		"/api/docs/",
		"/api/index.html",
		"/swagger/v1/swagger.yaml",
		"/swagger/v1/swagger.json",
		"/swagger.yaml",
		"/swagger.json",
		"/api-docs/swagger.yaml",
		"/api-docs/swagger.json",
	}

	for _, path := range paths {
		req, _ := http.NewRequest("GET", url+path, nil)
		resp, err := newRequest(req)

		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			bodyBytes, err2 := io.ReadAll(resp.Body)
			if err2 != nil {
				continue
			}
			bodyString := string(bodyBytes)

			if strings.Contains(bodyString, "Swagger UI") ||
				strings.Contains(bodyString, "swagger-ui.min.js") ||
				strings.Contains(bodyString, "swagger:") ||
				strings.Contains(bodyString, "Swagger 2.0") ||
				strings.Contains(bodyString, "\"swagger\":") {
				flags := Flagcve{url + path, "swagger未授权访问", ""}
				echoFlag(flags)
			}
		}
	}

}
