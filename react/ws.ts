import { Message } from "./types";
import { getDialogistID, getID, haveUnreaded, setID, setUsername } from "./utils";

let ws: WebSocket | null;
let alertAudio: HTMLAudioElement | null;

export default function getWS() {
  if (ws) {
    return ws;
  } else {
    ws = createWS();
    return ws;
  }
}

function createWS() {
  ws = new WebSocket("/ws/main");
  alertAudio = new Audio("/static/audio/alert.mp3");
  alertAudio.volume = 0.5;
  window.onbeforeunload = () => {
    if (ws) {
      ws.close();
    }
  };
  configureWS(ws);
  return ws;
}

function configureWS(ws: WebSocket) {
  ws.onopen = () => {
    console.debug("ws opened");
  };
  ws.onmessage = (e) => {
    const data = defaultReactOnMessage(e);
    if (data) {
      ws.send(data);
    }
  };
  ws.onclose = () => {
    console.debug("ws closed");
  };
  ws.onerror = tryReconnect;
}

function tryReconnect(e: Event) {
  console.error(`ws error: ${e}\nreconnecting after 5 seconds:`);
  for (let i = 5; i > 0; i--) {
    setTimeout(() => console.error(i), (5 - i) * 1000);
  }
  setTimeout(() => (ws = createWS()), 5 * 1000);
}

export function defaultReactOnMessage(e: MessageEvent): string | null {
  if (e.data === "ping") {
    return "pong";
  }
  const message = JSON.parse(e.data) as Message;
  switch (message.type) {
    case "info":
      setID(message.to.id);
      setUsername(message.to.username);
      return null;
    case "new_message":
      if (message.from.id !== getID()) {
        playAlert();
      }
      if (message.from.id !== getDialogistID() && message.from.id !== getID()) {
        haveUnreaded(true);
      }
      return null;
    default:
      console.error("wrong message type");
      return null;
  }
}

function playAlert() {
  if (alertAudio) {
    if (!alertAudio.ended) {
      alertAudio.pause();
      alertAudio.currentTime = 0;
    }
    alertAudio.play();
  }
}
