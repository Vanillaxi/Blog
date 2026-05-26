import type { ArticleStatus, FriendLink } from "./types";

export interface AdminCategory {
  id: number;
  name: string;
  slug: string;
  count: number;
  status: 0 | 1;
  sort: number;
}

export interface AdminTag {
  id: number;
  name: string;
  slug: string;
  count: number;
  status: 0 | 1;
  sort: number;
}

export interface AdminArticle {
  id: number;
  title: string;
  summary: string;
  cover?: string;
  categoryId: number;
  commentCount: number;
  status: ArticleStatus;
  isTop: boolean;
  deleted: boolean;
  publishedAt?: string;
  createdAt?: string;
  updatedAt?: string;
  tagIds: number[];
}

export interface AdminComment {
  id: number;
  targetType: 1 | 2;
  targetId: number;
  articleTitle: string;
  parentId: number;
  nickname: string;
  email: string;
  website?: string;
  content: string;
  ip: string;
  location: string;
  userAgent: string;
  deleted: boolean;
  createdAt?: string;
}

export interface AdminDashboardData {
  articleCount: number;
  publishedCount: number;
  draftCount: number;
  offlineCount: number;
  deletedArticleCount: number;
  categoryCount: number;
  tagCount: number;
  commentCount: number;
  guestbookCount: number;
  friendlinkCount: number;
  recentArticles: AdminArticle[];
}

type RecordLike = Record<string, unknown>;

function asRecord(value: unknown): RecordLike {
  return value && typeof value === "object" ? (value as RecordLike) : {};
}

function numberValue(value: unknown, fallback = 0) {
  const next = Number(value);
  return Number.isFinite(next) ? next : fallback;
}

function stringValue(value: unknown) {
  return typeof value === "string" ? value : value == null ? "" : String(value);
}

function statusValue(value: unknown): 0 | 1 {
  return numberValue(value, 1) === 0 ? 0 : 1;
}

function formatDate(value: unknown) {
  const raw = stringValue(value);
  if (!raw) return undefined;
  const date = new Date(raw);
  if (Number.isNaN(date.getTime())) return raw;
  return date.toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function mapArticleStatus(status: unknown, isDeleted: unknown): ArticleStatus {
  if (numberValue(isDeleted) === 1) return "deleted";
  if (numberValue(status) === 0) return "draft";
  if (numberValue(status) === 1) return "published";
  if (numberValue(status) === 2) return "offline";
  return "draft";
}

function mapTagIds(value: unknown) {
  if (!Array.isArray(value)) return [];
  return value
    .map((item) => {
      const record = asRecord(item);
      return numberValue(record.id ?? record.tag_id ?? item, NaN);
    })
    .filter((id) => Number.isFinite(id));
}

export function mapAdminArticle(value: unknown): AdminArticle {
  const item = asRecord(value);
  return {
    id: numberValue(item.id),
    title: stringValue(item.title),
    summary: stringValue(item.summary),
    cover: stringValue(item.cover_url) || undefined,
    categoryId: numberValue(item.category_id),
    commentCount: numberValue(item.comment_count),
    status: mapArticleStatus(item.status, item.is_deleted),
    isTop: numberValue(item.is_top) === 1,
    deleted: numberValue(item.is_deleted) === 1,
    publishedAt: formatDate(item.published_time),
    createdAt: formatDate(item.create_time),
    updatedAt: formatDate(item.update_time),
    tagIds: mapTagIds(item.tags ?? item.tag_ids),
  };
}

export function mapAdminCategory(value: unknown): AdminCategory {
  const item = asRecord(value);
  return {
    id: numberValue(item.id),
    name: stringValue(item.category_name ?? item.name),
    slug: stringValue(item.slug),
    count: numberValue(item.article_count ?? item.count),
    status: statusValue(item.status),
    sort: numberValue(item.sort),
  };
}

export function mapAdminTag(value: unknown): AdminTag {
  const item = asRecord(value);
  return {
    id: numberValue(item.id),
    name: stringValue(item.tag_name ?? item.name),
    slug: stringValue(item.slug),
    count: numberValue(item.article_count ?? item.count),
    status: statusValue(item.status),
    sort: numberValue(item.sort),
  };
}

export function mapFriendLink(value: unknown): FriendLink {
  const item = asRecord(value);
  return {
    id: numberValue(item.id),
    name: stringValue(item.name),
    url: stringValue(item.url),
    logo: stringValue(item.logo ?? item.logo_url ?? item.logoUrl) || undefined,
    description: stringValue(item.description) || undefined,
    sort: numberValue(item.sort),
    status: statusValue(item.status),
  };
}

export function mapAdminComment(value: unknown): AdminComment {
  const item = asRecord(value);
  return {
    id: numberValue(item.id),
    targetType: numberValue(item.target_type) === 2 ? 2 : 1,
    targetId: numberValue(item.target_id),
    articleTitle: stringValue(item.article_title),
    parentId: numberValue(item.parent_id),
    nickname: stringValue(item.nickname),
    email: stringValue(item.email),
    website: stringValue(item.website) || undefined,
    content: stringValue(item.content),
    ip: stringValue(item.ip),
    location: stringValue(item.ip_location ?? item.location),
    userAgent: stringValue(item.user_agent),
    deleted: numberValue(item.is_deleted) === 1,
    createdAt: formatDate(item.create_time),
  };
}

export function mapAdminDashboard(value: unknown): AdminDashboardData {
  const item = asRecord(value);
  return {
    articleCount: numberValue(item.article_count),
    publishedCount: numberValue(item.published_count),
    draftCount: numberValue(item.draft_count),
    offlineCount: numberValue(item.offline_count),
    deletedArticleCount: numberValue(item.deleted_article_count),
    categoryCount: numberValue(item.category_count),
    tagCount: numberValue(item.tag_count),
    commentCount: numberValue(item.comment_count),
    guestbookCount: numberValue(item.guestbook_count),
    friendlinkCount: numberValue(item.friendlink_count),
    recentArticles: Array.isArray(item.recent_articles) ? item.recent_articles.map(mapAdminArticle) : [],
  };
}

export function listData(value: unknown) {
  const data = asRecord(asRecord(value).data);
  const list = data.list ?? data.links ?? asRecord(value).data;
  return Array.isArray(list) ? list : [];
}
