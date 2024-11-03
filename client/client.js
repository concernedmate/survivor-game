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

canvas.width = 1000
canvas.height = 1000

const drawObject = (x, y, size) => {
    ctx.fillRect(x, y, size, size)
}

const render = () => {
    ctx.fillStyle = "black"
    ctx.fillRect(0, 0, canvas.width, canvas.height)
    ctx.clearRect(5, 5, canvas.width - 10, canvas.height - 10)

    if (PLAYERS_DATA != null) {
        PLAYERS_DATA.map((player) => {
            ctx.fillStyle = "black"
            drawObject(player.PosX-player.Size/2, player.PosY-player.Size/2, player.Size)
        })
    }
    if (MOBS_DATA != null) {
        MOBS_DATA.map((mob) => {
            ctx.fillStyle = "red"
            drawObject(mob.PosX-mob.Size/2, mob.PosY-mob.Size/2, mob.Size)
        })
    }
    if (PROJECTILES_DATA != null) {
        PROJECTILES_DATA.map((projectile) => {
            ctx.fillStyle = "green"
            drawObject(projectile.PosX-projectile.Size/2, projectile.PosY-projectile.Size/2, projectile.Size)
        })
    }

    // start rendering
    requestAnimationFrame(render)
}

render()
