const en = {
  send: "send",
  writeMessage: "write message",
  writeComment: "write comment",
  loading: "loading...",
  noMessages: "no messages",
  noComments: "no comments",
  comments: "comments",
  follow: "follow",
  unfollow: "unfollow",
  upload: "upload",
};

const ru = {
  send: "отправить",
  writeMessage: "напишите сообщение",
  writeComment: "напишите комментарий",
  loading: "загрузка...",
  noMessages: "нет сообщений",
  noComments: "нет комментариев",
  comments: "комментарии",
  follow: "подписаться",
  unfollow: "отписаться",
  upload: "загрузить",
};

// AI
const ja = {
  send: "送信",
  writeMessage: "メッセージを書く",
  writeComment: "コメントを書く",
  loading: "読み込み中...",
  noMessages: "メッセージがありません",
  noComments: "コメントがありません",
  comments: "コメント",
  follow: "フォロー",
  unfollow: "フォロー解除",
  upload: "アップロード",
};

// AI
const de = {
  send: "senden",
  writeMessage: "Nachricht schreiben",
  writeComment: "Kommentar schreiben",
  loading: "lädt...",
  noMessages: "keine Nachrichten",
  noComments: "keine Kommentare",
  comments: "Kommentare",
  follow: "folgen",
  unfollow: "nicht mehr folgen",
  upload: "hochladen",
};

type LocaleKey = keyof typeof en;

export default function i18n(text: LocaleKey, c: boolean = false): string {
  const lang = getLang();
  let locale: typeof en;
  switch (lang) {
    case "ru":
      locale = ru;
      break;
    case "en":
      locale = en;
      break;
    case "ja":
      locale = ja;
      break;
    case "de":
      locale = de;
      break;
    default:
      locale = en;
      break;
  }
  if (text in locale) {
    let result = locale[text];
    if (c) {
      return capitalize(result);
    }
    return result;
  }
  console.error(`translation of [${text}]`);
  return "translation error";
}

function capitalize(text: string) {
  if (text === "") {
    return "";
  }
  return text.charAt(0).toUpperCase() + text.slice(1);
}

function getLang(): string | null {
  const cookies = document.cookie.split(";");
  for (let cookie of cookies) {
    const [cookieName, cookieValue] = cookie.split("=").map((c) => c.trim());
    if (cookieName === "lang") {
      return decodeURIComponent(cookieValue);
    }
  }
  return null;
}
