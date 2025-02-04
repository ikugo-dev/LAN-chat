var conn;
var msg = document.getElementById("msg");
var log = document.getElementById("log");
const username = prompt("Enter username:", "Anonymous");

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
    var date = new Date();
    var time = date.getHours() + ":" + date.getMinutes();


    let message = {
        type: "chat",
        payload: msg.value,
        metadata: { username, time },
    };

    conn.send(JSON.stringify(message));

    msg.value = "";
};

if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/ws/chat");

    conn.onclose = function() {
        var item = document.createElement("h2");
        item.innerText = "Connection closed.";
        appendLog(item);
    };

    conn.onmessage = function(e) {
        var msg = JSON.parse(e.data);
        if (msg.type != "chat") {
            item.innerText = "Wrong message type: chat != " + msg.type;
            return;
        }

        var item = document.createElement("div");
        item.className = "message";

        const { time, username } = msg.metadata;

        const user = document.createElement("div");
        user.className = "user";
        user.innerText = `${username} (${time}):`;

        const text = document.createElement("div");
        text.className = "text";
        text.innerText = msg.payload;

        item.appendChild(user);
        item.appendChild(text);
        appendLog(item);
    };
} else {
    var item = document.createElement("h2");
    item.innerText = "Your browser does not support WebSockets :(";
    appendLog(item);
}
