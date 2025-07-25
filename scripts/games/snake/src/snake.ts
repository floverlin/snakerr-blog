const canvasElement = document.getElementById("canvas") as HTMLCanvasElement;
const ctx = canvasElement.getContext("2d") as CanvasRenderingContext2D;
const scoreElement = document.getElementById("score") as HTMLDivElement;
const speedElement = document.getElementById("speed") as HTMLDivElement;
const personalBestElement = document.getElementById(
  "personal-best"
) as HTMLDivElement;
const globalBestElement = document.getElementById(
  "global-best"
) as HTMLDivElement;
const inputElement = document.querySelector("input") as HTMLInputElement;

const SIZE = 10;
const START_SPEED = 2;
const SPEED_INCREASE = 0.2;
const SENSIVITY = 20;
const PROGRESS_SCORE = 20;
const CROSSFADE_TIME = 4;

const PATH = "/static/games/snake/";
const CSRF_TOKEN = inputElement.value;

let blockVal =
  window.innerWidth >= 1024
    ? Math.floor(768 / SIZE)
    : Math.floor((window.innerWidth - window.innerWidth / 4) / SIZE);
const BLOCK_SIZE = blockVal;

let startPos: Index | null = null;
let newPos: Index | null = null;

function isIndex(val: any): val is Index {
  return val?.x !== undefined && val?.y !== undefined;
}

function setTouch() {
  canvasElement.addEventListener(
    "touchstart",
    (e) => {
      // e.preventDefault();
      startPos = {
        x: e.changedTouches[0].clientX,
        y: e.changedTouches[0].clientY,
      };
    }
    // { passive: false }
  );

  canvasElement.addEventListener(
    "touchmove",
    (e) => {
      e.preventDefault();
      newPos = {
        x: e.changedTouches[0].clientX,
        y: e.changedTouches[0].clientY,
      };
    },
    { passive: false }
  );

  canvasElement.addEventListener("touchend", () => {
    const dir = checkSwipe();
    if (dir) direction = dir;
    startPos = null;
    newPos = null;
  });

  function checkSwipe(): Direction | null {
    if (!isIndex(newPos) || !isIndex(startPos)) {
      return null;
    }
    const absMove = Math.sqrt(
      Math.abs(newPos.x - startPos.x) ** 2 +
        Math.abs(newPos.y - startPos.y) ** 2
    );
    if (absMove < SENSIVITY) return null;
    if (Math.abs(newPos.x - startPos.x) > Math.abs(newPos.y - startPos.y)) {
      if (newPos.x - startPos.x > 0) {
        return "right";
      } else {
        return "left";
      }
    } else {
      if (newPos.y - startPos.y > 0) {
        return "down";
      } else {
        return "up";
      }
    }
  }
}

function setKeys() {
  window.addEventListener(
    "keydown",
    (event) => {
      if (event.key === "ArrowUp" || event.key === "ArrowDown") {
        event.preventDefault();
      }
    },
    { passive: false }
  );

  document.addEventListener("keydown", (event) => {
    switch (event.key.toLowerCase()) {
      case "arrowup":
        direction = "up";
        break;
      case "arrowdown":
        direction = "down";
        break;
      case "arrowleft":
        direction = "left";
        break;
      case "arrowright":
        direction = "right";
        break;

      case "w":
        direction = "up";
        break;
      case "s":
        direction = "down";
        break;
      case "a":
        direction = "left";
        break;
      case "d":
        direction = "right";
        break;

      case "ц":
        direction = "up";
        break;
      case "ы":
        direction = "down";
        break;
      case "ф":
        direction = "left";
        break;
      case "в":
        direction = "right";
        break;
    }
  });
}

const BEST_SCORE = "best_score";

canvasElement.width = SIZE * BLOCK_SIZE;
canvasElement.height = SIZE * BLOCK_SIZE;

function timerDeco(f: (...args: any[]) => any): (...args: any[]) => any {
  return function (...args: any[]): any {
    const start = performance.now();
    const res = f(...args);
    console.log(performance.now() - start);
    return res;
  };
}

function createFood(): Index {
  const free = [];
  for (let x = 0; x < canvasElement.width; x += BLOCK_SIZE) {
    for (let y = 0; y < canvasElement.height; y += BLOCK_SIZE) {
      if (
        snake.some((segment) => segment.x === x && segment.y === y) ||
        (x === food?.x && y === food?.y)
      ) {
        continue;
      }
      free.push({ x, y });
    }
  }
  return free[Math.floor(Math.random() * free.length)];
}

