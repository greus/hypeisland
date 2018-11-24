if (!stagecast) {
    var stagecast = {
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
};

var sc = stagecast;

const socket = new WebSocket('wss://stagecast.se/api/events/hypeisland/ws');

socket.addEventListener('open', function (event) {
    var userId = sc.getUserId();
    var json = JSON.stringify({
        userId: sc.getUserId(),
        eventId: sc.getEventId(),
        momentId: sc.getMomentId(),
        token: sc.getToken(),
        coordinates: sc.getCoordinates()
    });
    console.log("Client opened socket", json);
    socket.send(json);
});

socket.addEventListener('message', function (event) {
    console.log('Message from server ', event.data);
    document.getElementById("log").innerHTML = event.data;
});
