import { useCallback, useEffect, useState } from "react";
import { getGuestbookMessages } from "@/api/guestbook";
import type { GuestbookMessage } from "@/lib/types";
import { buildTwoLevelCommentTree } from "@/lib/commentTree";
import { getListPayload, mapGuestbookMessage } from "@/lib/publicMappers";
import { CommentForm } from "./CommentForm";
import { GuestbookList } from "./GuestbookList";

export function GuestbookSection() {
  const [items, setItems] = useState<GuestbookMessage[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [replyTo, setReplyTo] = useState<{ id: string; nickname: string } | null>(null);

  const loadMessages = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const response = await getGuestbookMessages(1, 50);
      setItems(buildTwoLevelCommentTree(getListPayload(response.data).map(mapGuestbookMessage)));
    } catch (err) {
      setError(err instanceof Error ? err.message : "留言加载失败");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadMessages();
  }, [loadMessages]);

  return (
    <div className="mx-auto grid max-w-4xl gap-6">
      <CommentForm mode="guestbook" targetType={2} targetId={0} replyTo={replyTo} onCancelReply={() => setReplyTo(null)} onSubmitted={loadMessages} />
      {error && <p className="rounded-2xl border border-rose-100 bg-rose-50/80 px-4 py-3 text-sm text-rose-600">{error}</p>}
      {loading ? (
        <p className="rounded-2xl border border-white/50 bg-white/40 px-4 py-3 text-sm text-slate-500">留言加载中...</p>
      ) : (
        <GuestbookList items={items} onReply={(item) => setReplyTo({ id: item.id, nickname: item.nickname })} />
      )}
    </div>
  );
}
