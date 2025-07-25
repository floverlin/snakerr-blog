import React, { useEffect, useRef, useState } from "react";
import { Comment } from "../types";
import Input from "./Input";
import { getID, getUsername } from "../utils";
import i18n from "../locales";

type Props = {
  postID: number;
  count: number;
};

export default function CommentButton({ postID, count }: Props) {
  const [isOpened, setIsOpened] = useState(false);
  const [comments, setComments] = useState<Comment[]>([]);
  const [placeholder, setPlaceholder] = useState<string>(i18n("loading"));
  const windowRef = useRef<HTMLUListElement>(null);
  const commentClass = "text-left w-full bg-gray-300 rounded-md p-2";

  function scrollChat() {
    if (windowRef.current) {
      windowRef.current.scrollTop = windowRef.current.scrollHeight;
    }
  }

  async function sendComment(text: string) {
    const comment: Comment = {
      id: 0,
      body: text,
      post_id: postID,
      user: { id: getID(), username: getUsername() },
      created_at: Math.floor(Date.now() / 1000),
    };
    const resp = await fetch("/api/comments", {
      method: "post",
      body: JSON.stringify(comment),
    });
    if (!resp.ok) {
      const res = await resp.json();
      console.error(
        `send comment error: ${resp.status} ${resp.statusText} ${res.error}`
      );
      return;
    }
    setComments((prev) => [...prev, comment]);
  }

  useEffect(() => {
    scrollChat();
  }, [comments]);

  useEffect(() => {
    if (!isOpened) {
      setPlaceholder(i18n("loading"));
      setComments([]);
      return;
    }
    async function getComments() {
      const resp = await fetch(`/api/comments/${postID}`, {
        method: "get",
      });
      if (!resp.ok) {
        console.error(
          `fetch comments error: ${resp.status} ${resp.statusText}`
        );
        return;
      }
      const res = (await resp.json()) as Comment[];
      console.debug(res);
      setComments(res);
      setPlaceholder(i18n("noComments"));
    }
    getComments();
  }, [isOpened]);

  return (
    <>
      <span
        className="hover:underline hover:cursor-pointer truncate"
        onClick={() => {
          setIsOpened(true);
        }}
      >
        {i18n("comments") + ` [${count}]`}
      </span>
      {isOpened && (
        <>
          <div className="fixed inset-0 bg-black opacity-80 z-40" />
          <div className="fixed inset-0 flex items-center justify-center z-50">
            <div className="bg-gray-200 rounded-xl py-6 px-12 w-full lg:w-200 relative">
              <button
                className="absolute top-2 right-4 text-2xl font-bold text-black hover:text-red-600 hover:cursor-pointer"
                onClick={() => setIsOpened(false)}
                aria-label="close"
              >
                x
              </button>
              <div className="flex flex-col items-center w-full">
                <ul
                  ref={windowRef}
                  className="mb-2 md:mb-8 h-60 md:h-120 overflow-y-auto overflow-x-hidden border-y-2 border-black border-dashed p-4 w-full space-y-2 md:text-xl"
                >
                  {comments.length > 0 ? (
                    comments.map((comment) => (
                      <li className={commentClass} key={comment.id}>
                        <a
                          className="text-blue-600 underline truncate"
                          href={`/user/${comment.user.id}`}
                        >
                          {comment.user.username}
                        </a>
                        <span>{" >> "}</span>
                        <span>{comment.body}</span>
                        <div className="text-[10px] text-right">
                          {new Date(comment.created_at * 1000).toLocaleString()}
                        </div>
                      </li>
                    ))
                  ) : (
                    <li className={commentClass}>{placeholder}</li>
                  )}
                </ul>
                <Input placeholder={i18n("writeComment")} func={sendComment} />
              </div>
            </div>
          </div>
        </>
      )}
    </>
  );
}
