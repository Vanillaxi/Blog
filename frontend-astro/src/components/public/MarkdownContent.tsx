import type { ReactNode } from "react";

interface Props {
  content: string;
  className?: string;
}

type Block =
  | { type: "paragraph"; lines: string[] }
  | { type: "list"; items: string[] };

function isSafeLink(href: string) {
  const value = href.trim();
  if (!value) return false;
  if (value.startsWith("/") || value.startsWith("#")) return true;

  try {
    const url = new URL(value);
    return ["http:", "https:", "mailto:"].includes(url.protocol);
  } catch {
    return false;
  }
}

function parseInline(text: string, keyPrefix: string): ReactNode[] {
  const nodes: ReactNode[] = [];
  let index = 0;

  while (index < text.length) {
    const linkStart = text.indexOf("[", index);
    const boldStart = text.indexOf("**", index);
    const candidates = [linkStart, boldStart].filter((value) => value >= 0);
    const next = candidates.length > 0 ? Math.min(...candidates) : -1;

    if (next < 0) {
      nodes.push(text.slice(index));
      break;
    }

    if (next > index) {
      nodes.push(text.slice(index, next));
    }

    if (next === boldStart) {
      const end = text.indexOf("**", next + 2);
      if (end < 0) {
        nodes.push(text.slice(next));
        break;
      }
      nodes.push(
        <strong key={`${keyPrefix}-strong-${next}`} className="font-semibold text-[#24314f]">
          {parseInline(text.slice(next + 2, end), `${keyPrefix}-strong-${next}`)}
        </strong>,
      );
      index = end + 2;
      continue;
    }

    const labelEnd = text.indexOf("]", next + 1);
    const urlStart = labelEnd >= 0 ? text.indexOf("(", labelEnd + 1) : -1;
    const urlEnd = urlStart === labelEnd + 1 ? text.indexOf(")", urlStart + 1) : -1;

    if (labelEnd < 0 || urlStart !== labelEnd + 1 || urlEnd < 0) {
      nodes.push(text.slice(next, next + 1));
      index = next + 1;
      continue;
    }

    const label = text.slice(next + 1, labelEnd);
    const href = text.slice(urlStart + 1, urlEnd).trim();

    if (isSafeLink(href)) {
      nodes.push(
        <a
          key={`${keyPrefix}-link-${next}`}
          href={href}
          target="_blank"
          rel="noopener noreferrer"
          className="break-words text-rose-500 underline decoration-rose-200 underline-offset-4 transition hover:text-rose-600"
        >
          {parseInline(label, `${keyPrefix}-link-${next}`)}
        </a>,
      );
    } else {
      nodes.push(text.slice(next, urlEnd + 1));
    }

    index = urlEnd + 1;
  }

  return nodes;
}

function parseBlocks(content: string): Block[] {
  const blocks: Block[] = [];
  const lines = content.replace(/\r\n/g, "\n").split("\n");
  let paragraph: string[] = [];
  let list: string[] = [];

  function flushParagraph() {
    if (paragraph.length > 0) {
      blocks.push({ type: "paragraph", lines: paragraph });
      paragraph = [];
    }
  }

  function flushList() {
    if (list.length > 0) {
      blocks.push({ type: "list", items: list });
      list = [];
    }
  }

  for (const line of lines) {
    const listMatch = line.match(/^\s*[-*]\s+(.+)$/);

    if (!line.trim()) {
      flushParagraph();
      flushList();
      continue;
    }

    if (listMatch) {
      flushParagraph();
      list.push(listMatch[1]);
      continue;
    }

    flushList();
    paragraph.push(line);
  }

  flushParagraph();
  flushList();

  return blocks;
}

export function MarkdownContent({ content, className = "" }: Props) {
  const blocks = parseBlocks(content);

  if (blocks.length === 0) {
    return null;
  }

  return (
    <div className={className}>
      {blocks.map((block, blockIndex) => {
        if (block.type === "list") {
          return (
            <ul key={`list-${blockIndex}`} className="mt-3 list-disc space-y-1 pl-5">
              {block.items.map((item, itemIndex) => (
                <li key={`list-${blockIndex}-${itemIndex}`}>{parseInline(item, `list-${blockIndex}-${itemIndex}`)}</li>
              ))}
            </ul>
          );
        }

        return (
          <p key={`paragraph-${blockIndex}`} className="mt-3 whitespace-pre-wrap">
            {block.lines.map((line, lineIndex) => (
              <span key={`paragraph-${blockIndex}-${lineIndex}`}>
                {lineIndex > 0 && <br />}
                {parseInline(line, `paragraph-${blockIndex}-${lineIndex}`)}
              </span>
            ))}
          </p>
        );
      })}
    </div>
  );
}
