function smoking() {
  const chat = document.getElementById("smoke_chat") as HTMLUListElement;
  const chatWindow = document.getElementById(
    "smoke_chat-window"
  ) as HTMLDivElement;
  const connectionStatus = document.getElementById(
    "smoke_status"
  ) as HTMLSpanElement;
  const input = document.getElementById("smoke_input") as HTMLInputElement;
  const send = document.getElementById("smoke_send") as HTMLButtonElement;
  const disconnect = document.getElementById(
    "smoke_disconnect"
  ) as HTMLButtonElement;
  const connect = document.getElementById("smoke_connect") as HTMLButtonElement;

  type Message = {
    username: string;
    user_id: number;
    body: string;
  };
  let ws: WebSocket | null;
  let messages: Message[] = [];

  window.onbeforeunload = () => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.close(1001, "");
    }
  };

  function createConnection() {
    if (ws?.readyState === WebSocket.OPEN) return;
    ws = new WebSocket("ws/smoking");
    ws.onopen = () => {
      connectionStatus.innerText = connectionStatus.dataset.connected as string;
      connectionStatus.style = "color: green;";
    };
    ws.onmessage = (e) => {
      if (e.data === "ping") {
        if (ws?.readyState === WebSocket.OPEN) ws.send("pong");
        return;
      }
      const data = JSON.parse(e.data) as Message;
      messages.push(data);
      if (messages.length > 10) {
        messages.shift();
      }
      chat.innerHTML = messages
        .map((m) => {
          return `<li class="message">${
            m.user_id > 0
              ? `<a class="hover:underline" href="/user/${m.user_id}">${m.username}</a>`
              : `${m.username}`
          }  >>  ${
            m.body.length > 60 ? m.body.slice(0, 1024) + "..." : m.body
          }</li>`;
        })
        .join("\n");
      chatWindow.scrollTop = chat.scrollHeight;
    };
    ws.onerror = () => {
      console.error("websocket error");
    };
    ws.onclose = () => {
      connectionStatus.innerText = connectionStatus.dataset
        .disconnected as string;
      connectionStatus.style = "color: red;";
    };
  }

  function sendMessage() {
    if (ws?.readyState === WebSocket.OPEN && input.value) {
      ws.send(input.value);
      input.value = "";
      input.focus();
    }
  }

  input.onkeydown = (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      sendMessage();
    }
  };

  send.onclick = sendMessage;

  disconnect.onclick = () => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.close(1000, "");
      chat.innerHTML = `<li class="message">
          <span>${chat.dataset.noMessages}</span>
        </li>`;
      messages = [];
    }
  };

  connect.onclick = createConnection;
}

smoking();
