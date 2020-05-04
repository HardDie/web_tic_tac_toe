package main

var pageStyle = `
body {
	margin: 0;
	padding: 0;
}
.gameTable {
	border-collapse: collapse;
	border-spacing: 0;
}
td {
	width: 80px;
	height: 80px;
	border: 1px solid #1c1c1c;
	font-family: Consolas;
	text-align: center;
	font-size: 30px;
}
td:hover {
	background: #e4e4e4;
	cursor: pointer;
}
.noselect {
	-webkit-touch-callout: none; /* iOS Safari */
	-webkit-user-select: none;   /* Safari */
	-khtml-user-select: none;    /* Konqueror HTML */
	-moz-user-select: none;      /* Old versions of Firefox */
	-ms-user-select: none;       /* Internet Explorer/Edge */
	user-select: none;           /* Non-prefixed version, currently
	                                supported by Chrome, Opera and Firefox */
}
`

var pageIndexTmpl = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Tic Tac Toe</title>
	<link rel="stylesheet prefetch" type="text/css" href="style.css">
</head>

<body>
	<table class="gameTable">
		{{range $line, $element := .Array}}
		<tr class="row">
			{{range $row, $value := $element}}
			<td class="cell noselect" name="{{$line}}_{{$row}}">{{$value}}</td>
			{{end}}
		</tr>
		{{end}}
	</table>

	<script type="text/javascript" src="script.js"></script>
</body>
</html>
`

var pageScript = `
(function() {
	class Http {
		static Post(path, body, callback) {
			const xhr = new XMLHttpRequest();
			xhr.open("POST", path, true);
			xhr.withCredentials = true;
			xhr.setRequestHeader('Content-Type', 'application/json; charset=utf-8');

			xhr.onreadystatechange = function() {
				if (xhr.readyState !== 4) return;
				if (+xhr.status !== 200) {
					return callback(xhr, null);
				}

				const response = JSON.parse(xhr.responseText);
				callback(null, response);
			}

			xhr.send(JSON.stringify(body));
		}
	}

	function responseHandler(error, response) {
		response.forEach(function(elem) {
			cells = document.getElementsByName(elem.Line + "_" + elem.Row);
			cells[0].innerText = elem.Value;
		});
	}


	class Game {
		static registerHandler(callback) {
			const nodeCells = document.getElementsByClassName("cell");
			const arrayCells = Array.from(nodeCells);
			arrayCells.forEach(function(elem) {
				elem.addEventListener('click', callback, false);
			});
		}

		static eventHandler() {
			const Name = this.getAttribute('name');
			const body = {Name}
			Http.Post("/api", body, responseHandler);
		}
	}

	Game.registerHandler(Game.eventHandler);

	window.Game = Game
}());
`
