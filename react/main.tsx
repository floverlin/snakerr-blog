import React from "react";
import { createRoot } from "react-dom/client";

import LikePannel from "./components/LikePannel";
import Clock from "./components/Clock";
import Previever from "./components/Previever";
import Chat from "./components/Chat";
import FollowButton from "./components/Follow";

import getWS from "./ws";
import CommentButton from "./components/Comment";

document.addEventListener("DOMContentLoaded", () => {
  const ws = getWS();

  const likes = document.querySelectorAll<HTMLElement>(".__island-like-pannel");
  likes.forEach((isl) => {
    const props = {
      id: Number(isl.dataset.id),
      liked: isl.dataset.liked === "true",
      likes: Number(isl.dataset.likes),
      disliked: isl.dataset.disliked === "true",
      dislikes: Number(isl.dataset.dislikes),
    };
    const root = createRoot(isl);
    root.render(<LikePannel {...props} />);
  });

  const clocks = document.querySelectorAll<HTMLElement>(".__island-clock");
  clocks.forEach((isl) => {
    const props = {
      initTime: Number(isl.dataset.initTime),
    };
    const root = createRoot(isl);
    root.render(<Clock {...props} />);
  });

  const previevers = document.querySelectorAll<HTMLElement>(
    ".__island-previever"
  );
  previevers.forEach((isl) => {
    const root = createRoot(isl);
    root.render(<Previever />);
  });

  const chats = document.querySelectorAll<HTMLElement>(".__island-chat");
  chats.forEach((isl) => {
    const props = { ws };
    const root = createRoot(isl);
    root.render(<Chat {...props} />);
  });

  const follows = document.querySelectorAll<HTMLElement>(".__island-follow");
  follows.forEach((isl) => {
    const props = { followed: isl.dataset.followed === "true" };
    const root = createRoot(isl);
    root.render(<FollowButton {...props} />);
  });

  const comments = document.querySelectorAll<HTMLElement>(".__island-comment");
  comments.forEach((isl) => {
    const props = {
      postID: Number(isl.dataset.postId),
      count: Number(isl.dataset.count),
    };
    const root = createRoot(isl);
    root.render(<CommentButton {...props} />);
  });
});