function hasValue(val: string | null): val is string {
  return typeof val === "string";
}

enum Music {
  background,
  progress,
  death,
  eat,
  puff,
}
type Direction = "right" | "left" | "up" | "down";
type Index = { x: number; y: number };

async function setInitValues() {
  const ss = Math.floor(SIZE / 2) * BLOCK_SIZE - BLOCK_SIZE;

  direction = "right";
  oldDirection = "right";

  score = 0;
  record = await getPersonalBest();
  speed = START_SPEED;

  snake = [
    {
      x: ss,
      y: ss,
    },
    {
      x: ss - BLOCK_SIZE,
      y: ss,
    },
    {
      x: ss - 2 * BLOCK_SIZE,
      y: ss,
    },
  ];
  food = createFood();
}

let snake: Index[];
let food: Index;
let score: number;
let record: number;
let direction: Direction;
let oldDirection: Direction;
let speed: number;

const altDirection = new Map<Direction, Direction>([
  ["right", "left"],
  ["left", "right"],
  ["up", "down"],
  ["down", "up"],
]);

const deathAudio = new Audio(PATH + "audio/death.mp3");

const backgroundAudio = new Audio(PATH + "audio/background.mp3");
backgroundAudio.loop = true;

const backgroundProgressAudio = new Audio(PATH + "audio/progress.mp3");
backgroundProgressAudio.loop = true;

const eatAudio = new Audio(PATH + "audio/eat.mp3");

const puffAudio = new Audio(PATH + "audio/puff.mp3");

deathAudio.volume = 1;
backgroundAudio.volume = 1;
backgroundProgressAudio.volume = 0;
eatAudio.volume = 0.6;
puffAudio.volume = 1;

const cover = new Image();
cover.src = PATH + "img/cover.jpg";

let audioCrossfadeTimeoutIds: number[] = [];

function audioCrossfade(
  decr: HTMLAudioElement,
  incr: HTMLAudioElement,
  time: number,
  steps: number,
  curve: number
): number[] {
  const ids: number[] = [];
  const volumeStep = 1 / steps;
  for (let i = 0; i < steps; i++) {
    const cf = i + 1;
    const id = setTimeout(() => {
      incr.volume = (volumeStep * cf) ** (1 / curve);
      decr.volume = (1 - volumeStep * cf) ** curve;
      // console.log(incr.volume.toPrecision(4), " - ", decr.volume.toPrecision(4));
    }, (time / steps) * i);
    ids.push(id);
  }
  return ids;
}

function isIOS() {
  return /iPad|iPhone|iPod/.test(navigator.userAgent);
}

function playBackground(bg: Music) {
  if (isIOS()) {
    switch (bg) {
      case Music.background:
        if (!deathAudio.ended) {
          deathAudio.pause();
          deathAudio.currentTime = 0;
        }
        backgroundAudio.play();
        break;
      case Music.progress:
        backgroundAudio.pause();
        backgroundAudio.currentTime = 0;
        backgroundProgressAudio.play();
        break;
      case Music.death:
        backgroundAudio.pause();
        backgroundAudio.currentTime = 0;
        backgroundProgressAudio.pause();
        backgroundProgressAudio.currentTime = 0;
        deathAudio.play();
        break;
      case Music.eat:
        if (!eatAudio.ended) {
          eatAudio.currentTime = 0;
        }
        eatAudio.play();
        break;
      case Music.puff:
        if (!puffAudio.ended) {
          puffAudio.currentTime = 0;
        }
        puffAudio.play();
        break;
    }
    return;
  }
  switch (bg) {
    case Music.background:
      if (!deathAudio.ended) {
        deathAudio.pause();
        deathAudio.currentTime = 0;
      }
      backgroundAudio.volume = 1;
      backgroundProgressAudio.volume = 0;
      backgroundAudio.play();
      backgroundProgressAudio.play();
      break;
    case Music.progress:
      audioCrossfadeTimeoutIds = audioCrossfade(
        backgroundAudio,
        backgroundProgressAudio,
        CROSSFADE_TIME * 1000,
        10,
        2
      );
      break;
    case Music.death:
      backgroundAudio.pause();
      backgroundAudio.currentTime = 0;
      backgroundProgressAudio.pause();
      backgroundProgressAudio.currentTime = 0;
      for (const id of audioCrossfadeTimeoutIds) {
        clearTimeout(id);
      }
      audioCrossfadeTimeoutIds = [];
      deathAudio.play();
      break;
    case Music.eat:
      if (!eatAudio.ended) {
        eatAudio.currentTime = 0;
      }
      eatAudio.play();
      break;
    case Music.puff:
      if (!puffAudio.ended) {
        puffAudio.currentTime = 0;
      }
      puffAudio.play();
      break;
  }
}

