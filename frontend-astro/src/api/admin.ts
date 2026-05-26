import { getAdminArticles } from "./article";
import { getAdminCategories } from "./category";
import { getAdminFriendLinks } from "./friend";
import { getAdminTags } from "./tag";
import { getComments } from "./comment";
import { http, unwrap } from "./http";

export function getAdminDashboard() {
  return unwrap(http.get("/api/admin/dashboard"));
}

export function getDashboardData() {
  return Promise.all([
    getAdminArticles({ page: 1, pageSize: 10, status: -1, is_deleted: -1 }),
    getComments({ target_type: 1, target_id: 0, page: 1, pageSize: 10 }),
    getComments({ target_type: 2, target_id: 0, page: 1, pageSize: 10 }),
  ]);
}

export function getAdminBootstrapData() {
  return Promise.all([getAdminCategories(), getAdminTags(), getAdminFriendLinks()]);
}
