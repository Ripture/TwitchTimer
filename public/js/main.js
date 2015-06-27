
var firstload = true;     //first load variable to track when first time opening streams so we can leave one unmuted
var frontPlayerPlaying = false; //track the states of players
var backPlayerPlaying = false;  // and whether or not they are playing

$(function () {
    window.onFrontPlayerEvent = function (data) {
    data.forEach(function(event) {
      if (event.event == "playerInit") {
        $("#twitch_player_front")[0].mute();
        if (firstload){
          $("#twitch_player_front")[0].unmute();
          firstload = false;
        }

      }

      if (event.event == "videoPlaying"){
        frontPlayerPlaying = true;
        if(frontPlayerPlaying & backPlayerPlaying) {
          enableNextButton()
        }
      }

      if(event.event == "offline"){
        frontPlayerPlaying = true;
        if(frontPlayerPlaying & backPlayerPlaying) {
          enableNextButton()
        }
      }

    });
  }

  window.onBackPlayerEvent = function (data) {
    data.forEach(function(event) {
      if (event.event == "playerInit") {
    var player = $("#twitch_player_back")[0];
          player.mute();
        }

        if (event.event == "videoPlaying") {
          backPlayerPlaying = true;
          if(frontPlayerPlaying & backPlayerPlaying) {
            enableNextButton()
          }
        }

        if(event.event == "offline"){
          backPlayerPlaying = true;
          if(frontPlayerPlaying & backPlayerPlaying) {
            enableNextButton()
          }
        }

      });
    }

      replaceFrontPlayer("nl_Kripp");
      replaceBackPlayer("Calebhart42");
  });

function enableNextButton(){
  $("#nextButton")[0].disabled = "";
  $("#nextButton")[0].value = "NEXT";
  $("#nextButton").css("width", window.innerWidth);
  $("#nextButton").css("height", 50);
}

function disableNextButton(){
  $("#nextButton")[0].disabled = "true";
  $("#nextButton")[0].value= "LOADING NEXT STREAM...";
  $("#nextButton").css("width", window.innerWidth);
  $("#nextButton").css("height", 50);
}

function fixButtonCSS(){
    $("#nextButton").css("width", window.innerWidth-30);
    $("#nextButton").css("height", 50);
    $("#nextButton").css("font-size", 30);
  }

function replaceFrontPlayer(name) {
  swfobject.embedSWF("//www-cdn.jtvnw.net/swflibs/TwitchPlayer.swf", "twitch_player_front", "640", "400", "11", null,
      { "eventsCallback":"onFrontPlayerEvent",
        "embed":1,
        "channel":name,
        "auto_play":"true",
        "start_volume":"50"},
      { "allowScriptAccess":"always",
        "allowFullScreen":"true"});
}

function replaceBackPlayer(name) {

    swfobject.embedSWF("//www-cdn.jtvnw.net/swflibs/TwitchPlayer.swf", "twitch_player_back", "640", "400", "11", null,
      { "eventsCallback":"onBackPlayerEvent",
        "embed":1,
        "channel":name,
        "auto_play":"true",
        "start_volume":"50"},
      { "allowScriptAccess":"always",
        "allowFullScreen":"true"});
}
// window.onload = new Function("changeChannel('nl_Kripp');");

var counter = 20;
var counting;
var paused = false;
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

function changeChannel(name){
  if(frontontop){                       //IF THE FRONT PLAYER IS ON TOP
    $("#twitch_player_back").empty();   //remove and replace the back player
    replaceBackPlayer(name);            //create back player
    $("#twitch_player_back")[0].mute(); //mute the back player
  }
  else {                                  //IF THE BACK PLAYER IS ON TOP
    $("#twitch_player_front").empty();    //remove and replace the front player
    replaceFrontPlayer(name);             //create the front player
    $("#twitch_player_front")[0].mute();  //mute the front player
  }
}

var frontontop = true;
function nextStreamer(){

  if(frontontop){                             //IF THE FRONT PLAYER IS ON TOP (CURRENT PLAYER)
    disableNextButton();
    $("#twitch_player_front").hide();         //hide the front player
    $("#twitch_player_back")[0].unmute();     //unmute the back player
    $("#back_player_div").css('z-index', 1);  //move the back player to top of z
    $("#front_player_div").css('z-index', -1);//move the front player to back of z

    frontontop = false                        //set variable to track which player is on top
    serversocket.send("requestStreamer");     //request a new streamer
  }
  else{                                       //IF THE BACK PLAYER IS ON TOP (CURRENT PLAYER)
    disableNextButton();
    $("#twitch_player_back").hide();          //hide the back player
    $("#twitch_player_front")[0].unmute();    //unmute the front player
    $("#front_player_div").css('z-index', 1); //move the front player to top of z
    $("#back_player_div").css('z-index', -1); //move the back player to back of z

    frontontop = true                        //set variable to track which player is on top
    serversocket.send("requestStreamer");    //request a new streamer
  }

}

//var serversocket = new WebSocket ("ws://70.161.150.36:1935/requestStreamer");
var serversocket = new WebSocket ("ws://70.161.150.36:1935/requestStreamer");

serversocket.onmessage = function(e) {
  changeChannel(e.data);
  //it never gets past the above func call, i dont know why
  // $("#streamer").innerHTML = "Streamer: " + e.data;
};

function resize(){
  $("#nextButton").css("width", window.innerWidth-30);
}