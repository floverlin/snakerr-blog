import React, { useState } from "react";
import i18n from "../locales";

type Props = {
  followed: boolean;
};

function getUserID() {
  return Number(document.location.pathname.split("/")[2]);
}

export default function FollowButton({ followed }: Props) {
  const [isFollowed, setIsFollowed] = useState(followed);

  async function toggleFollow() {
    const data = {
      type: isFollowed ? "unfollow" : "follow",
      to: getUserID(),
    };
    const resp = await fetch("/api/follow", {
      method: "post",
      body: JSON.stringify(data),
    });
    if (!resp.ok) {
      const res = await resp.json();
      console.error(`follow error: ${res.error}`);
      return;
    }
    setIsFollowed((prev) => !prev);
  }

  return (
    <button
      onClick={() => toggleFollow()}
      className="btn px-2 py-1 text-base font-sans"
    >
      {isFollowed ? i18n("unfollow", true) : i18n("follow", true)}
    </button>
  );
}
