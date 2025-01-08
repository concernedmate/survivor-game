/* =============================== GLOBAL OBJECT =============================== */
/* =============================== GLOBAL OBJECT =============================== */
let PLAYERS_DATA = null
let MOBS_DATA = null
let PROJECTILES_DATA = null
let OBSTACLES_DATA = null

let PLAYER_ID = null

let ID_ROOM = (new URLSearchParams(window.location.search)).get("id_room")

/* =============================== CONNECTION =============================== */
/* =============================== CONNECTION =============================== */
const websocket_client = new WebSocket("ws://" + document.location.host + "/ws_client" + `?id_room=${ID_ROOM}`)
const websocket_server = new WebSocket("ws://" + document.location.host + "/ws_server" + `?id_room=${ID_ROOM}`)

const typeSizes = {
    "undefined": () => 0,
    "boolean": () => 4,
    "number": () => 8,
    "string": item => 2 * item.length,
    "object": item => !item ? 0 : Object
        .keys(item)
        .reduce((total, key) => sizeOf(key) + sizeOf(item[key]) + total, 0)
};
const sizeOf = value => typeSizes[typeof value](value);


let ping = Date.now();
websocket_server.onmessage = (event) => {
    let data = null
    try { data = JSON.parse(event.data) } catch (error) { console.log("Failure to parse JSON data!") }
    if (data != null) {
        PLAYERS_DATA = data.players
        MOBS_DATA = data.mobs
        PROJECTILES_DATA = data.projectiles
    }
    console.log(data.players)
    console.log("Received TOTAL data:", sizeOf(event.data), "bytes")

    // console.log("Ping: ", Date.now() - ping)
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
    event.preventDefault()
    if (event.key == 's') {
        pressedKeys[0] = 1
    }
    if (event.key == 'a') {
        pressedKeys[1] = 1
    }
    if (event.key == 'w') {
        pressedKeys[2] = 1
    }
    if (event.key == 'd') {
        pressedKeys[3] = 1
    }
    if (event.key == ' ') {
        pressedKeys[4] = 1
    }
    if (pressedKeys.toString('') != lastInput.toString('')) {
        websocket_client.send(pressedKeys)
        for (let i = 0; i < pressedKeys.length; i++) {
            lastInput[i] = pressedKeys[i]
        }
    }
})
addEventListener("keyup", (event) => {
    event.preventDefault()
    if (event.key == 's') {
        pressedKeys[0] = 0
    }
    if (event.key == 'a') {
        pressedKeys[1] = 0
    }
    if (event.key == 'w') {
        pressedKeys[2] = 0
    }
    if (event.key == 'd') {
        pressedKeys[3] = 0
    }
    if (event.key == ' ') {
        pressedKeys[4] = 0
    }
    if (pressedKeys.toString('') != lastInput.toString('')) {
        websocket_client.send(pressedKeys)
        for (let i = 0; i < pressedKeys.length; i++) {
            lastInput[i] = pressedKeys[i]
        }
    }
})

/* =============================== GAME WORLD =============================== */
/* =============================== GAME WORLD =============================== */
const canvas = document.getElementById("canvas")
const ctx = canvas.getContext("2d")

canvas.width = 1280
canvas.height = 720

const drawObject = (x, y, size) => {
    ctx.fillRect(x, y, size, size)
}

const render = () => {
    ctx.fillStyle = "black"
    ctx.fillRect(0, 0, canvas.width, canvas.height)
    ctx.clearRect(5, 5, canvas.width - 10, canvas.height - 10)

    let CurrPlayerData = null

    if (PLAYERS_DATA != null) {
        PLAYERS_DATA.map((player) => {
            if (player.Uid == PLAYER_ID) {
                CurrPlayerData = player
                ctx.fillStyle = "grey"
            } else {
                ctx.fillStyle = "black"
            }
            drawObject(player.PosX - player.Size / 2, player.PosY - player.Size / 2, player.Size)
        })
    }
    if (MOBS_DATA != null) {
        ctx.fillStyle = "red"
        const keys = Object.keys(MOBS_DATA)
        for (let i = 0; i < keys.length; i++) {
            let Size = keys[i]
            for (let j = 0; j < MOBS_DATA[Size].length; j += 2) {
                let PosX = MOBS_DATA[Size][j]
                let PosY = MOBS_DATA[Size][j + 1]
                drawObject(PosX - (Size / 2), PosY - (Size / 2), Size)
            }
        }
    }
    if (PROJECTILES_DATA != null) {
        ctx.fillStyle = "green"
        const keys = Object.keys(PROJECTILES_DATA)
        for (let i = 0; i < keys.length; i++) {
            let Size = keys[i]
            for (let j = 0; j < PROJECTILES_DATA[Size].length; j += 2) {
                let PosX = PROJECTILES_DATA[Size][j]
                let PosY = PROJECTILES_DATA[Size][j + 1]
                drawObject(PosX - (Size / 2), PosY - (Size / 2), Size)
            }
        }
    }

    if (CurrPlayerData != null){
        ctx.fillStyle = "black"
        ctx.font = "24px serif";
        
        ctx.fillText("Health:" + CurrPlayerData.Health, 10, 30);  
        ctx.fillText("Score:" + CurrPlayerData.Score, 10, 60);   
    }

    // start rendering
    requestAnimationFrame(render)
}

render()