function drawSnakeHead(x: number, y: number) {
  ctx.fillStyle = "darkgreen";
  ctx.fillRect(x, y, BLOCK_SIZE, BLOCK_SIZE);

  const toothX = Math.floor(BLOCK_SIZE / 5);
  const toothW = Math.floor(BLOCK_SIZE / 6);
  const toothH = Math.floor(BLOCK_SIZE / 3);

  const mouthH = Math.floor(BLOCK_SIZE / 8);

  const eyeX = Math.floor(BLOCK_SIZE / 8);
  const eyeY = Math.floor(BLOCK_SIZE / 2);
  const eyeW = Math.floor(BLOCK_SIZE / 3);
  const eyeH = Math.floor(BLOCK_SIZE / 5);

  switch (direction) {
    case "right":
      ctx.fillStyle = "white";
      ctx.fillRect(x + BLOCK_SIZE, y + toothX, -toothH, toothW);
      ctx.fillRect(
        x + BLOCK_SIZE,
        y + BLOCK_SIZE - (toothX + toothW),
        -toothH,
        toothW
      );
      ctx.fillStyle = "red";
      ctx.fillRect(
        x + BLOCK_SIZE,
        y + toothX + toothW,
        -mouthH,
        BLOCK_SIZE - (toothX + toothW) * 2
      );
      ctx.fillStyle = "yellow";
      ctx.fillRect(x + BLOCK_SIZE - eyeY, y + eyeX, -eyeH, eyeW);
      ctx.fillRect(
        x + BLOCK_SIZE - eyeY,
        y + BLOCK_SIZE - (eyeX + eyeW),
        -eyeH,
        eyeW
      );
      break;
    case "left":
      ctx.fillStyle = "white";
      ctx.fillRect(x, y + toothX, toothH, toothW);
      ctx.fillRect(x, y + BLOCK_SIZE - (toothX + toothW), toothH, toothW);
      ctx.fillStyle = "red";
      ctx.fillRect(
        x,
        y + toothX + toothW,
        mouthH,
        BLOCK_SIZE - (toothX + toothW) * 2
      );
      ctx.fillStyle = "yellow";
      ctx.fillRect(x + eyeY, y + eyeX, eyeH, eyeW);
      ctx.fillRect(x + eyeY, y + BLOCK_SIZE - (eyeX + eyeW), eyeH, eyeW);
      break;
    case "up":
      ctx.fillStyle = "white";
      ctx.fillRect(x + toothX, y, toothW, toothH);
      ctx.fillRect(x + BLOCK_SIZE - (toothX + toothW), y, toothW, toothH);
      ctx.fillStyle = "red";
      ctx.fillRect(
        x + toothX + toothW,
        y,
        BLOCK_SIZE - (toothX + toothW) * 2,
        mouthH
      );
      ctx.fillStyle = "yellow";
      ctx.fillRect(x + eyeX, y + eyeY, eyeW, eyeH);
      ctx.fillRect(x + BLOCK_SIZE - (eyeX + eyeW), y + eyeY, eyeW, eyeH);
      break;
    case "down":
      ctx.fillStyle = "white";
      ctx.fillRect(x + toothX, y + BLOCK_SIZE, toothW, -toothH);
      ctx.fillRect(
        x + BLOCK_SIZE - (toothX + toothW),
        y + BLOCK_SIZE,
        toothW,
        -toothH
      );
      ctx.fillStyle = "red";
      ctx.fillRect(
        x + toothX + toothW,
        y + BLOCK_SIZE,
        BLOCK_SIZE - (toothX + toothW) * 2,
        -mouthH
      );
      ctx.fillStyle = "yellow";
      ctx.fillRect(x + eyeX, y + BLOCK_SIZE - eyeY, eyeW, -eyeH);
      ctx.fillRect(
        x + BLOCK_SIZE - (eyeX + eyeW),
        y + BLOCK_SIZE - eyeY,
        eyeW,
        -eyeH
      );
      break;
  }
}

