export function setID(id: number) {
  localStorage.setItem("id", String(id));
}

export function getID() {
  return Number(localStorage.getItem("id"));
}
export function setUsername(username: string) {
  localStorage.setItem("username", username);
}

export function getUsername() {
  let username = localStorage.getItem("username");
  if (!username) {
    username = "xx_username_xx";
  }
  return username;
}

export function getDialogistID() {
  return Number(window.location.href.split("/")[4]);
}

export function haveUnreaded(have: boolean) {
  const chat = document.getElementById("my-chats");
  if (chat) {
    if (have) {
      chat.style.color = "red";
    } else {
      chat.style.color = "black";
    }
  }
}
