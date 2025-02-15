var conn;
var msg = document.getElementById("msg");
var log = document.getElementById("log");
const username = prompt("Enter username:", "Anonymous");
function getTime() {
    var date = new Date();
    return date.getHours() + ":" + date.getMinutes();
}

document.getElementById("form").onsubmit = function(event) {
    event.preventDefault();
    if (!conn) {
        return;
    }
    let message = window.createMessage(msg.value);
    message.metadata = {
        time: getTime(),
        username: username,
    }
    conn.send(JSON.stringify(message));
    msg.value = "";
};

function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

if (!window["WebSocket"]) {
    var item = document.createElement("h2");
    item.innerText = "Your browser does not support WebSockets :(";
    appendLog(item);
} else {
    conn = new WebSocket("wss://" + document.location.host + "/ws/" + window.serviceName);
    conn.onclose = function() {
        var item = document.createElement("h2");
        item.innerText = "Connection closed.";
        appendLog(item);
    };
    conn.onmessage = function(e) {
        var msg = JSON.parse(e.data);
        if (msg.type != window.serviceName) {
            item.innerText = "Wrong message type: " + window.serviceName + " != " + msg.type;
            return;
        }
        var message = window.createMessageHTML(msg)
        appendLog(message);
    };
}
