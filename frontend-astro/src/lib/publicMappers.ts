import type { Article, Category, Comment, FriendLink, GuestbookMessage, Tag } from "@/lib/types";

type AnyRecord = Record<string, unknown>;

export type PaginatedData<T = unknown> = {
  list?: T[];
  links?: T[];
  total?: number;
  page?: number;
  size?: number;
  page_size?: number;
};

function asRecord(value: unknown): AnyRecord {
  return value && typeof value === "object" ? (value as AnyRecord) : {};
}

function asArray(value: unknown): unknown[] {
  return Array.isArray(value) ? value : [];
}

function text(value: unknown, fallback = "") {
  if (typeof value === "string") return value;
  if (typeof value === "number") return String(value);
  return fallback;
}

function number(value: unknown, fallback = 0) {
  if (typeof value === "number" && Number.isFinite(value)) return value;
  if (typeof value === "string") {
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : fallback;
  }
  return fallback;
}

function boolean(value: unknown) {
  if (typeof value === "boolean") return value;
  if (typeof value === "number") return value === 1;
  if (typeof value === "string") return value === "1" || value.toLowerCase() === "true";
  return false;
}

function padDatePart(value: number) {
  return String(value).padStart(2, "0");
}

function parseLocalDateTime(value: unknown) {
  if (!value) return "";

  const raw = String(value).trim();
  if (!raw) return "";

  const hasTimezone = /(?:z|[+-]\d{2}:?\d{2})$/i.test(raw);
  const localDateTimeMatch = raw.match(/^(\d{4})[-/](\d{1,2})[-/](\d{1,2})(?:[ T](\d{1,2}):(\d{1,2})(?::(\d{1,2})(?:\.\d+)?)?)?/);

  if (localDateTimeMatch && !hasTimezone) {
    const year = Number(localDateTimeMatch[1]);
    const month = Number(localDateTimeMatch[2]);
    const day = Number(localDateTimeMatch[3]);
    const hour = Number(localDateTimeMatch[4] ?? 0);
    const minute = Number(localDateTimeMatch[5] ?? 0);
    const second = Number(localDateTimeMatch[6] ?? 0);
    const date = new Date(year, month - 1, day, hour, minute, second);

    if (
      date.getFullYear() !== year ||
      date.getMonth() !== month - 1 ||
      date.getDate() !== day ||
      date.getHours() !== hour ||
      date.getMinutes() !== minute
    ) {
      return "";
    }

    return date;
  }

  const date = new Date(raw);
  return Number.isNaN(date.getTime()) ? "" : date;
}

function formatDate(value: unknown) {
  const date = parseLocalDateTime(value);
  if (!date) return "";
  return `${date.getFullYear()}/${padDatePart(date.getMonth() + 1)}/${padDatePart(date.getDate())}`;
}

export function formatDateTime(value: unknown) {
  const date = parseLocalDateTime(value);
  if (!date) return "";

  return `${date.getFullYear()}/${padDatePart(date.getMonth() + 1)}/${padDatePart(date.getDate())} ${padDatePart(date.getHours())}:${padDatePart(date.getMinutes())}`;
}

function commentTime(item: AnyRecord) {
  return item.create_time ?? item.createTime ?? item.created_at ?? item.createdAt ?? item.update_time ?? item.updateTime;
}

export function getListPayload<T = unknown>(value: unknown): T[] {
  if (Array.isArray(value)) return value as T[];
  const data = asRecord(value);
  return asArray(data.list ?? data.articles ?? data.rows ?? data.items ?? data.records ?? data.comments ?? data.replies ?? data.links) as T[];
}

export function mapCategory(value: unknown): Category {
  const item = asRecord(value);
  return {
    id: number(item.id),
    name: text(item.name ?? item.category_name, "未分类"),
    slug: text(item.slug),
    count: number(item.count ?? item.article_count),
    sort: number(item.sort),
  };
}

export function mapTag(value: unknown): Tag {
  const item = asRecord(value);
  return {
    id: number(item.id),
    name: text(item.name ?? item.tag_name, "未命名"),
    slug: text(item.slug),
    sort: number(item.sort),
  };
}

