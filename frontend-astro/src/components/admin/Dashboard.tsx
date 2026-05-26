import { useEffect, useState } from "react";
import { getAdminDashboard } from "@/api/admin";
import { mapAdminDashboard, type AdminDashboardData } from "@/lib/adminAdapters";
import { StatusBadge } from "./StatusBadge";

const emptyDashboard: AdminDashboardData = {
  articleCount: 0,
  publishedCount: 0,
  draftCount: 0,
  offlineCount: 0,
  deletedArticleCount: 0,
  categoryCount: 0,
  tagCount: 0,
  commentCount: 0,
  guestbookCount: 0,
  friendlinkCount: 0,
  recentArticles: [],
};

export function Dashboard() {
  const [dashboard, setDashboard] = useState<AdminDashboardData>(emptyDashboard);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function loadDashboard() {
    setLoading(true);
    setError("");
    try {
      const response = await getAdminDashboard();
      setDashboard(mapAdminDashboard(response.data));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Dashboard 加载失败");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadDashboard();
  }, []);

  const stats = [
    ["文章总数", dashboard.articleCount],
    ["已发布", dashboard.publishedCount],
    ["草稿数", dashboard.draftCount],
    ["已下架", dashboard.offlineCount],
    ["回收站", dashboard.deletedArticleCount],
    ["分类数", dashboard.categoryCount],
    ["Tag 数", dashboard.tagCount],
    ["评论总数", dashboard.commentCount],
    ["留言总数", dashboard.guestbookCount],
    ["友链数", dashboard.friendlinkCount],
  ];

  return (
    <div className="space-y-5">
      {error && <p className="rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">{error}</p>}
      {loading && <p className="rounded-2xl bg-white/70 p-4 text-sm text-slate-500">加载中...</p>}
      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
        {stats.map(([label, value]) => (
          <div key={label} className="rounded-[1.5rem] border border-white/70 bg-white/70 p-5 shadow-soft backdrop-blur">
            <p className="text-sm text-slate-500">{label}</p>
            <p className="mt-3 text-3xl font-semibold">{value}</p>
          </div>
        ))}
      </section>
      <section className="rounded-[2rem] border border-white/70 bg-white/70 p-6 shadow-soft backdrop-blur">
        <h2 className="font-semibold">最近文章</h2>
        <div className="mt-4 space-y-3">
          {dashboard.recentArticles.map((article) => (
            <a href={`/admin/articles/${article.id}/edit`} key={article.id} className="flex items-center justify-between gap-4 rounded-2xl bg-white/70 p-4 hover:bg-white">
              <div>
                <p className="font-medium">{article.title}</p>
                <p className="mt-1 text-sm text-slate-500">{article.updatedAt}</p>
              </div>
              <StatusBadge status={article.status} />
            </a>
          ))}
          {!loading && dashboard.recentArticles.length === 0 && <p className="rounded-2xl bg-white/70 p-4 text-sm text-slate-500">暂无文章。</p>}
        </div>
      </section>
    </div>
  );
}
