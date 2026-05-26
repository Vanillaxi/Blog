import { useEffect, useMemo, useState } from "react";
import { createArticle, getAdminArticleDetail, updateArticle, updateArticleStatus, type ArticlePayload, type ArticleSaveStatus } from "@/api/article";
import { getAdminCategories } from "@/api/category";
import { getAdminTags } from "@/api/tag";
import { listData, mapAdminCategory, mapAdminTag, type AdminCategory, type AdminTag } from "@/lib/adminAdapters";

type FormState = {
  title: string;
  summary: string;
  content: string;
  coverUrl: string;
  categoryId: string;
  tagIds: number[];
  isTop: boolean;
  status: ArticleSaveStatus;
};

const initialForm: FormState = {
  title: "",
  summary: "",
  content: "",
  coverUrl: "",
  categoryId: "",
  tagIds: [],
  isTop: false,
  status: 0,
};

function asRecord(value: unknown): Record<string, unknown> {
  return value && typeof value === "object" ? (value as Record<string, unknown>) : {};
}

function numberValue(value: unknown, fallback = 0) {
  const next = Number(value);
  return Number.isFinite(next) ? next : fallback;
}

function stringValue(value: unknown) {
  return typeof value === "string" ? value : value == null ? "" : String(value);
}

function getEditId(articleId?: number) {
  if (articleId) return articleId;
  if (typeof window === "undefined") return undefined;

  const queryID = Number(new URLSearchParams(window.location.search).get("edit"));
  if (Number.isFinite(queryID) && queryID > 0) return queryID;

  const match = window.location.pathname.match(/\/admin\/articles\/(\d+)\/edit/);
  return match ? Number(match[1]) : undefined;
}

function mapDetailToForm(value: unknown): FormState {
  const item = asRecord(value);
  const rawTagIDs = Array.isArray(item.tag_ids) ? item.tag_ids : Array.isArray(item.tags) ? item.tags : [];

  return {
    title: stringValue(item.title),
    summary: stringValue(item.summary),
    content: stringValue(item.content),
    coverUrl: stringValue(item.cover_url),
    categoryId: String(numberValue(item.category_id) || ""),
    tagIds: rawTagIDs
      .map((tag) => {
        const record = asRecord(tag);
        return numberValue(record.id ?? record.tag_id ?? tag, NaN);
      })
      .filter((id) => Number.isFinite(id)),
    isTop: numberValue(item.is_top) === 1,
    status: numberValue(item.status, 0) as ArticleSaveStatus,
  };
}

function statusLabel(status: ArticleSaveStatus) {
  if (status === 1) return "发布";
  if (status === 2) return "下架";
  return "草稿";
}

