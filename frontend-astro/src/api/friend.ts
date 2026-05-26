import { http, unwrap } from "./http";

export interface FriendLinkPayload {
  id?: number;
  name: string;
  url: string;
  logo?: string;
  description?: string;
  sort: number;
  status: 0 | 1;
}

export function getFriendLinks(params?: { page?: number; pageSize?: number }) {
  return unwrap(http.get("/api/friendlinks", { params }));
}

export function getAdminFriendLinks(params?: { page?: number; pageSize?: number }) {
  return unwrap(http.get("/api/admin/friendlinks", { params }));
}

export function createFriendLink(payload: FriendLinkPayload) {
  return unwrap(http.post("/api/admin/friendlinks/add", payload));
}

export const createAdminFriendLink = createFriendLink;

export function updateFriendLink(id: number | string, payload: FriendLinkPayload) {
  return unwrap(http.put(`/api/admin/friendlinks/${id}/update`, { ...payload, id: Number(id) }));
}

export const updateAdminFriendLink = updateFriendLink;

export function deleteFriendLink(id: number | string) {
  return unwrap(http.delete(`/api/admin/friendlinks/${id}/delete`));
}

export const deleteAdminFriendLink = deleteFriendLink;