export function mapArticle(value: unknown, categories: Category[] = [], tags: Tag[] = []): Article {
  const item = asRecord(value);
  const categoryId = number(item.category_id ?? item.categoryId ?? asRecord(item.category).id);
  const category =
    categories.find((candidate) => candidate.id === categoryId) ??
    (item.category
      ? mapCategory(item.category)
      : {
          id: categoryId,
          name: text(item.category_name ?? item.categoryName, categoryId ? `分类 ${categoryId}` : "未分类"),
          slug: "",
          count: 0,
          sort: 0,
        });
  const rawTags = asArray(item.tags);
  const rawTagIds = asArray(item.tag_ids ?? item.tagIds);

  return {
    id: number(item.id),
    title: text(item.title, "未命名文章"),
    summary: text(item.summary),
    content: text(item.content),
    cover: text(item.cover ?? item.cover_url ?? item.coverUrl ?? item.image) || undefined,
    category,
    tags: rawTags.length > 0 ? rawTags.map(mapTag) : tags.filter((tag) => rawTagIds.map((id) => number(id, -1)).includes(tag.id)),
    status: number(item.status, 1) === 0 ? "draft" : number(item.status, 1) === 2 ? "offline" : "published",
    updatedAt: formatDate(item.updatedAt ?? item.updateTime ?? item.update_time ?? item.published_time ?? item.publishedTime ?? item.date ?? item.create_time),
    publishedAt: formatDate(item.publishedAt ?? item.publishedTime ?? item.published_time ?? item.date ?? item.create_time),
    deleted: boolean(item.deleted ?? item.is_deleted ?? item.isDeleted),
    commentCount: number(item.commentCount ?? item.comment_count),
    isTop: boolean(item.isTop ?? item.is_top),
  };
}

export function mapComment(value: unknown): Comment {
  const item = asRecord(value);
  const targetId = number(item.articleId ?? item.target_id);
  const parentId = number(item.parentId ?? item.parent_id);
  const replyToId = number(item.replyToId ?? item.reply_to_id ?? parentId);
  const replyToNickname = text(item.replyToNickname ?? item.reply_to_nickname);
  const children = asArray(item.children ?? item.replies).map(mapComment);
  const isReply = parentId > 0;

  return {
    id: number(item.id),
    articleId: targetId,
    articleTitle: text(item.articleTitle),
    parentId,
    nickname: text(item.nickname, "匿名"),
    email: text(item.email),
    website: text(item.website) || undefined,
    avatar: text(item.avatar) || undefined,
    content: text(item.content),
    createdAt: formatDateTime(commentTime(item)),
    deleted: Boolean(item.deleted ?? item.is_deleted),
    browser: text(item.browser),
    os: text(item.os),
    ip: text(item.ip),
    location: text(item.location ?? item.ip_location ?? item.ipLocation),
    userAgent: text(item.userAgent ?? item.user_agent),
    replyTo: isReply && replyToNickname ? { id: replyToId, nickname: replyToNickname } : undefined,
    children,
  };
}

export function mapGuestbookMessage(value: unknown): GuestbookMessage {
  const item = asRecord(value);
  const parentId = String(item.parentId ?? item.parent_id ?? "");
  const replyToId = String(item.replyToId ?? item.reply_to_id ?? parentId);
  const replyToNickname = text(item.replyToNickname ?? item.reply_to_nickname);
  const children = asArray(item.children ?? item.replies).map(mapGuestbookMessage);
  const isReply = parentId !== "" && parentId !== "0";

  return {
    id: String(item.id ?? ""),
    parentId,
    nickname: text(item.nickname, "匿名"),
    email: text(item.email),
    website: text(item.website) || undefined,
    avatar: text(item.avatar) || undefined,
    content: text(item.content),
    createdAt: formatDateTime(commentTime(item)),
    deleted: Boolean(item.deleted ?? item.is_deleted),
    browser: text(item.browser),
    browserVersion: text(item.browserVersion),
    os: text(item.os),
    device: text(item.device),
    ip: text(item.ip),
    location: text(item.location ?? item.ip_location ?? item.ipLocation),
    userAgent: text(item.userAgent ?? item.user_agent),
    replyTo: isReply && replyToNickname ? { id: replyToId, nickname: replyToNickname } : undefined,
    children,
  };
}

export function mapFriendLink(value: unknown): FriendLink {
  const item = asRecord(value);

  return {
    id: number(item.id),
    name: text(item.name, "未命名站点"),
    url: text(item.url),
    logo: text(item.logo ?? item.logo_url ?? item.logoUrl) || undefined,
    description: text(item.description) || undefined,
    sort: number(item.sort),
    status: number(item.status, 1) === 1 ? 1 : 0,
  };
}
