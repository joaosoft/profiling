package web

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
<head>
	<title>Profiling</title>
</head>
<style>
.content {
	width: 250px;
	height: 300px;
	
	position: absolute;
	top:0;
	bottom: 0;
	left: 0;
	right: 0;
  	
	margin: auto;
}

th {
  font-size: 40px;
}

a {
  font-size: 30px;
}

td {
  font-size: 30px;
}
</style>
<body>
	<div class="content">
		<table>
			<th>Menu</th>
			<tr>
				<td>
					<a href="http://localhost:%d/debug/index">Access pprof tools</a>
				</td>
			</tr>
			<tr>
				<td>
					<a href="http://localhost:%d/report">Generate report</a>
				</td>
			</tr>
		</table>
	</div>
</body>
</html>
`, httpWebServerPort, httpWebServerPort)))
}
