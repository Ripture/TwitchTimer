$(function () {
    window.onPlayerEvent = function (data) {
    data.forEach(function(event) {
      if (event.event == "playerInit") {
        $("#twitch_player")[0].unmute()
      }

      if (event.event == "videoPlaying"){
        enableNextButton()
      }

      if(event.event == "offline"){
        enableNextButton()
      }

    })
    }
    $("#nextButton").css("width", window.innerWidth-15);
    $("#nextButton").css("height", 50);
    $("#nextButton").css("font-size", 30);
  }
  )

var firstload = true     //first load variable to track when first time opening streams so we can leave one unmuted
var counter = 20
var counting
var paused = false
var frontontop = true

var serversocket = new WebSocket ("ws://" + ip + "/requestStreamer");
serversocket.onmessage = function(e) {
  $("#twitch_player").empty()
  createPlayer(e.data, "twitch_player")
  $("#streamerName")[0].innerHTML = e.data
}
serversocket.onopen = function(e) {
  //socket is now opened, get initial streamer
  getStreamer()
}

function getStreamer(){
  disableNextButton()
  serversocket.send("requestStreamer");
}

function createPlayer(streamerName, playerName) {
  swfobject.embedSWF("//www-cdn.jtvnw.net/swflibs/TwitchPlayer.swf", playerName, window.innerWidth-15, window.innerHeight-200, "11", null,
      { "eventsCallback":"onPlayerEvent",
        "embed":1,
        "channel":streamerName,
        "auto_play":"true",
        "start_volume":"50"},
      { "allowScriptAccess":"always",
        "allowFullScreen":"true"});
}

function enableNextButton(){
  $("#nextButton")[0].disabled = "";
  $("#nextButton")[0].value = "NEXT";
}

function disableNextButton(){
  $("#nextButton")[0].disabled = "true";
  $("#nextButton")[0].value= "LOADING NEXT STREAM...";
}

// id = setInterval(timer, 1000);

// function pauseTimer(){
//   if (paused){
//     id = setInterval(timer, 1000);
//     document.getElementById("pausebtn").innerHTML = "Pause";
//   }
//   else {
//     document.getElementById("pausebtn").innerHTML = "Unpause";
//     clearInterval(id)
//   }
//   paused = !paused;
// }

function timer(){
  counter--;
  if(counter < 0) {
    document.getElementById("go").click();
  } else {
    document.getElementById("timer").innerHTML = counter.toString();
    }
}

function resize(){
  $("#nextButton").css("width", window.innerWidth-15);
  $("#twitch_player").attr("width", window.innerWidth-15)
  $("#twitch_player").attr("height", window.innerHeight-200)
}