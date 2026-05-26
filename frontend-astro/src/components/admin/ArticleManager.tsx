import { useEffect, useMemo, useState } from "react";
import { deleteArticle, getAdminArticles, updateArticleStatus } from "@/api/article";
import { getAdminCategories } from "@/api/category";
import { getAdminTags } from "@/api/tag";
import { listData, mapAdminArticle, mapAdminCategory, mapAdminTag, type AdminArticle, type AdminCategory, type AdminTag } from "@/lib/adminAdapters";
import type { ArticleStatus } from "@/lib/types";
import { StatusBadge } from "./StatusBadge";

const tabs: Array<{ key: "all" | ArticleStatus; label: string }> = [
  { key: "all", label: "全部" },
  { key: "draft", label: "草稿箱" },
  { key: "published", label: "已发布" },
  { key: "offline", label: "已下架" },
  { key: "deleted", label: "回收站" },
];

const actionButtonClass = "inline-flex h-8 items-center justify-center rounded-full px-3 text-xs leading-none transition disabled:cursor-not-allowed disabled:opacity-50";

function statusQuery(status: "all" | ArticleStatus) {
  if (status === "draft") return 0;
  if (status === "published") return 1;
  if (status === "offline") return 2;
  return -1;
}

function deletedQuery(status: "all" | ArticleStatus) {
  if (status === "deleted") return 1;
  return status === "all" ? -1 : 0;
}

