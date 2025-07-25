import { useEffect, useState } from "react";

type Props = {
  initTime: number;
};

export default function Clock({ initTime }: Props) {
  const [seconds, setSeconds] = useState(initTime);
  useEffect(() => {
    const interval = setInterval(() => setSeconds((prev) => prev + 1), 1000);
    return () => clearInterval(interval);
  }, []);

  return new Date(seconds * 1000).toLocaleTimeString();
}
