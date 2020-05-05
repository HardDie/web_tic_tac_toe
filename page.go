package main

var pageStyle = `
*, *::after, *::before {
  box-sizing: border-box;
}

:root {
  --cell-size: 100px;
  --mark-size: calc(var(--cell-size) * .9);
}

body {
  margin: 0;
}

.board {
  width: 100vw;
  height: 100vh;
  display: grid;
  justify-content: center;
  align-content: center;
  justify-items: center;
  align-items: center;
  grid-template-columns: repeat(3, auto)
}

.cell {
  width: var(--cell-size);
  height: var(--cell-size);
  border: 1px solid black;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  cursor: pointer;
}

.cell:first-child,
.cell:nth-child(2),
.cell:nth-child(3) {
  border-top: none;
}

.cell:nth-child(3n + 1) {
  border-left: none;
}

.cell:nth-child(3n + 3) {
  border-right: none;
}

.cell:last-child,
.cell:nth-child(8),
.cell:nth-child(7) {
  border-bottom: none;
}

.cell.x,
.cell.circle {
  cursor: not-allowed;
}

.cell.x::before,
.cell.x::after,
.cell.circle::before {
  background-color: black;
}

.board.x .cell:not(.x):not(.circle):hover::before,
.board.x .cell:not(.x):not(.circle):hover::after,
.board.circle .cell:not(.x):not(.circle):hover::before {
  background-color: lightgrey;
}

.cell.x::before,
.cell.x::after,
.board.x .cell:not(.x):not(.circle):hover::before,
.board.x .cell:not(.x):not(.circle):hover::after {
  content: '';
  position: absolute;
  width: calc(var(--mark-size) * .15);
  height: var(--mark-size);
}

.cell.x::before,
.board.x .cell:not(.x):not(.circle):hover::before {
  transform: rotate(45deg);
}

.cell.x::after,
.board.x .cell:not(.x):not(.circle):hover::after {
  transform: rotate(-45deg);
}

.cell.circle::before,
.cell.circle::after,
.board.circle .cell:not(.x):not(.circle):hover::before,
.board.circle .cell:not(.x):not(.circle):hover::after {
  content: '';
  position: absolute;
  border-radius: 50%;
}

.cell.circle::before,
.board.circle .cell:not(.x):not(.circle):hover::before {
  width: var(--mark-size);
  height: var(--mark-size);
}

.cell.circle::after,
.board.circle .cell:not(.x):not(.circle):hover::after {
  width: calc(var(--mark-size) * .7);
  height: calc(var(--mark-size) * .7);
  background-color: white;
}

.winning-message {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, .9);
  justify-content: center;
  align-items: center;
  color: white;
  font-size: 5rem;
  flex-direction: column;
}

.winning-message button {
  font-size: 3rem;
  background-color: white;
  border: 1px solid black;
  padding: .25em .5em;
  cursor: pointer;
}

.winning-message button:hover {
  background-color: black;
  color: white;
  border-color: white;
}

.winning-message.show {
  display: flex;
}
/* Forbid select cells text */
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
	<div class="board x" id="board">
	{{range $line, $element := .}}
		{{range $row, $value := $element}}
		<div class="cell noselect" data-cell name="{{$line}}_{{$row}}"></div>
		{{end}}
	{{end}}
	</div>
	<div class="winning-message" id="winningMessage">
		<div data-winning-message-text class="noselect"></div>
		<button id="restartButton" class="noselect">Restart</button>
	</div>

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
	function responseHandler(error, resp) {
		if (error) {
			console.error(error)
			return
		}
		if (resp.Status == "Step made") {
			cells = document.getElementsByName(resp.Step.Line+"_"+resp.Step.Row);
			if (resp.Step.Player === "X") {
				// Set X in cell
				Game.setMark(cells[0], "x");
				// Set O as help
				Game.setHighlight("circle");
			} else if (resp.Step.Player === "O") {
				// Set O in cell
				Game.setMark(cells[0], "circle");
				// Set X as help
				Game.setHighlight("x");
			}
		}
		else if (resp.Status == "Game done") {
			console.log("Game done!");

			cells = document.getElementsByName(resp.Step.Line+"_"+resp.Step.Row);
			if (resp.Step.Player === "X") {
				// Set X in cell
				Game.setMark(cells[0], "x");
				// Set O as help
				Game.setHighlight("circle");
			} else if (resp.Step.Player === "O") {
				// Set O in cell
				Game.setMark(cells[0], "circle");
				// Set X as help
				Game.setHighlight("x");
			}

			Game.showOverMessage("Game over!");
		}
	}


	class Game {
		static setMark(elem, mark) {
			elem.classList.add(mark)
		}
		static setHighlight(mark) {
			board = document.getElementById("board")

			if (mark === "x") {
				board.classList.remove("circle")
				board.classList.add("x")
			} else if (mark === "circle") {
				board.classList.remove("x")
				board.classList.add("circle")
			}
		}
		static showOverMessage(message) {
			const winMsg = document.getElementById("winningMessage");
			const winMsgText = document.querySelector('[data-winning-message-text]')
			winMsgText.innerText = message;
			winMsg.classList.add("show");
		}
		static hideMessage() {
			const winMsg = document.getElementById("winningMessage");
			winMsg.classList.remove("show");
		}
		static resetField() {
			const nodeCells = document.getElementsByClassName("cell");
			const arrayCells = Array.from(nodeCells);
			arrayCells.forEach(function(elem) {
				elem.classList.remove("x");
				elem.classList.remove("circle");
			});
			Game.setHighlight("x");
		}

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
			const button = document.getElementById("restartButton");
			button.addEventListener('click', callback, false);
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
			Http.Post("/api", body, function(){});

			Game.resetField();
			Game.hideMessage();
		}
	}

	Game.registerFieldHandler(Game.eventFieldHandler);
	Game.registerResetHandler(Game.eventResetHandler);

	window.Game = Game
}());
`
