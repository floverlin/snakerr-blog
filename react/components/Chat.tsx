import React, { useEffect, useRef, useState } from "react";
import { Message } from "../types";
import { getDialogistID, getID, haveUnreaded } from "../utils";
import { defaultReactOnMessage } from "../ws";
import i18n from "../locales";
import Input from "./Input";

type Props = {
  ws: WebSocket;
};

function messageShorter(message: string, length: number) {
  return message.length > length ? message.slice(0, length) + "..." : message;
}

export default function Chat({ ws }: Props) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const chatWindowRef = useRef<HTMLDivElement | null>(null);

  function scrollChat() {
    if (chatWindowRef.current) {
      chatWindowRef.current.scrollTop = chatWindowRef.current.scrollHeight;
    }
  }

  useEffect(() => {
    async function getInitialMessages() {
      const resp = await fetch("/api/chat/messages", {
        method: "post",
        body: JSON.stringify({
          offset: 0,
          limit: 40,
          dialogist_id: getDialogistID(),
        }),
      });
      if (!resp.ok) {
        console.error(`fetch init messages: ${resp.status} ${resp.statusText}`);
        return;
      }
      const result = await resp.json();
      setMessages(result);
      setLoading(false);
    }

    async function markChatReaded() {
      type Resp = {
        all_readed: boolean;
      };
      const resp = await fetch("/api/chat/mark_readed", {
        method: "post",
        body: JSON.stringify({
          me: getID(),
          dialogist: getDialogistID(),
        }),
      });
      if (!resp.ok) {
        console.error(`mark chat readed: ${resp.status} ${resp.statusText}`);
      }
      const data = (await resp.json()) as Resp;
      if (data.all_readed) {
        haveUnreaded(false);
      }
    }

    getInitialMessages();
    markChatReaded();
  }, []);

  useEffect(() => {
    const prevHandler = ws.onmessage;
    ws.onmessage = (e: MessageEvent) => {
      const defaultData = defaultReactOnMessage(e);
      if (defaultData) {
        ws.send(defaultData);
        return;
      }
      const data = JSON.parse(e.data) as Message;
      if (data.type === "new_message") {
        setMessages((prev) => [...prev, data]);
      }
    };
    return () => {
      ws.onmessage = prevHandler;
    };
  }, [ws]);

  useEffect(() => {
    scrollChat();
  }, [messages]);

  function sendMessage(text: string) {
    const message: Message = {
      type: "message_to_user",
      from: { id: getID(), username: "" },
      to: { id: getDialogistID(), username: "" },
      body: text,
      created_at: Math.floor(Date.now() / 1000),
    };
    if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify(message));
  }

  return (
    <>
      <div
        className="border-4 border-black rounded-xl w-full h-96 lg:h-124 p-4 bg-slate-600 overflow-y-scroll mb-2 md:mb-8"
        id="chat-window"
        ref={chatWindowRef}
      >
        <ul id="chat" className="space-y-2">
          {messages.length === 0 ? (
            <li className="message">
              <span>{loading ? i18n("loading") : i18n("noMessages")}</span>
            </li>
          ) : (
            messages.map((message, idx) => (
              <li className="message" key={idx}>
                <span
                  className="text-xl"
                  style={
                    message.from.id === getID()
                      ? { color: "skyblue" }
                      : { color: "cornflowerblue" }
                  }
                >
                  {message.from.username}
                </span>
                <span className="text-base text-slate-400">{" >> "}</span>
                <span>{messageShorter(message.body, 1024)}</span>
                <div className="text-right text-[12px] text-slate-400 font-sans italic">
                  {new Date(message.created_at * 1000).toLocaleString()}
                </div>
              </li>
            ))
          )}
        </ul>
      </div>
      <Input placeholder={i18n("writeMessage")} func={sendMessage} />
    </>
  );
}
