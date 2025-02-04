var conn;
var msg = document.getElementById("msg");
var log = document.getElementById("log");
const username = prompt("Enter username:");

function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

document.getElementById("form").onsubmit = function(event) {
    event.preventDefault();
    if (!conn) {
        return;
    }

    let message = {
        type: "chat",
        payload: msg.value,
        metadata: { username },
    };

    conn.send(JSON.stringify(message));

    msg.value = "";
};

if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/ws/chat");

    conn.onclose = function() {
        var item = document.createElement("div");
        item.innerHTML = "<h2>Connection closed.</h2>";
        appendLog(item);
    };

    conn.onmessage = function(e) {
        var msg = JSON.parse(e.data);
        var item = document.createElement("div");
        item.className = "message";

        if (msg.type == "chat") {
            item.innerText = msg.metadata["username"] + ": " + msg.payload;
        } else {
            item.innerText = "Wrong message type: chat != " + msg.type;
        }

        appendLog(item);
    };
} else {
    var item = document.createElement("div");
    item.innerHTML = "<h2>Your browser does not support WebSockets :(</h2>";
    appendLog(item);
}
