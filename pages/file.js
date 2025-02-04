let fileInput = document.getElementById("fileInput");
if (fileInput.files.length > 0) {
    let file = fileInput.files[0];
    let reader = new FileReader();

    reader.onload = function(e) {
        messageType = "file";
        payload = e.target.result;
        metadata = { filename: file.name };

        let message = {
            type: messageType,
            payload: btoa(payload), // Convert to base64
            metadata: metadata
        };

        conn.send(JSON.stringify(message));
    };
    reader.readAsBinaryString(file);
}
