var stagecast = stagecast || {
    getUserId: function () {
        return "1234"
    }
};

// Create WebSocket connection.
const socket = new WebSocket('wss://stagecast.se/api/events/hypeisland/ws');

// Connection opened
socket.addEventListener('open', function (event) {
    console.log(stagecast);
    var userId = stagecast.getUserId();
    var json = JSON.stringify({
        userId: stagecast.getUserId()
    });
    socket.send(json);
});

// Listen for messages
socket.addEventListener('message', function (event) {
    console.log('Message from server ', event.data);
});
