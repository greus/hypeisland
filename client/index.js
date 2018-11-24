var stagecast = stagecast || {
    getUserId: function () {
        return "1234";
    },
    getEventId: function () {
        return "5678";
    },
    getMomentId: function () {
        return "9876";
    },
    getToken: function () {
        return "foobar1337";
    },
    getCoordinates: function () {
        return "56 57";
    },
};

// Create WebSocket connection.
const socket = new WebSocket('wss://stagecast.se/api/events/hypeisland/ws');

// Connection opened
socket.addEventListener('open', function (event) {
    var userId = stagecast.getUserId();
    var json = JSON.stringify({
        userId: stagecast.getUserId(),
        eventId: stagecast.getEventId(),
        momentId: stagecast.getMomentId(),
        token: stagecast.getToken(),
        coordinates: stagecast.getCoordinates()
    });
    console.log("Client opened socket", json);
    socket.send(json);
});

// Listen for messages
socket.addEventListener('message', function (event) {
    console.log('Message from server ', event.data);
});
