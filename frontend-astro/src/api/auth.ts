import { http, unwrap } from "./http";

export interface LoginPayload {
  username: string;
  password: string;
}

export interface LoginData {
  token: string;
  username?: string;
}

export function login(payload: LoginPayload) {
  return unwrap<LoginData>(http.post("/api/admin/login", payload));
}
