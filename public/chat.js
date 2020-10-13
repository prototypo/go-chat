var el = '#app';

var ws = null; // The WebSocket this session
var newMsg = ''; // New message to be sent to the server
var chatContent = ''; // All chat messages displayed on the screen
var username = null; // Username for this session
var joined = false; // True if a user has joined (by submitting a username)

function getWebSocket() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function (e) {
        var msg = JSON.parse(e.data);

        // DEBUG
        console.log("Received from server:\nUsername: " + msg.username + "\nMessage: " + msg.message + "\n\n")

        var content = document.getElementById('chatcontent');
        content.innerHTML += '<div><span class="username">'
            + msg.username
            + ': </span>'
            + msg.message + '</div>';

        content.scrollTop = content.scrollHeight; // Auto scroll to the bottom
    });
    }

function send() {
    this.newMsg = document.getElementById("message").value;

    if (this.username == null || this.username.length < 1) {
        alert("Must join prior to sending messages");
        return
    }
    if (this.newMsg != '') {

        // DEBUG
        console.log("Sending to server: " + this.newMsg);

        this.ws.send(
            JSON.stringify({
                username: this.username,
                message: this.newMsg
            }
            ));
        this.newMsg = ''; // Reset newMsg
        document.getElementById("message").value = ""; // Reset message field
    } else {
        // DEBUG
        console.log("No message to send!");
    }
}

function join() {
    if (!this.username) {
        document.getElementById("userstatus").textContent = "Please choose a username";
        document.getElementById("userstatus").className = "warning";
        return
    }
    this.joined = true;
}

function setUsername() {
    this.username = document.getElementById("username").value;
    document.getElementById("userstatus").textContent = "User: " + this.username;
    document.getElementById("userstatus").className = "displayuser";
    join();
}

window.onload = function () {
    getWebSocket();
    join();
};
