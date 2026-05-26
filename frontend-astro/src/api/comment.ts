import { http, unwrap } from "./http";

export interface CommentQuery {
  target_type: 1 | 2;
  target_id: number;
  page?: number;
  pageSize?: number;
}

export interface AdminCommentQuery {
  target_type?: 1 | 2;
  target_id?: number;
  include_deleted?: boolean;
  page?: number;
  page_size?: number;
}

export interface CommentPayload {
  target_type: 1 | 2;
  target_id: number;
  parent_id?: number;
  nickname: string;
  email: string;
  website?: string;
  content: string;
}

export function getComments(params: CommentQuery) {
  return unwrap(http.get("/api/cmments/get", { params }));
}

export function addComment(payload: CommentPayload) {
  return unwrap(http.post("/api/comments/add", payload));
}

export function deleteComment(id: number | string) {
  return unwrap(http.delete(`/api/admin/comments/delete/${id}`));
}

export function getAdminComments(params?: AdminCommentQuery) {
  return unwrap(http.get("/api/admin/comments", { params }));
}

export function deleteAdminComment(id: number | string) {
  return deleteComment(id);
}

export function restoreAdminComment(id: number | string) {
  return unwrap(http.put(`/api/admin/comments/${id}/restore`));
}

export function getAdminGuestbookMessages(params?: Omit<AdminCommentQuery, "target_type">) {
  return getAdminComments({
    ...params,
    target_type: 2,
  });
}
