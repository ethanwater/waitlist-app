var socket = new WebSocket("wss://localhost:2695/ws", "protocolOne");

socket.onopen = function(event) {
    console.log("WebSocket connection established.");
};

socket.onmessage = function(event) {
    document.getElementById("timestamp").innerText = event.data;
};

socket.onclose = function(event) {
    console.log("WebSocket connection closed.");
};

socket.onerror = function(event) {
    console.error("WebSocket error:", event);
};
