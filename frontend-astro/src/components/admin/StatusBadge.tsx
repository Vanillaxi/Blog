import type { ArticleStatus } from "@/lib/types";

const styles: Record<ArticleStatus, string> = {
  draft: "bg-slate-100 text-slate-600",
  published: "bg-emerald-50 text-emerald-700",
  offline: "bg-amber-50 text-amber-700",
  deleted: "bg-rose-50 text-rose-700",
};

const labels: Record<ArticleStatus, string> = {
  draft: "草稿",
  published: "已发布",
  offline: "已下架",
  deleted: "已删除",
};

export function StatusBadge({ status }: { status: ArticleStatus }) {
  return <span className={`rounded-full px-3 py-1 text-xs ${styles[status]}`}>{labels[status]}</span>;
}
