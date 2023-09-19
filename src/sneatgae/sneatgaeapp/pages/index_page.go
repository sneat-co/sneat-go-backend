package pages

import (
	"fmt"
	"net/http"
	"os"
)

// IndexHandler returns homepage
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	//w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
	<title>Backend @ Sneat.app</title>
	<style>body{font-family: Verdana;}</style>
</head>
<body>
	<h1>Backend for <a href=https://sneat.app>Sneat.app</a></h1>
	<p>
		Developed in <a href=https://golang.org/>Go</a> language
		& hosted @ Google <a href=https://cloud.google.com/appengine/>App Engine</a>
		<a href=https://cloud.google.com/appengine/docs/standard/>Standard</a>
	</p>

	<h2>Google App Engine hosting details</h2>
	<table>
		<tr><td>GAE_RUNTIME</td><td>%v</td><tr>
		<tr><td>GAE_VERSION</td><td>%v</td><tr>
		<tr><td>GAE_DEPLOYMENT_ID</td><td>%v</td><tr>
		<tr><td>GAE_INSTANCE</td><td>%v</td><tr>
		<tr><td>GAE_MEMORY_MB</td><td>%v</td><tr>
	</table>

	<p>&copy; 2020 <a href=https://sneat.team>Sneat.team</a></p>
</body>
</html>
`,
		os.Getenv("GAE_RUNTIME"),
		os.Getenv("GAE_VERSION"),
		os.Getenv("GAE_DEPLOYMENT_ID"),
		os.Getenv("GAE_INSTANCE"),
		os.Getenv("GAE_MEMORY_MB"),
	)
}
