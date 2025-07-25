import React, { useState } from "react";
import i18n from "../locales";

export default function Previever() {
  const [imgURL, setImgURL] = useState<string | null>(null);

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0] || null;
    if (file) {
      const url = URL.createObjectURL(file);
      setImgURL(url);
    } else {
      setImgURL(null);
    }
  }

  return (
    <div className="flex flex-col gap-8 items-center">
      <label className="btn md:py-2 md:text-xl text-center w-full">
        {i18n("upload", true)}
        <input
          className="hidden"
          accept="image/*"
          type="file"
          name="avatar"
          onChange={(e) => handleFileChange(e)}
        />
      </label>
      {imgURL && (
        <div className="w-48 h-48 overflow-hidden border-4 border-black rounded-md">
          <img
            className="w-full h-full object-cover"
            src={imgURL}
            alt="uploaded image"
          />
        </div>
      )}
    </div>
  );
}
