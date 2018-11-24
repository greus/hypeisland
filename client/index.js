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
    var json = broadcast({});
    console.log("Client opened socket", json);
});

socket.addEventListener('message', function (event) {
    console.log('Message from server ', event.data);
    document.getElementById("log").innerHTML = event.data;
});

function broadcast(obj) {
    obj.userId = sc.getUserId();
    obj.eventId = sc.getEventId();
    obj.momentId = sc.getMomentId();

    var stringified = JSON.stringify(obj);
    socket.send(stringified);
    return obj;
};

// UI code
function hide(id) {
    document.getElementById(id).classList.add('hidden');
};

function selectOption(option) {
    console.log("You chose: ", option);
    if (option === "rock") {
        hide("paperIcon");
        hide("scissorIcon");
    } else if (option === "paper") {
        hide("rockIcon");
        hide("scissorIcon");
    } else {
        hide("paperIcon");
        hide("rockIcon");
    }
};
