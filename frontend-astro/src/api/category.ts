import { http, unwrap } from "./http";

export interface CategoryPayload {
  category_name: string;
  slug?: string;
  sort: number;
}

export function getCategories() {
  return unwrap(http.get("/api/categories"));
}

export function getAdminCategories() {
  return unwrap(http.post("/api/admin/categories"));
}

export function createCategory(payload: CategoryPayload) {
  return unwrap(http.post("/api/admin/categories/create", payload));
}

export function updateCategoryStatus(id: number | string, status: 0 | 1) {
  return unwrap(http.put(`/api/admin/categories/${id}/status`, { status }));
}

export function updateCategorySort(id: number | string, sort: number) {
  return unwrap(http.put(`/api/admin/categories/${id}/sort`, { sort }));
}
