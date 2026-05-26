import { useCallback, useEffect, useState } from "react";
import { getComments } from "@/api/comment";
import type { Comment } from "@/lib/types";
import { buildTwoLevelCommentTree } from "@/lib/commentTree";
import { getListPayload, mapComment } from "@/lib/publicMappers";
import { CommentForm } from "./CommentForm";
import { MessageList } from "./MessageList";

export function CommentSection({ articleId }: { articleId: number }) {
  const [items, setItems] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [replyTo, setReplyTo] = useState<{ id: number; nickname: string } | null>(null);

  const loadComments = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const response = await getComments({
        target_type: 1,
        target_id: articleId,
        page: 1,
        pageSize: 50,
      });
      setItems(buildTwoLevelCommentTree(getListPayload(response.data).map(mapComment)));
    } catch (err) {
      setError(err instanceof Error ? err.message : "评论加载失败");
    } finally {
      setLoading(false);
    }
  }, [articleId]);

  useEffect(() => {
    void loadComments();
  }, [loadComments]);

  return (
    <section className="mx-auto mt-8 grid max-w-4xl gap-6">
      <CommentForm mode="comment" targetType={1} targetId={articleId} replyTo={replyTo} onCancelReply={() => setReplyTo(null)} onSubmitted={loadComments} />
      {error && <p className="rounded-2xl border border-rose-100 bg-rose-50/80 px-4 py-3 text-sm text-rose-600">{error}</p>}
      {loading ? (
        <p className="rounded-2xl border border-white/50 bg-white/40 px-4 py-3 text-sm text-slate-500">评论加载中...</p>
      ) : (
        <MessageList items={items} onReply={(item) => setReplyTo({ id: Number(item.id), nickname: item.nickname })} />
      )}
    </section>
  );
}
