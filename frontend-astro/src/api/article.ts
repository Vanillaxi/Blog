import { http, unwrap } from "./http";

export interface PaginationParams {
  page?: number;
  pageSize?: number;
}

export interface ArticleQuery extends PaginationParams {
  status?: number;
  is_deleted?: number;
  category_id?: number;
  keyword?: string;
}

export interface ArticlePayload {
  title: string;
  category_id: number;
  summary: string;
  content: string;
  is_top: 0 | 1;
  cover_url: string;
  tag_ids: number[];
}

export type ArticleSaveStatus = 0 | 1 | 2;

export function getArticleDetail(id: number | string) {
  return unwrap(http.get(`/api/articles/${id}`));
}

export function getArticlesTimeline(params?: PaginationParams) {
  return unwrap(http.get("/api/articles/timeline", { params }));
}

export function searchArticles(params: PaginationParams & { keyword: string }) {
  return unwrap(http.get("/api/articles/search", { params }));
}

export function getArticlesByCategory(categoryId: number | string, params?: PaginationParams) {
  return unwrap(http.get(`/api/categories/${categoryId}/articles`, { params }));
}

export function getArticlesByTag(tagId: number | string, params?: PaginationParams) {
  return unwrap(http.get(`/api/tags/${tagId}/articles`, { params }));
}

export function getAdminArticles(params?: ArticleQuery) {
  return unwrap(http.get("/api/admin/articles", { params }));
}

export function getAdminArticleDetail(id: number | string) {
  return unwrap(http.get(`/api/admin/articles/${id}`));
}

export function createArticle(payload: ArticlePayload) {
  return unwrap(http.post("/api/admin/articles/create", payload));
}

export function updateArticle(id: number | string, payload: ArticlePayload) {
  return unwrap(http.put(`/api/admin/articles/update/${id}`, payload));
}

export function updateArticleStatus(id: number | string, status: ArticleSaveStatus) {
  return unwrap(http.put(`/api/admin/articles/updateStatus/${id}`, { status }));
}

export function deleteArticle(id: number | string) {
  return unwrap(http.delete(`/api/admin/articles/delete/${id}`));
}