function draw() {
  ctx.clearRect(0, 0, canvasElement.width, canvasElement.height);

  snake.forEach((segment, index) => {
    if (index === 0) {
      drawSnakeHead(segment.x, segment.y);
      return;
    }
    ctx.fillStyle = index % 2 === 0 ? "darkgreen" : "green";
    ctx.fillRect(segment.x, segment.y, BLOCK_SIZE, BLOCK_SIZE);
  });

  ctx.fillStyle = "red";
  ctx.fillRect(food.x, food.y, BLOCK_SIZE, BLOCK_SIZE);

  scoreElement.innerText = score.toString();
  speedElement.innerText = speed.toFixed(1);
}

function update(): boolean {
  let newX = snake[0].x;
  let newY = snake[0].y;

  if (direction === altDirection.get(oldDirection)) {
    direction = oldDirection;
  }

  oldDirection = direction;
  switch (direction) {
    case "right":
      newX += BLOCK_SIZE;
      break;
    case "left":
      newX -= BLOCK_SIZE;
      break;
    case "up":
      newY -= BLOCK_SIZE;
      break;
    case "down":
      newY += BLOCK_SIZE;
      break;
  }

  if (
    newX < 0 ||
    newX >= canvasElement.width ||
    newY < 0 ||
    newY >= canvasElement.height ||
    snake
      .slice(0, snake.length - 1)
      .some(
        (segment, index) => segment.x === newX && segment.y === newY && index
      )
  ) {
    return false;
  }

  if (newX === food.x && newY === food.y) {
    playBackground(Music.eat);
    score++;
    if (score === PROGRESS_SCORE) playBackground(Music.progress);
    speed += SPEED_INCREASE;
    food = createFood();
  } else {
    snake.pop();
  }

  snake.unshift({ x: newX, y: newY });

  return true;
}

function gameLoop() {
  const cont = update();
  draw();

  if (!cont) {
    playBackground(Music.death);
    if (score > record) {
      sendRecord(score);
      setBestRecords(score);
    }
    playAgainQuestion();
    return;
  } else if (snake.length === SIZE * SIZE) {
    sendRecord(score);
    playAgainQuestion();
    return;
  }
  setTimeout(gameLoop, 1000 / speed);
}

function playAgainQuestion() {
  const text = "click to play again";
  const font = "px arial";
  const center = (SIZE * BLOCK_SIZE) / 2;
  let textSize = 8;
  ctx.font = textSize + font;
  while (ctx.measureText(text).width < canvasElement.width * 0.8) {
    textSize += 4;
    ctx.font = textSize + font;
  }
  ctx.textAlign = "center";
  ctx.textBaseline = "middle";
  ctx.fillStyle = "black";
  ctx.fillText(text, center, center);

  setInitValues();
  canvasElement.addEventListener("click", playGame);
}

setKeys();
setTouch();

async function playGame() {
  canvasElement.removeEventListener("click", playGame);
  await setInitValues();
  draw();
  playBackground(Music.puff);
  playBackground(Music.background);
  gameLoop();
}

cover.onload = () => {
  canvasElement.addEventListener("click", playGame);
  ctx.drawImage(cover, 0, 0, canvasElement.width, canvasElement.height);
};

async function sendRecord(record: number) {
  const body = {
    record,
    csrf_token: CSRF_TOKEN,
  };
  const resp = await fetch("/api/snake/record", {
    method: "post",
    body: JSON.stringify(body),
  });
  if (!resp.ok) {
    console.error("record send error");
    return;
  }
  const res = await resp.json();
}

async function getPersonalBest(): Promise<number> {
  const resp = await fetch("/api/snake/personal", {
    method: "get",
  });
  if (!resp.ok) {
    console.error("record send error");
    return 0;
  }
  const res = await resp.json();
  return res.record;
}

async function setBestRecords(score: number) {
  const resp = await fetch("/api/snake/global", {
    method: "get",
  });
  if (!resp.ok) {
    console.error("record send error");
    return 0;
  }
  const res = await resp.json();
  personalBestElement.innerText = score.toString();
  globalBestElement.innerText = res.record.toString();
}
