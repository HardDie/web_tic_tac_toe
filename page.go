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
.button {
	box-shadow:inset 0px 1px 0px 0px #ffffff;
	background:linear-gradient(to bottom, #f9f9f9 5%, #e9e9e9 100%);
	background-color:#f9f9f9;
	border-radius:6px;
	border:1px solid #dcdcdc;
	display:inline-block;
	cursor:pointer;
	color:#666666;
	font-family:Arial;
	font-size:15px;
	font-weight:bold;
	padding:6px 24px;
	text-decoration:none;
	text-shadow:0px 1px 0px #ffffff;
}
.button:hover {
	background:linear-gradient(to bottom, #e9e9e9 5%, #f9f9f9 100%);
	background-color:#e9e9e9;
}
.button:active {
	position:relative;
	top:1px;
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

	<input type="button" value="Restart" class="button" name="resetButton">

	<script type="text/javascript" src="script.js"></script>
</body>
</html>
`

var pageScript = `
(function() {
	/**
	 * Help send http request to server
	 */
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

	/**
	 * Process response from server and fill cells
	 */
	function responseHandler(error, response) {
		response.forEach(function(elem) {
			cells = document.getElementsByName(elem.Line + "_" + elem.Row);
			cells[0].innerText = elem.Value;
		});
	}


	class Game {
		/**
		 * Register event handler on field buttons
		 */
		static registerFieldHandler(callback) {
			const nodeCells = document.getElementsByClassName("cell");
			const arrayCells = Array.from(nodeCells);
			arrayCells.forEach(function(elem) {
				elem.addEventListener('click', callback, false);
			});
		}

		/**
		 * Register event handler on reset buttons
		 */
		static registerResetHandler(callback) {
			const nodeButtons = document.getElementsByClassName("button");
			const arrayButtons = Array.from(nodeButtons);
			arrayButtons.forEach(function(elem) {
				elem.addEventListener('click', callback, false);
			});
		}

		/**
		 * Handler field button
		 */
		static eventFieldHandler() {
			const Name = this.getAttribute('name');
			const Type = "Set"
			const body = {Type, Name}
			Http.Post("/api", body, responseHandler);
		}

		/**
		 * Handler reset button
		 */
		static eventResetHandler() {
			const Type = "Reset"
			const body = {Type}
			Http.Post("/api", body, responseHandler);
		}
	}

	Game.registerFieldHandler(Game.eventFieldHandler);
	Game.registerResetHandler(Game.eventResetHandler);

	window.Game = Game
}());
`
