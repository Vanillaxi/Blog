export type ArticleStatus = "draft" | "published" | "offline" | "deleted";

export interface Category {
  id: number;
  name: string;
  slug: string;
  count: number;
  sort: number;
}

export interface Tag {
  id: number;
  name: string;
  slug: string;
  sort: number;
}

export interface Article {
  id: number;
  title: string;
  summary: string;
  content: string;
  cover?: string;
  category: Category;
  tags: Tag[];
  status: ArticleStatus;
  updatedAt: string;
  publishedAt?: string;
  deleted?: boolean;
  commentCount?: number;
  isTop?: boolean;
}

export interface Comment {
  id: number;
  articleId: number;
  articleTitle: string;
  parentId: number;
  nickname: string;
  email: string;
  website?: string;
  avatar?: string;
  content: string;
  createdAt: string;
  deleted?: boolean;
  browser?: string;
  os?: string;
  ip?: string;
  location?: string;
  userAgent?: string;
  replyTo?: {
    id: number;
    nickname: string;
  };
  children?: Comment[];
}

export interface GuestbookMessage {
  id: string;
  parentId: string;
  nickname: string;
  email: string;
  website?: string;
  avatar?: string;
  content: string;
  createdAt: string;
  deleted?: boolean;
  browser?: string;
  browserVersion?: string;
  os?: string;
  device?: string;
  ip?: string;
  location?: string;
  userAgent?: string;
  replyTo?: {
    id: string;
    nickname: string;
  };
  children?: GuestbookMessage[];
}

export interface FriendLink {
  id: number;
  name: string;
  url: string;
  logo?: string;
  description?: string;
  sort: number;
  status: 0 | 1;
}
