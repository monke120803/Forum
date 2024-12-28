let socket;

const createWebsocket = () => {
  return new WebSocket("ws://localhost:8080/ws");
};

function showRegistrationUI() {
  document.querySelector("#registration").style.display = "flex";
  document.querySelector("#login").style.display = "none";
}

function showLoginUI() {
  document.querySelector("#login").style.display = "flex";
  document.querySelector("#registration").style.display = "none";
}

function showChat() {
  document.querySelector("#chat").style.display = "block";
  document.querySelector("#login").style.display = "none";
  document.querySelector("#registration").style.display = "none";
}

var conn;
var msg = document.getElementById("msg");
var log = document.getElementById("log");

function appendLog(item) {
  var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
  log.appendChild(item);
  if (doScroll) {
    log.scrollTop = log.scrollHeight - log.clientHeight;
  }
}

document.getElementById("form").onsubmit = function () {
  // console.log("form");
  if (!socket) {
    return false;
  }
  if (!msg.value) {
    return false;
  }

  var newDate = new Date();
  let sendername = document.querySelector(".username").textContent;
  let userID = document
    .querySelector("#chat > div.profile-header > div > p")
    .getAttribute("data-reciverid");
  let senderID = document.querySelector("p.username").dataset.userid;
  let msgObj = {
    message: msg.value,
    messagesender: sendername,
    messagerecipient: userID,
    SenderID: senderID,
    createdAt: newDate,
  };
  // console.log("Our Message Object: ", msgObj);
  // console.log("item from form onsubmit");
  document.querySelector("#log").innerHTML += `
  <div class="bubbleWrapper">
    <div class="inlineContainer own">
      <div class="ownBubble own">
        ${msgObj.message}
      </div>
    </div>
    <span class="own"> 
    ${newDate.getHours()}:${newDate.getMinutes()}
    </span>
  </div>`;
  socket.send(JSON.stringify(msgObj));

  msg.value = "";
  return false;
};
function showFeed() {
  document.querySelector(".auth-container").style.display = "none";
  document.querySelector("main").style.display = "block";
  socket = createWebsocket();
  socket.onopen = () => {
    console.log("Socket open", socket);
  };
  socket.onmessage = function (evt) {
    var messages = JSON.parse(evt.data);
    // console.log("Notification:", messages.notification);
    console.log("receiving msg", messages);
    let receiver = messages.messagesender;
    receiver = receiver.replace("@", "#");
    let rec = document.querySelector(receiver);
    let userlist = document.querySelector(".user-prompt");
    userlist.insertBefore(rec.parentElement, userlist.firstChild);

    //////
    let notification = document.querySelector(
      "#" + messages.messagesender.slice(1) + "-notification"
    );
    notification.classList.add("-newNotification");
    let chatReceiver = document.querySelector("#chatReceiver");
    // console.log(chatReceiver);
    let allUsersInList = Array.from(document.getElementById("recently-joined"));
    console.log("ALL USERS IN LIST: ", allUsersInList);

    var item = document.querySelector("#log");
    chatDiv = document.querySelector(".chat");

    // console.log(
    //   "xxxxx",
    //   messages.messagesender.slice(1),
    //   chatReceiver.innerHTML
    // );
    // console.log(messages.messagesender.slice(1) === chatReceiver.innerHTML);

    let onlineusers = Array.from(document.querySelectorAll(".notification"));
    // console.log("online users", onlineusers);
    for (let i = 0; i < onlineusers.length; i++) {
      // console.log(onlineusers[i].messages);
      if (
        onlineusers[i].id ==
        messages.messagesender.slice(1) + "-notification"
      ) {
        let isChatOpen =
          document.querySelector(".chat").style.display === "flex"
            ? true
            : false;
        if (
          !isChatOpen ||
          chatReceiver.innerHTML !== messages.messagesender.slice(1)
        ) {
          // notification.classList.remove("-newNotification");
          console.log(
            chatReceiver.innerHTML,
            messages.messagesender.slice(1),
            chatReceiver.innerHTML === messages.messagesender.slice(1)
          );
        } else {
          notification.classList.remove("-newNotification");
          console.log(
            chatReceiver.innerHTML,
            messages.messagesender.slice(1),
            chatReceiver.innerHTML === messages.messagesender.slice(1)
          );
        }
      }
    }

    // if (messages.messagesender.slice(1) === chatReceiver.innerHTML) {
    //   item.innerHTML +=
    //     `  <div class="bubbleWrapper">
    //       <div class="inlineContainer own">
    //        <div class="ownBubble own">${messages.message} </div>
    //     </div><span class="own">` +
    //     convertTime(messages.createdAt) +
    //     `</span>
    //     </div>`;
    // }
    if (messages.messagesender.slice(1) === chatReceiver.innerHTML) {
      item.innerHTML += `
        <div class="bubbleWrapper">
          <div class="inlineContainer">
              <div class="otherBubble other">
               ${messages.message}
              </div>
          </div>
          <span class="other">
            ${convertTime(messages.createdAt)}
          </span>
        </div>`;
    } else {
      //send notification to the backend
      //Display notrification on specific chat
    }
    // appendLog(item);
  };
}