export function ArticleManager() {
  const [active, setActive] = useState<"all" | ArticleStatus>("all");
  const [keyword, setKeyword] = useState("");
  const [categoryId, setCategoryId] = useState("all");
  const [tagId, setTagId] = useState("all");
  const [articles, setArticles] = useState<AdminArticle[]>([]);
  const [categories, setCategories] = useState<AdminCategory[]>([]);
  const [tags, setTags] = useState<AdminTag[]>([]);
  const [loading, setLoading] = useState(true);
  const [savingId, setSavingId] = useState<number | null>(null);
  const [error, setError] = useState("");

  async function loadData() {
    setLoading(true);
    setError("");
    try {
      const [articleResult, categoryResult, tagResult] = await Promise.all([
        getAdminArticles({
          page: 1,
          pageSize: 100,
          status: statusQuery(active),
          is_deleted: deletedQuery(active),
          category_id: categoryId === "all" ? 0 : Number(categoryId),
          keyword: keyword.trim() || undefined,
        }),
        getAdminCategories(),
        getAdminTags(),
      ]);

      setArticles(listData(articleResult).map(mapAdminArticle));
      setCategories(listData(categoryResult).map(mapAdminCategory));
      setTags(listData(tagResult).map(mapAdminTag));
    } catch (err) {
      setError(err instanceof Error ? err.message : "加载文章失败");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadData();
  }, [active, categoryId]);

  const categoryNames = useMemo(() => new Map(categories.map((item) => [item.id, item.name])), [categories]);

  const filtered = useMemo(() => {
    const value = keyword.trim().toLowerCase();
    return articles.filter((article) => {
      if (value && !article.title.toLowerCase().includes(value) && !article.summary.toLowerCase().includes(value)) return false;
      if (tagId !== "all" && article.tagIds.length > 0 && !article.tagIds.includes(Number(tagId))) return false;
      return true;
    });
  }, [articles, keyword, tagId]);

  async function changeStatus(article: AdminArticle, nextStatus: 1 | 2) {
    setSavingId(article.id);
    setError("");
    try {
      await updateArticleStatus(article.id, nextStatus);
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "更新文章状态失败");
    } finally {
      setSavingId(null);
    }
  }

  async function removeArticle(article: AdminArticle) {
    if (!window.confirm(`确认删除「${article.title}」？`)) return;
    setSavingId(article.id);
    setError("");
    try {
      await deleteArticle(article.id);
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "删除文章失败");
    } finally {
      setSavingId(null);
    }
  }

  return (
    <div className="rounded-[2rem] border border-white/70 bg-white/70 p-5 shadow-soft backdrop-blur">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-semibold">文章管理</h2>
          <p className="mt-1 text-sm text-slate-500">整理草稿、发布记录和回收站。</p>
        </div>
        <a href="/admin/articles/new" className="rounded-full bg-slate-950 px-5 py-2.5 text-sm text-white">写文章</a>
      </div>
      <div className="mt-6 flex flex-wrap gap-2">
        {tabs.map((tab) => (
          <button key={tab.key} onClick={() => setActive(tab.key)} className={`rounded-full px-4 py-2 text-sm ${active === tab.key ? "bg-slate-950 text-white" : "bg-white/70 text-slate-600 hover:bg-white"}`}>
            {tab.label}
          </button>
        ))}
      </div>
      <div className="mt-5 grid gap-3 md:grid-cols-[1fr_180px_180px_auto]">
        <input value={keyword} onChange={(event) => setKeyword(event.target.value)} placeholder="按标题搜索" className="rounded-2xl border border-white/70 bg-white/80 px-4 py-3 text-sm outline-none" />
        <select value={categoryId} onChange={(event) => setCategoryId(event.target.value)} className="rounded-2xl border border-white/70 bg-white/80 px-4 py-3 text-sm outline-none">
          <option value="all">全部分类</option>
          {categories.map((item) => <option value={item.id} key={item.id}>{item.name}</option>)}
        </select>
        <select value={tagId} onChange={(event) => setTagId(event.target.value)} className="rounded-2xl border border-white/70 bg-white/80 px-4 py-3 text-sm outline-none">
          <option value="all">全部 Tags</option>
          {tags.map((item) => <option value={item.id} key={item.id}>{item.name}</option>)}
        </select>
        <button onClick={() => void loadData()} className="rounded-full bg-white px-4 py-2 text-sm text-slate-700 shadow-sm">刷新</button>
      </div>
      {error && <p className="mt-4 rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">{error}</p>}
      <div className="mt-5 space-y-3">
        {loading && <p className="rounded-2xl bg-white/70 p-5 text-sm text-slate-500">加载中...</p>}
        {!loading && filtered.length === 0 && <p className="rounded-2xl bg-white/70 p-5 text-sm text-slate-500">暂无文章。</p>}
        {!loading && filtered.map((article) => (
          <article key={article.id} className="rounded-[1.5rem] border border-slate-100 bg-white/75 p-5">
            <div className="grid gap-5 lg:grid-cols-[minmax(0,1fr)_auto]">
              <div className="min-w-0">
                <div className="flex flex-wrap items-center gap-2">
                  <h3 className="min-w-0 break-words font-semibold text-slate-950">{article.title}</h3>
                  {article.isTop && <span className="rounded-full bg-rose-50 px-2.5 py-1 text-xs text-rose-700">置顶</span>}
                  <StatusBadge status={article.status} />
                </div>
                <p className="mt-2 line-clamp-2 text-sm leading-6 text-slate-500">{article.summary}</p>
                <div className="mt-3 flex flex-wrap items-center gap-2 text-xs text-slate-500">
                  <span className="rounded-full bg-slate-100 px-3 py-1.5 text-slate-600">{categoryNames.get(article.categoryId) || `分类 #${article.categoryId}`}</span>
                  {article.tagIds.map((id) => (
                    <span key={id} className="rounded-full bg-white px-2.5 py-1 text-slate-500">
                      #{tags.find((tag) => tag.id === id)?.name || id}
                    </span>
                  ))}
                  <span>评论：{article.commentCount}</span>
                  <span>置顶：{article.isTop ? "是" : "否"}</span>
                  <span>发布时间：{article.publishedAt || "-"}</span>
                  <span>更新时间：{article.updatedAt || "-"}</span>
                </div>
              </div>

              <div className="flex shrink-0 flex-wrap items-center gap-2 lg:justify-end">
                <a href={`/admin/articles/new?edit=${article.id}`} className={`${actionButtonClass} bg-slate-100 text-slate-700 hover:bg-slate-200/80`}>编辑</a>
                {article.status !== "published" && article.status !== "deleted" && <button disabled={savingId === article.id} onClick={() => void changeStatus(article, 1)} className={`${actionButtonClass} bg-emerald-50 text-emerald-700 hover:bg-emerald-100/80`}>发布</button>}
                {article.status === "published" && <button disabled={savingId === article.id} onClick={() => void changeStatus(article, 2)} className={`${actionButtonClass} bg-amber-50 text-amber-700 hover:bg-amber-100/80`}>下架</button>}
                {article.status !== "deleted" && <button disabled={savingId === article.id} onClick={() => void removeArticle(article)} className={`${actionButtonClass} bg-rose-50 text-rose-700 hover:bg-rose-100/80`}>删除</button>}
              </div>
            </div>
          </article>
        ))}
      </div>
    </div>
  );
}
