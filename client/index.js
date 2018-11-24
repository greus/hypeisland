function sc() {
    try {
        return stagecast;
    } catch (error) {
        // ignore target defense
    };

    return {
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

const socket = new WebSocket('wss://stagecast.se/api/events/hypeisland/ws');

socket.addEventListener('open', function (event) {
    var json = broadcast({});
    console.log("Client opened socket", json);
});

socket.addEventListener('message', function (event) {
    console.log('Message from server ', event.data);
    document.getElementById("log").innerHTML = event.data;

    var message = JSON.parse(event.data);

    var regardsCurrentUser = message.userId === "1"; // todo sc().userId()

    if (message.type === "client_info" && regardsCurrentUser) {
        setScene(message);
    }
});

function broadcast(obj) {
    obj.userId = sc().getUserId();
    obj.eventId = sc().getEventId();
    obj.momentId = sc().getMomentId();

    var stringified = JSON.stringify(obj);
    socket.send(stringified);
    return obj;
};

document.getElementById('rockIcon').addEventListener('ontouchdstart', () => selectOption('rock'));
document.getElementById('paperIcon').addEventListener('ontouchstart', () => selectOption('paper'));
document.getElementById('scissorIcon').addEventListener('ontouchstart', () => selectOption('scissors'));
setScene();

// UI code
function setScene(message) {
    if (!message) {
        show("messageTop");
        hide("messageBottom");
        hide("rockIcon");
        hide("paperIcon");
        hide("scissorIcon");
        setText("messageTop", "Waiting for more players...");
    } else if (message.view === "match") {
        show("messageTop");
        show("messageBottom");
        show("rockIcon");
        show("paperIcon");
        show("scissorIcon");
        setText("messageTop", "Pick your move!");
        setText("messageBottom", message.info);
    } else if (message.view === "result") {
        hide("rockIcon");
        hide("paperIcon");
        hide("scissorIcon");
        hide("messageTop");
        show("messageBottom");
        setText("messageBottom", message.info);
    } else {
        // end
        hide("messageTop");
        hide("rockIcon");
        hide("paperIcon");
        hide("scissorIcon");
        show("messageBottom");
        setText("messageBottom", "This is the end");
    }
};

function setText(id, text) {
    document.getElementById(id).innerHTML = text;
};

function show(id) {
    document.getElementById(id).classList.remove('hidden');
};

function hide(id) {
    document.getElementById(id).classList.add('hidden');
};

function selectOption(option) {
    console.log("You chose: ", option);
    document.getElementById("messageTop").innerHTML = "You selected " + option + ".";
    document.getElementById("messageBottom").innerHTML = "Waiting for opponents...";
    if (option === "rock") {
        hide("paperIcon");
        hide("scissorIcon");
    } else if (option === "paper") {
        hide("rockIcon");
        hide("scissorIcon");
    } else if (option === "scissors") {
        hide("paperIcon");
        hide("rockIcon");
    }
};