export function ArticleEditor({ articleId }: { articleId?: number }) {
  const [editId, setEditId] = useState<number | undefined>(articleId);
  const [form, setForm] = useState<FormState>(initialForm);
  const [originalStatus, setOriginalStatus] = useState<ArticleSaveStatus>(0);
  const [categories, setCategories] = useState<AdminCategory[]>([]);
  const [tags, setTags] = useState<AdminTag[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const isEdit = Boolean(editId);
  const activeCategories = useMemo(() => categories.filter((item) => item.status === 1), [categories]);
  const activeTags = useMemo(() => tags.filter((item) => item.status === 1), [tags]);

  useEffect(() => {
    setEditId(getEditId(articleId));
  }, [articleId]);

  useEffect(() => {
    async function loadEditorData() {
      setLoading(true);
      setError("");
      try {
        const [categoryResult, tagResult] = await Promise.all([getAdminCategories(), getAdminTags()]);
        const nextCategories = listData(categoryResult).map(mapAdminCategory);
        const nextTags = listData(tagResult).map(mapAdminTag);
        setCategories(nextCategories);
        setTags(nextTags);

        if (editId) {
          const detail = await getAdminArticleDetail(editId);
          const nextForm = mapDetailToForm(asRecord(detail).data);
          setForm(nextForm);
          setOriginalStatus(nextForm.status);
        } else {
          setForm((current) => ({
            ...current,
            categoryId: current.categoryId || String(nextCategories.find((item) => item.status === 1)?.id ?? ""),
          }));
          setOriginalStatus(0);
        }
      } catch (err) {
        console.error("[admin] article editor load failed", err);
        setError(err instanceof Error ? err.message : "加载写作页面失败");
      } finally {
        setLoading(false);
      }
    }

    void loadEditorData();
  }, [editId]);

  function updateForm(patch: Partial<FormState>) {
    setForm((current) => ({ ...current, ...patch }));
  }

  function toggleTag(tagId: number) {
    setForm((current) => ({
      ...current,
      tagIds: current.tagIds.includes(tagId) ? current.tagIds.filter((id) => id !== tagId) : [...current.tagIds, tagId],
    }));
  }

  function validate() {
    if (!form.title.trim()) return "请填写文章标题";
    if (!form.content.trim()) return "请填写文章内容";
    if (!Number(form.categoryId)) return "请选择分类";
    return "";
  }

  async function submitArticle(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validationError = validate();
    setError("");
    setSuccess("");
    if (validationError) {
      setError(validationError);
      return;
    }

    const payload: ArticlePayload = {
      title: form.title.trim(),
      summary: form.summary.trim(),
      content: form.content,
      cover_url: form.coverUrl.trim(),
      category_id: Number(form.categoryId),
      tag_ids: form.tagIds,
      is_top: form.isTop ? 1 : 0,
    };
    console.debug("[admin] article save payload", payload, { status: form.status });

    setSaving(true);
    try {
      const response = editId ? await updateArticle(editId, payload) : await createArticle(payload);
      const savedID = editId ?? numberValue(asRecord(asRecord(response).data).id);

      if (savedID && form.status !== originalStatus) {
        await updateArticleStatus(savedID, form.status);
      }

      setSuccess("文章保存成功");
      setOriginalStatus(form.status);
      if (!editId && savedID) {
        setEditId(savedID);
        window.history.replaceState(null, "", `/admin/articles/new?edit=${savedID}`);
      }
    } catch (err) {
      console.error("[admin] article save failed", err);
      setError(err instanceof Error ? err.message : "文章保存失败");
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return <p className="rounded-[2rem] border border-white/70 bg-white/70 p-6 text-sm text-slate-500 shadow-soft backdrop-blur">写作页面加载中...</p>;
  }

  return (
    <form onSubmit={submitArticle} className="grid gap-5 lg:grid-cols-[1fr_320px]">
      <section className="rounded-[2rem] border border-white/70 bg-white/70 p-6 shadow-soft backdrop-blur">
        <input value={form.title} onChange={(event) => updateForm({ title: event.target.value })} placeholder="文章标题" className="w-full border-0 bg-transparent text-3xl font-semibold outline-none placeholder:text-slate-300" />
        <textarea value={form.summary} onChange={(event) => updateForm({ summary: event.target.value })} placeholder="摘要" className="mt-5 min-h-24 w-full rounded-2xl border border-slate-100 bg-white/70 px-4 py-3 text-sm outline-none" />
        <textarea value={form.content} onChange={(event) => updateForm({ content: event.target.value })} placeholder="Markdown 内容" className="mt-5 min-h-[420px] w-full rounded-[1.5rem] border border-slate-100 bg-white/70 px-4 py-4 font-mono text-sm leading-7 outline-none" />
      </section>
      <aside className="space-y-4">
        <div className="rounded-[2rem] border border-white/70 bg-white/70 p-5 shadow-soft backdrop-blur">
          <h2 className="font-semibold">发布设置</h2>
          <label className="mt-4 block text-sm text-slate-500">封面图 URL</label>
          <input value={form.coverUrl} onChange={(event) => updateForm({ coverUrl: event.target.value })} className="mt-2 w-full rounded-2xl border border-slate-100 bg-white/80 px-4 py-3 text-sm outline-none" />
          <label className="mt-4 block text-sm text-slate-500">分类</label>
          <select value={form.categoryId} onChange={(event) => updateForm({ categoryId: event.target.value })} className="mt-2 w-full rounded-2xl border border-slate-100 bg-white/80 px-4 py-3 text-sm outline-none">
            <option value="">请选择分类</option>
            {activeCategories.map((item) => <option value={item.id} key={item.id}>{item.name}</option>)}
          </select>
          {activeCategories.length === 0 && <p className="mt-2 rounded-2xl bg-amber-50 px-3 py-2 text-xs text-amber-700">请先创建分类。</p>}
          <label className="mt-4 block text-sm text-slate-500">Tags</label>
          <div className="mt-2 flex flex-wrap gap-2">
            {activeTags.map((tag) => (
              <label key={tag.id} className="rounded-full bg-white px-3 py-1.5 text-xs text-slate-600">
                <input name="tag_ids" value={tag.id} type="checkbox" checked={form.tagIds.includes(tag.id)} onChange={() => toggleTag(tag.id)} className="mr-1" />
                {tag.name}
              </label>
            ))}
          </div>
          {activeTags.length === 0 && <p className="mt-2 rounded-2xl bg-amber-50 px-3 py-2 text-xs text-amber-700">请先创建标签。</p>}
          <label className="mt-4 flex items-center justify-between gap-4 rounded-2xl border border-slate-100 bg-white/70 px-4 py-3">
            <span>
              <span className="block text-sm text-slate-600">置顶文章</span>
              <span className="mt-0.5 block text-xs text-slate-400">保存时对应后端字段 is_top</span>
            </span>
            <input name="is_top" type="checkbox" checked={form.isTop} onChange={(event) => updateForm({ isTop: event.target.checked })} className="h-5 w-5 accent-rose-400" />
          </label>
          <label className="mt-4 block text-sm text-slate-500">状态</label>
          <select value={form.status} onChange={(event) => updateForm({ status: Number(event.target.value) as ArticleSaveStatus })} className="mt-2 w-full rounded-2xl border border-slate-100 bg-white/80 px-4 py-3 text-sm outline-none">
            <option value={0}>草稿</option>
            <option value={1}>发布</option>
            <option value={2}>下架</option>
          </select>
          <p className="mt-2 text-xs text-slate-400">当前将保存为：{statusLabel(form.status)}</p>
        </div>
        {error && <p className="rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">{error}</p>}
        {success && <p className="rounded-2xl bg-emerald-50 px-4 py-3 text-sm text-emerald-700">{success}</p>}
        <div className="flex gap-3">
          <button type="button" onClick={() => updateForm({ status: 0 })} className="flex-1 rounded-full bg-white px-4 py-3 text-sm text-slate-700 shadow-sm">设为草稿</button>
          <button disabled={saving || activeCategories.length === 0} className="flex-1 rounded-full bg-slate-950 px-4 py-3 text-sm text-white disabled:cursor-not-allowed disabled:opacity-50">
            {saving ? "保存中..." : isEdit ? "保存文章" : "保存文章"}
          </button>
        </div>
      </aside>
    </form>
  );
}
