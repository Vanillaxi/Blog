import { useEffect, useState } from "react";
import { getFriendLinks } from "@/api/friend";
import type { FriendLink } from "@/lib/types";
import { getListPayload, mapFriendLink } from "@/lib/publicMappers";
import { ExternalLink } from "lucide-react";

const defaultDescription = "这个人还没有留下介绍。";

function getDisplayUrl(url: string) {
  try {
    const parsed = new URL(url);
    return parsed.hostname.replace(/^www\./, "");
  } catch {
    return url;
  }
}

function FriendAvatar({ friend }: { friend: FriendLink }) {
  const [imageFailed, setImageFailed] = useState(false);
  const initial = friend.name.trim().slice(0, 1).toUpperCase() || "?";

  if (!friend.logo || imageFailed) {
    return (
      <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-[#24314f] text-base font-semibold text-white shadow-[0_10px_28px_rgba(36,49,79,0.18)]">
        {initial}
      </div>
    );
  }

  return (
    <img
      src={friend.logo}
      alt={`${friend.name} logo`}
      loading="lazy"
      onError={() => setImageFailed(true)}
      className="h-14 w-14 shrink-0 rounded-full border border-white/70 bg-white/70 object-cover shadow-[0_10px_28px_rgba(148,163,184,0.16)]"
    />
  );
}

function FriendCard({ friend }: { friend: FriendLink }) {
  return (
    <a
      href={friend.url}
      target="_blank"
      rel="noopener noreferrer"
      className="group relative flex min-w-0 gap-4 rounded-[1.35rem] border border-white/60 bg-white/50 p-4 shadow-[0_18px_48px_rgba(100,116,139,0.10)] backdrop-blur-xl transition duration-200 hover:-translate-y-1 hover:border-white/80 hover:bg-white/64 hover:shadow-[0_22px_54px_rgba(100,116,139,0.14)] sm:p-5"
    >
      <FriendAvatar friend={friend} />
      <div className="min-w-0 flex-1 pr-7">
        <h2 className="truncate text-[0.98rem] font-semibold text-[#24314f]">{friend.name}</h2>
        <p className="mt-1.5 line-clamp-2 text-sm leading-6 text-slate-500">{friend.description || defaultDescription}</p>
        <p className="mt-2 truncate text-xs text-slate-400">{getDisplayUrl(friend.url)}</p>
      </div>
      <span className="absolute right-4 top-4 flex h-8 w-8 items-center justify-center rounded-full border border-white/60 bg-white/54 text-slate-400 transition group-hover:text-rose-400">
        <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
        <span className="sr-only">打开 {friend.name}</span>
      </span>
    </a>
  );
}

export function FriendLinks() {
  const [items, setItems] = useState<FriendLink[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let ignore = false;

    async function loadFriendLinks() {
      setLoading(true);
      setError("");
      try {
        const response = await getFriendLinks({ page: 1, pageSize: 50 });
        const links = getListPayload(response.data).map(mapFriendLink);
        if (!ignore) setItems(links);
      } catch (err) {
        console.error("[friends] parse or request failed", err);
        if (!ignore) setError(err instanceof Error ? err.message : "友链加载失败");
      } finally {
        if (!ignore) setLoading(false);
      }
    }

    void loadFriendLinks();
    return () => {
      ignore = true;
    };
  }, []);

  if (loading) {
    return <p className="rounded-2xl border border-white/50 bg-white/40 px-4 py-3 text-sm text-slate-500">友链加载中...</p>;
  }

  if (error) {
    return <p className="rounded-2xl border border-rose-100 bg-rose-50/80 px-4 py-3 text-sm text-rose-600">{error}</p>;
  }

  if (items.length === 0) {
    return <p className="rounded-2xl border border-white/50 bg-white/40 px-4 py-3 text-sm text-slate-500">暂时还没有友链。</p>;
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
      {items.map((friend) => (
        <FriendCard friend={friend} key={friend.id} />
      ))}
    </div>
  );
}
