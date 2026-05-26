import { useEffect, useMemo, useState } from "react";

const titleNames = ["Vanillaxi", "ヴァニラシ"];
const subtitles = [
  "A quiet place for writings, notes, and memories.",
  "言葉と記憶を、そっと綴る場所。",
];
const typingSpeed = 110;
const japaneseTypingSpeed = 135;
const deleteSpeed = 55;
const holdTime = 2200;

function prefersReducedMotion() {
  return typeof window !== "undefined" && window.matchMedia("(prefers-reduced-motion: reduce)").matches;
}

export function HeroCenter() {
  const [titleIndex, setTitleIndex] = useState(0);
  const [subtitleIndex, setSubtitleIndex] = useState(0);
  const [typedText, setTypedText] = useState(subtitles[0]);
  const [reducedMotion, setReducedMotion] = useState(false);

  useEffect(() => {
    setReducedMotion(prefersReducedMotion());
  }, []);

  useEffect(() => {
    if (reducedMotion) return;
    const timer = window.setInterval(() => {
      setTitleIndex((current) => (current + 1) % titleNames.length);
    }, 5600);
    return () => window.clearInterval(timer);
  }, [reducedMotion]);

  useEffect(() => {
    if (reducedMotion) {
      setTypedText(subtitles[0]);
      return;
    }

    const text = subtitles[subtitleIndex];
    const speed = subtitleIndex === 1 ? japaneseTypingSpeed : typingSpeed;
    const timers: number[] = [];
    let position = 0;
    setTypedText("");

    const typeNext = () => {
      position += 1;
      setTypedText(text.slice(0, position));

      if (position >= text.length) {
        timers.push(
          window.setTimeout(() => {
            deleteNext();
          }, holdTime),
        );
        return;
      }

      timers.push(window.setTimeout(typeNext, speed));
    };

    const deleteNext = () => {
      position -= 1;
      setTypedText(text.slice(0, position));

      if (position <= 0) {
        timers.push(
          window.setTimeout(() => {
            setSubtitleIndex((current) => (current + 1) % subtitles.length);
          }, 260),
        );
        return;
      }

      timers.push(window.setTimeout(deleteNext, deleteSpeed));
    };

    timers.push(window.setTimeout(typeNext, 260));

    return () => {
      timers.forEach((timer) => {
        window.clearTimeout(timer);
      });
    };
  }, [subtitleIndex, reducedMotion]);

  const title = useMemo(() => titleNames[titleIndex], [titleIndex]);
  const isJapaneseTitle = title === "ヴァニラシ";

  return (
    <div className="-translate-y-10 text-center sm:-translate-y-14">
      <h1
        key={title}
        className={`hero-title-fade hero-title-soft-shadow group font-serif italic text-[#3B4A6B] transition duration-700 hover:bg-gradient-to-r hover:from-[#F5A9C8] hover:via-[#F7B7A3] hover:to-[#A9D8FF] hover:text-transparent hover:[background-clip:text] hover:[-webkit-background-clip:text] ${
          isJapaneseTitle ? "text-[3.35rem] sm:text-[5.4rem]" : "text-6xl sm:text-8xl"
        }`}
      >
        {title}
      </h1>

      <div className="mx-auto mt-10 inline-flex min-w-[min(460px,88vw)] max-w-[88vw] items-center justify-center rounded-full border border-white/55 bg-white/35 px-6 py-2.5 shadow-[0_10px_30px_rgba(148,163,184,0.18)] backdrop-blur-xl sm:mt-16">
        <span className="text-center font-serif text-sm italic leading-6 text-[#3B4A6B] sm:text-base">
          {typedText}
        </span>
        {!reducedMotion && <span className="typing-caret ml-1 h-4 w-px bg-rose-400" aria-hidden="true" />}
      </div>
    </div>
  );
}
