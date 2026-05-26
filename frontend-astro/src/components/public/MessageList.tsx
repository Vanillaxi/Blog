import type { Comment, GuestbookMessage } from "@/lib/types";
import { publicMetaParts } from "@/lib/clientInfo";
import { MarkdownContent } from "./MarkdownContent";
import { PublicAvatar } from "./PublicAvatar";

type MessageItem = Comment | GuestbookMessage;

export function MessageList({ items, onReply }: { items: MessageItem[]; onReply?: (item: MessageItem) => void }) {
  const visibleItems = items.filter((item) => !item.deleted);

  if (import.meta.env.DEV) {
    visibleItems.flatMap((item) => item.children ?? []).forEach((reply) => {
      console.log("reply item:", reply);
    });
  }

  return (
    <section className="overflow-hidden rounded-[1.75rem] border border-white/55 bg-white/55 shadow-[0_18px_50px_rgba(148,163,184,0.14)] backdrop-blur-xl">
      {visibleItems.map((item, index) => {
        const metaParts = publicMetaParts(item);

        return (
          <article key={item.id} className={`flex gap-4 px-5 py-5 sm:px-6 ${index > 0 ? "border-t border-slate-200/40" : ""}`}>
            <PublicAvatar name={item.nickname} avatar={item.avatar} />

            <div className="min-w-0 flex-1">
              <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                <div className="min-w-0">
                  <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
                    {item.website ? (
                      <a href={item.website} target="_blank" rel="noopener noreferrer" className="font-medium text-[#24314f] transition hover:text-rose-500">
                        {item.nickname}
                      </a>
                    ) : (
                      <span className="font-medium text-[#24314f]">{item.nickname}</span>
                    )}
                  </div>

                  <div className="mt-1.5 text-xs text-slate-400">
                    <span>📍 {metaParts.join(" · ")}</span>
                  </div>
                </div>

                <time className="shrink-0 text-xs text-slate-400">{item.createdAt}</time>
              </div>

              <MarkdownContent content={item.content} className="text-sm leading-7 text-slate-600" />
              <div className="mt-3">
                <button type="button" onClick={() => onReply?.(item)} className="rounded-full bg-white/50 px-3 py-1 text-xs text-slate-500 transition hover:bg-white/80 hover:text-rose-500">
                  回复
                </button>
              </div>
              {item.children && item.children.length > 0 && (
                <div className="mt-4 space-y-3 sm:ml-10">
                  {item.children.filter((child) => !child.deleted).map((child) => (
                    <ReplyCard key={child.id} item={child} fallbackReplyTo={String(child.parentId) === String(item.id) ? item.nickname : ""} onReply={onReply} />
                  ))}
                </div>
              )}
            </div>
          </article>
        );
      })}
    </section>
  );
}

function ReplyCard({ item, fallbackReplyTo, onReply }: { item: MessageItem; fallbackReplyTo: string; onReply?: (item: MessageItem) => void }) {
  const replyToName = item.replyTo?.nickname || fallbackReplyTo;
  const metaParts = publicMetaParts(item);

  return (
    <article className="flex gap-3 rounded-2xl border border-white/45 bg-white/40 px-4 py-3">
      <PublicAvatar name={item.nickname} avatar={item.avatar} size="sm" />
      <div className="min-w-0 flex-1">
        <div className="flex flex-wrap items-center justify-between gap-2 text-xs text-slate-400">
          <div className="flex min-w-0 flex-wrap items-center gap-1">
            <span className="font-medium text-[#24314f]">{item.nickname}</span>
            {replyToName && (
              <>
                <span>回复</span>
                <span className="font-medium text-slate-500">@{replyToName}</span>
              </>
            )}
          </div>
          <time>{item.createdAt}</time>
        </div>
        <div className="mt-1 text-xs text-slate-400">
          <span>📍 {metaParts.join(" · ")}</span>
        </div>
        <MarkdownContent content={item.content} className="text-sm leading-6 text-slate-600" />
        <button type="button" onClick={() => onReply?.(item)} className="mt-2 rounded-full bg-white/50 px-3 py-1 text-xs text-slate-500 transition hover:bg-white/80 hover:text-rose-500">
          回复
        </button>
      </div>
    </article>
  );
}
