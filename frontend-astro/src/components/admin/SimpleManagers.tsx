import { useEffect, useState } from "react";
import { createCategory, getAdminCategories, updateCategorySort, updateCategoryStatus } from "@/api/category";
import { deleteAdminComment, getAdminComments, getAdminGuestbookMessages, restoreAdminComment } from "@/api/comment";
import { createFriendLink, deleteFriendLink, getAdminFriendLinks, updateFriendLink } from "@/api/friend";
import { createTag, getAdminTags, updateTagSort, updateTagStatus } from "@/api/tag";
import { listData, mapAdminCategory, mapAdminComment, mapAdminTag, mapFriendLink, type AdminCategory, type AdminComment, type AdminTag } from "@/lib/adminAdapters";
import type { FriendLink } from "@/lib/types";
import { MarkdownContent } from "../public/MarkdownContent";

type TaxonomyItem = AdminCategory | AdminTag;
type TaxonomyFormData = {
  name: string;
  slug: string;
  sort: string;
};

type FriendFormData = {
  name: string;
  url: string;
  logo: string;
  description: string;
  sort: string;
  status: "0" | "1";
};

const initialTaxonomyFormData: TaxonomyFormData = {
  name: "",
  slug: "",
  sort: "",
};

const initialFriendFormData: FriendFormData = {
  name: "",
  url: "",
  logo: "",
  description: "",
  sort: "0",
  status: "1",
};

function logError(context: string, err: unknown) {
  console.error(`[admin] ${context}`, err);
}

export function CategoryManager() {
  return (
    <TaxonomyManager
      title="分类管理"
      addLabel="新增分类"
      nameLabel="分类名"
      load={async () => listData(await getAdminCategories()).map(mapAdminCategory)}
      create={(payload) => createCategory({ category_name: payload.name, slug: payload.slug, sort: payload.sort })}
      updateStatus={updateCategoryStatus}
      updateSort={updateCategorySort}
    />
  );
}

export function TagManager() {
  return (
    <TaxonomyManager
      title="Tag 管理"
      addLabel="新增 Tag"
      nameLabel="Tag 名"
      load={async () => listData(await getAdminTags()).map(mapAdminTag)}
      create={(payload) => createTag({ tag_name: payload.name, slug: payload.slug, sort: payload.sort })}
      updateStatus={updateTagStatus}
      updateSort={updateTagSort}
    />
  );
}

export function CommentManager() {
  return <AdminCommentList title="评论管理" subtitle="查看和管理文章评论。" targetType={1} />;
}

export function ArticleCommentManager({ articleId }: { articleId: number }) {
  return <AdminCommentList title="文章评论" subtitle={`文章 #${articleId} 下的评论。`} targetType={1} targetId={articleId} backHref="/admin/comments" />;
}

export function GuestbookManager() {
  return <AdminCommentList title="留言管理" subtitle="查看留言板消息，删除或恢复显示状态。" targetType={2} />;
}

