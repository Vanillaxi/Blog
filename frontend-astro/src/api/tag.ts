import { http, unwrap } from "./http";

export interface TagPayload {
  tag_name: string;
  slug: string;
  sort: number;
}

export function getTags() {
  return unwrap(http.get("/api/tags"));
}

export function getAdminTags() {
  return unwrap(http.get("/api/admin/tags"));
}

export function createTag(payload: TagPayload) {
  return unwrap(http.post("/api/admin/tag/create", payload));
}

export function updateTagStatus(id: number | string, status: 0 | 1) {
  return unwrap(http.put(`/api/admin/tag/${id}/status`, { status }));
}

export function updateTagSort(id: number | string, sort: number) {
  return unwrap(http.put(`/api/admin/tag/${id}/sort`, { sort }));
}
