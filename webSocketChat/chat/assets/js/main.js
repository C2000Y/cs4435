$(function() {

    var conn;
    var msg = $("#msg");
    var log = $("#log");

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    $("#form").submit(function() {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            return false;
        }
        console.log("end")
        conn.send(msg.val())
        msg.val("");
        console.log("end--")
        return false
    });

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://localhost:8011/ws");
        conn.onclose = function(evt) {
            appendLog($("<div><b>Connection Closed.</b></div>"))
        }
        conn.onmessage = function(evt) {
            appendLog($("<div/>").text(evt.data))
        }
    } else {
        appendLog($("<div><b>WebSockets Not Support.</b></div>"))
    }
    });