function AdminCommentList({
  title,
  subtitle,
  targetType,
  targetId,
  backHref,
}: {
  title: string;
  subtitle: string;
  targetType: 1 | 2;
  targetId?: number;
  backHref?: string;
}) {
  const [items, setItems] = useState<AdminComment[]>([]);
  const [loading, setLoading] = useState(true);
  const [savingId, setSavingId] = useState<number | null>(null);
  const [error, setError] = useState("");

  async function loadItems() {
    setLoading(true);
    setError("");
    try {
      const params = {
        target_id: targetId,
        include_deleted: true,
        page: 1,
        page_size: 100,
      };
      const response = targetType === 2
        ? await getAdminGuestbookMessages(params)
        : await getAdminComments({ ...params, target_type: 1 });
      setItems(listData(response).map(mapAdminComment));
    } catch (err) {
      logError(`${title} load failed`, err);
      setError(err instanceof Error ? err.message : "加载失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadItems();
  }, [targetId, targetType]);

  async function removeItem(item: AdminComment) {
    if (!window.confirm(`确认删除「${item.nickname}」的${targetType === 1 ? "评论" : "留言"}？`)) return;
    setSavingId(item.id);
    setError("");
    try {
      await deleteAdminComment(item.id);
      await loadItems();
    } catch (err) {
      logError(`${title} delete failed`, err);
      setError(err instanceof Error ? err.message : "删除失败，请稍后重试");
    } finally {
      setSavingId(null);
    }
  }

  async function restoreItem(item: AdminComment) {
    setSavingId(item.id);
    setError("");
    try {
      await restoreAdminComment(item.id);
      await loadItems();
    } catch (err) {
      logError(`${title} restore failed`, err);
      setError(err instanceof Error ? err.message : "恢复失败，请稍后重试");
    } finally {
      setSavingId(null);
    }
  }

  return (
    <div className="rounded-[2rem] border border-white/70 bg-white/70 p-6 shadow-soft backdrop-blur">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-semibold">{title}</h2>
          <p className="mt-1 text-sm text-slate-500">{subtitle}</p>
        </div>
        {backHref && (
          <a href={backHref} className="rounded-full bg-white/75 px-4 py-2 text-sm text-slate-600 hover:bg-white">
            返回评论管理
          </a>
        )}
      </div>
      {error && <p className="mt-4 rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">{error}</p>}
      <div className="mt-5 space-y-3">
        {loading && <p className="rounded-2xl bg-white/70 p-4 text-sm text-slate-500">加载中...</p>}
        {!loading && items.length === 0 && <p className="rounded-2xl bg-white/70 p-4 text-sm text-slate-500">暂无数据。</p>}
        {!loading && items.map((item) => (
          <AdminMessageCard
            key={item.id}
            item={item}
            saving={savingId === item.id}
            onDelete={() => void removeItem(item)}
            onRestore={() => void restoreItem(item)}
          />
        ))}
      </div>
    </div>
  );
}

export function FriendLinkManager() {
  const [items, setItems] = useState<FriendLink[]>([]);
  const [editing, setEditing] = useState<FriendLink | "new" | null>(null);
  const [formData, setFormData] = useState<FriendFormData>(initialFriendFormData);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  async function loadItems() {
    setLoading(true);
    setError("");
    try {
      setItems(listData(await getAdminFriendLinks({ page: 1, pageSize: 100 })).map(mapFriendLink));
    } catch (err) {
      logError("load friend links failed", err);
      setError("操作失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadItems();
  }, []);

  async function removeItem(item: FriendLink) {
    if (!window.confirm(`确认删除「${item.name}」？`)) return;
    setSaving(true);
    setError("");
    try {
      await deleteFriendLink(item.id);
      await loadItems();
    } catch (err) {
      logError("delete friend link failed", err);
      setError("操作失败，请稍后重试");
    } finally {
      setSaving(false);
    }
  }

  async function submitFriend(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const payload = {
      name: formData.name,
      url: formData.url,
      logo: formData.logo,
      description: formData.description,
      sort: Number(formData.sort || 0),
      status: Number(formData.status) as 0 | 1,
    };
    setSaving(true);
    setError("");
    try {
      if (editing && editing !== "new") {
        await updateFriendLink(editing.id, { ...payload, id: editing.id });
      } else {
        await createFriendLink(payload);
      }
      setEditing(null);
      setFormData(initialFriendFormData);
      await loadItems();
    } catch (err) {
      logError("save friend link failed", err);
      setError("操作失败，请稍后重试");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="rounded-[2rem] border border-white/70 bg-white/70 p-6 shadow-soft backdrop-blur">
      <div className="flex items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-semibold">友链管理</h2>
          <p className="mt-1 text-sm text-slate-500">新增、编辑和删除朋友链接。</p>
        </div>
        <button onClick={() => {
          setEditing("new");
          setFormData(initialFriendFormData);
        }} className="rounded-full bg-slate-950 px-4 py-2 text-sm text-white">新增友链</button>
      </div>
      {error && <p className="mt-4 rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">{error}</p>}
      {editing && <FriendForm formData={formData} setFormData={setFormData} saving={saving} onSubmit={submitFriend} onClose={() => setEditing(null)} />}
      <div className="mt-5 grid gap-3">
        {loading && <p className="rounded-2xl bg-white/70 p-4 text-sm text-slate-500">加载中...</p>}
        {!loading && items.map((item) => (
          <div key={item.id} className="flex flex-wrap items-center justify-between gap-4 rounded-2xl bg-white/75 p-4">
            <div className="min-w-0">
              <p className="break-words font-medium">{item.name}</p>
              {item.description && <p className="mt-1 break-words text-sm text-slate-500">{item.description}</p>}
              <p className="mt-1 break-all text-sm text-slate-500">{item.url}</p>
              <p className="mt-1 text-xs text-slate-400">排序 {item.sort} · {item.status === 1 ? "显示" : "隐藏"}</p>
            </div>
            <div className="flex gap-2">
              <button onClick={() => {
                setEditing(item);
                setFormData({
                  name: item.name,
                  url: item.url,
                  logo: item.logo || "",
                  description: item.description || "",
                  sort: String(item.sort),
                  status: String(item.status) as "0" | "1",
                });
              }} className="rounded-full bg-slate-100 px-3 py-1.5 text-xs">编辑</button>
              <button disabled={saving} onClick={() => void removeItem(item)} className="rounded-full bg-rose-50 px-3 py-1.5 text-xs text-rose-700 disabled:opacity-50">删除</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function TaxonomyManager({
  title,
  addLabel,
  nameLabel,
  load,
  create,
  updateStatus,
  updateSort,
}: {
  title: string;
  addLabel: string;
  nameLabel: string;
  load: () => Promise<TaxonomyItem[]>;
  create: (payload: { name: string; slug: string; sort: number }) => Promise<unknown>;
  updateStatus: (id: number | string, status: 0 | 1) => Promise<unknown>;
  updateSort: (id: number | string, sort: number) => Promise<unknown>;
}) {
  const [items, setItems] = useState<TaxonomyItem[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState<TaxonomyFormData>(initialTaxonomyFormData);
  const [loading, setLoading] = useState(true);
  const [savingId, setSavingId] = useState<number | "new" | null>(null);
  const [error, setError] = useState("");

  async function loadItems() {
    setLoading(true);
    setError("");
    try {
      setItems(await load());
    } catch (err) {
      logError(`${title} load failed`, err);
      setError("操作失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadItems();
  }, []);

  async function submitCreate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSavingId("new");
    setError("");
    try {
      await create({
        name: formData.name,
        slug: formData.slug,
        sort: Number(formData.sort || 0),
      });
      setShowForm(false);
      setFormData(initialTaxonomyFormData);
      await loadItems();
    } catch (err) {
      logError(`${title} create failed`, err);
      setError("操作失败，请稍后重试");
    } finally {
      setSavingId(null);
    }
  }

  async function toggleStatus(item: TaxonomyItem) {
    setSavingId(item.id);
    setError("");
    try {
      await updateStatus(item.id, item.status === 1 ? 0 : 1);
      await loadItems();
    } catch (err) {
      logError(`${title} update status failed`, err);
      setError("操作失败，请稍后重试");
    } finally {
      setSavingId(null);
    }
  }

  async function submitSort(event: React.FormEvent<HTMLFormElement>, item: TaxonomyItem) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    setSavingId(item.id);
    setError("");
    try {
      await updateSort(item.id, Number(form.get("sort") || 0));
      await loadItems();
    } catch (err) {
      logError(`${title} update sort failed`, err);
      setError("操作失败，请稍后重试");
    } finally {
      setSavingId(null);
    }
  }

  return (
    <div className="rounded-[2rem] border border-white/70 bg-white/70 p-6 shadow-soft backdrop-blur">
      <div className="flex items-center justify-between gap-4">
        <h2 className="text-xl font-semibold">{title}</h2>
        <button onClick={() => {
          setShowForm((current) => !current);
          setFormData(initialTaxonomyFormData);
        }} className="rounded-full bg-slate-950 px-4 py-2 text-sm text-white">{addLabel}</button>
      </div>
      {error && <p className="mt-4 rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">{error}</p>}
      {showForm && (
        <form onSubmit={submitCreate} className="mt-5 grid gap-3 rounded-2xl bg-white/70 p-4 sm:grid-cols-[1fr_1fr_120px_auto]">
          <input value={formData.name} onChange={(event) => setFormData((current) => ({ ...current, name: event.target.value }))} placeholder={nameLabel} className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none" />
          <input value={formData.slug} onChange={(event) => setFormData((current) => ({ ...current, slug: event.target.value }))} placeholder="slug" className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none" />
          <input value={formData.sort} onChange={(event) => setFormData((current) => ({ ...current, sort: event.target.value }))} type="number" placeholder="排序" className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none" />
          <button disabled={savingId === "new"} className="rounded-full bg-slate-950 px-4 py-2 text-sm text-white disabled:opacity-50">保存</button>
        </form>
      )}
      <div className="mt-5 space-y-3">
        {loading && <p className="rounded-2xl bg-white/70 p-4 text-sm text-slate-500">加载中...</p>}
        {!loading && items.map((item) => (
          <div key={item.id} className="grid gap-4 rounded-2xl bg-white/75 p-4 lg:grid-cols-[minmax(0,1fr)_auto_auto] lg:items-center">
            <div className="min-w-0">
              <p className="break-words font-medium">{item.name}</p>
              <p className="mt-1 text-sm text-slate-500">{item.slug || "-"}</p>
              <p className="mt-1 text-xs text-slate-400">文章数 {item.count} · {item.status === 1 ? "启用" : "停用"}</p>
            </div>
            <form onSubmit={(event) => void submitSort(event, item)} className="flex items-center gap-2">
              <input name="sort" type="number" defaultValue={item.sort} className="w-24 rounded-2xl border border-slate-100 px-3 py-2 text-sm outline-none" />
              <button disabled={savingId === item.id} className="rounded-full bg-slate-100 px-3 py-1.5 text-xs disabled:opacity-50">保存排序</button>
            </form>
            <button disabled={savingId === item.id} onClick={() => void toggleStatus(item)} className={`rounded-full px-3 py-1.5 text-xs disabled:opacity-50 ${item.status === 1 ? "bg-amber-50 text-amber-700" : "bg-emerald-50 text-emerald-700"}`}>
              {item.status === 1 ? "停用" : "启用"}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

function AdminMessageCard({
  item,
  saving,
  onDelete,
  onRestore,
}: {
  item: AdminComment;
  saving: boolean;
  onDelete: () => void;
  onRestore: () => void;
}) {
  return (
    <article className="rounded-2xl bg-white/75 p-4">
      <div className="flex flex-wrap justify-between gap-3">
        <div>
          <div className="flex flex-wrap items-center gap-2">
            <p className="font-medium">{item.nickname}</p>
            {item.articleTitle && <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs text-slate-600">{item.articleTitle}</span>}
            {item.deleted && <span className="rounded-full bg-rose-50 px-2 py-0.5 text-xs text-rose-700">已删除</span>}
          </div>
          <p className="mt-1 text-xs text-slate-400">
            {item.email}
            {item.website && <span> · {item.website}</span>}
            <span> · {item.createdAt}</span>
          </p>
          {(item.ip || item.location || item.userAgent) && (
            <p className="mt-1 text-xs text-slate-400">
              {[item.ip, item.location, item.userAgent].filter(Boolean).join(" · ")}
            </p>
          )}
        </div>
        <div className="flex gap-2">
          {!item.deleted && <button disabled={saving} onClick={onDelete} className="rounded-full bg-rose-50 px-3 py-1.5 text-xs text-rose-700 disabled:opacity-50">删除</button>}
          {item.deleted && <button disabled={saving} onClick={onRestore} className="rounded-full bg-slate-100 px-3 py-1.5 text-xs disabled:opacity-50">恢复</button>}
        </div>
      </div>
      <MarkdownContent content={item.content} className="text-sm leading-6 text-slate-600" />
    </article>
  );
}

function FriendForm({
  formData,
  setFormData,
  saving,
  onSubmit,
  onClose,
}: {
  formData: FriendFormData;
  setFormData: React.Dispatch<React.SetStateAction<FriendFormData>>;
  saving: boolean;
  onSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
  onClose: () => void;
}) {
  return (
    <form onSubmit={onSubmit} className="mt-5 grid gap-3 rounded-2xl bg-white/70 p-4 sm:grid-cols-2">
      <input value={formData.name} onChange={(event) => setFormData((current) => ({ ...current, name: event.target.value }))} placeholder="站点名称" className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none" />
      <input value={formData.url} onChange={(event) => setFormData((current) => ({ ...current, url: event.target.value }))} placeholder="https://example.com" className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none" />
      <input value={formData.logo} onChange={(event) => setFormData((current) => ({ ...current, logo: event.target.value }))} placeholder="Logo URL" className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none" />
      <input value={formData.description} onChange={(event) => setFormData((current) => ({ ...current, description: event.target.value }))} placeholder="这个站点的一句话介绍" maxLength={255} className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none sm:col-span-2" />
      <input value={formData.sort} onChange={(event) => setFormData((current) => ({ ...current, sort: event.target.value }))} type="number" placeholder="排序" className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none" />
      <select value={formData.status} onChange={(event) => setFormData((current) => ({ ...current, status: event.target.value as "0" | "1" }))} className="rounded-2xl border border-slate-100 px-4 py-3 text-sm outline-none">
        <option value={1}>显示</option>
        <option value={0}>隐藏</option>
      </select>
      <div className="flex gap-3">
        <button disabled={saving} className="flex-1 rounded-full bg-slate-950 px-4 py-2 text-sm text-white disabled:opacity-50">保存</button>
        <button type="button" onClick={onClose} className="flex-1 rounded-full bg-white px-4 py-2 text-sm text-slate-700">取消</button>
      </div>
    </form>
  );
}
