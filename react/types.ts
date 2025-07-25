export type MessageType = "message_to_user" | "new_message" | "info";

export type User = {
  id: number;
  username: string;
  avatar?: string;
};

export type Message = {
  type: MessageType;
  from: User;
  to: User;
  body: string;
  created_at: number;
};

export type Comment = {
  id: number;
  user: User;
  post_id: number;
  body: string;
  created_at: number;
};
