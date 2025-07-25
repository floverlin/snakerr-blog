import React, { useState } from "react";

type Props = {
  id: number;
  liked: boolean;
  likes: number;
  disliked: boolean;
  dislikes: number;
};

type Resp = {
  liked: boolean;
  likes: number;
  disliked: boolean;
  dislikes: number;
};

export default function LikePannel(props: Props) {
  const [isLiked, setIsLiked] = useState(props.liked);
  const [likeCount, setLikeCount] = useState(props.likes);
  const [isDisliked, setIsDisliked] = useState(props.disliked);
  const [dislikeCount, setDislikeCount] = useState(props.dislikes);

  async function toggleButton(type: "like" | "dislike") {
    const body = { post_id: props.id, type };
    const resp = await fetch("/api/like", {
      method: "post",
      body: JSON.stringify(body),
    });
    if (!resp.ok) {
      const res = await resp.json();
      console.error(res.error);
      return;
    }
    const res = (await resp.json()) as Resp;
    setIsLiked(res.liked);
    setLikeCount(res.likes);
    setIsDisliked(res.disliked);
    setDislikeCount(res.dislikes);
  }

  return (
    <div className="flex items-center">
      <Button
        func={() => toggleButton("dislike")}
        active={isDisliked}
        counter={dislikeCount}
        text="👎"
        filter="hue-rotate(340deg) saturate(400%) brightness(90%) contrast(200%)"
      />
      <span className="font-bold">|</span>
      <Button
        func={() => toggleButton("like")}
        active={isLiked}
        counter={likeCount}
        filter="hue-rotate(60deg) saturate(100%) brightness(75%) contrast(200%)"
        text="👍"
      />
    </div>
  );
}

function Button({
  func,
  active,
  counter,
  text,
  filter = "none",
}: {
  func: () => void;
  active: boolean;
  counter: number;
  text: string;
  filter?: string;
}) {
  return (
    <button
      onClick={func}
      className="w-12 text-sm whitespace-nowrap hover:cursor-pointer"
    >
      <span
        style={{
          filter: active ? filter : "saturate(800%) grayscale(100%)",
        }}
      >
        {text + " "}
      </span>
      <span>{counter}</span>
    </button>
  );
}
