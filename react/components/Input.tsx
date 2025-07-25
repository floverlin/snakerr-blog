import React, { useRef, useState } from "react";
import i18n from "../locales";

type Props = {
  func: (text: string) => void;
  placeholder: string;
};

export default function Input({ func, placeholder }: Props) {
  const [value, setValue] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  const send = (text: string) => {
    setValue("");
    if (value.trim() === "") return;
    func(text);
    inputRef.current?.focus();
  };

  return (
    <div className="flex gap-2 md:gap-8 w-full flex-col md:flex-row">
      <input
        autoComplete="off"
        ref={inputRef}
        value={value}
        onChange={(e) => setValue(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            e.preventDefault();
            send(value);
          }
        }}
        className="w-full h-10"
        id="input"
        type="text"
        placeholder={placeholder}
      />
      <button
        className="btn px-4 py-1 ml-auto w-full md:w-32"
        id="send"
        onClick={() => {
          send(value);
        }}
      >
        {i18n("send")}
      </button>
    </div>
  );
}
