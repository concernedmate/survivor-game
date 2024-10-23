/* =============================== GLOBAL OBJECT =============================== */
/* =============================== GLOBAL OBJECT =============================== */
let PLAYERS_DATA = null
let MOBS_DATA = null
let PROJECTILES_DATA = null
let OBSTACLES_DATA = null

let PLAYER_ID = null

/* =============================== CONNECTION =============================== */
/* =============================== CONNECTION =============================== */
const websocket_client = new WebSocket("ws://" + document.location.host + "/ws_client")
const websocket_server = new WebSocket("ws://" + document.location.host + "/ws_server")

let ping = Date.now();
websocket_server.onmessage = (event) => {
    let data = null
    try { data = JSON.parse(event.data) } catch (error) { console.log("Failure to parse JSON data!") }
    if (data != null) {
        PLAYERS_DATA = data.players
        MOBS_DATA = data.mobs
        PROJECTILES_DATA = data.projectiles
        console.log("Player count: ", data.players.length)
    }
    console.log("Ping: ", Date.now() - ping)
    ping = Date.now()
}
websocket_client.onmessage = (event) => {
    let data = null
    try { data = JSON.parse(event.data) } catch (error) { console.log("Failure to parse JSON data!") }
    console.log("Player ID: ", data)
    PLAYER_ID = data.id
}

/* =============================== CONTROLS =============================== */
/* =============================== CONTROLS =============================== */
const pressedKeys = [0, 0, 0, 0, 0]
const lastInput = [0, 0, 0, 0, 0]

addEventListener("keydown", (event) => {
    if (event.key == 'w') {
        pressedKeys[0] = 1
    }
    if (event.key == 'a') {
        pressedKeys[1] = 1
    }
    if (event.key == 's') {
        pressedKeys[2] = 1
    }
    if (event.key == 'd') {
        pressedKeys[3] = 1
    }
    if (event.key == ' ') {
        pressedKeys[4] = 1
    }
})
addEventListener("keyup", (event) => {
    if (event.key == 'w') {
        pressedKeys[0] = 0
    }
    if (event.key == 'a') {
        pressedKeys[1] = 0
    }
    if (event.key == 's') {
        pressedKeys[2] = 0
    }
    if (event.key == 'd') {
        pressedKeys[3] = 0
    }
    if (event.key == ' ') {
        pressedKeys[4] = 0
    }
})

const interval = setInterval(() => {
    if (pressedKeys.toString('') != lastInput.toString('')) {
        websocket_client.send(pressedKeys)
        for (let i = 0; i < pressedKeys.length; i++) {
            lastInput[i] = pressedKeys[i]
        }
    }
}, 10)

/* =============================== GAME WORLD =============================== */
/* =============================== GAME WORLD =============================== */
const canvas = document.getElementById("canvas")
const canvasCtx = canvas.getContext("2d")