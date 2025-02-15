window.serviceName = "chat";
window.createMessage = function(inputText) {
    let message = {
        type: "chat",
        payload: inputText,
    };
    return message;
}
window.createMessageHTML = function(msg) {
    const message = document.createElement("div");
    message.className = "message";
    message.innerHTML = `
        <div class="user">${msg.metadata.username} (${msg.metadata.time}):</div>
        <div class="text">${msg.payload}</div>
    `;
    return message;
}
