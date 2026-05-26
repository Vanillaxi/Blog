import axios from "axios";
import { clearAdminToken, getAdminToken } from "./adminAuth";

export interface ApiResponse<T = unknown> {
  code?: number;
  data: T;
  msg?: string;
  success?: boolean;
}

const API_BASE_URL = import.meta.env.PROD ? "" : (import.meta.env.PUBLIC_API_BASE_URL ?? "");

export const http = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

http.interceptors.request.use((config) => {
  const origin = typeof window === "undefined" ? "http://localhost" : window.location.origin;
  const baseURL = config.baseURL || origin;
  const requestUrl = new URL(config.url || "", baseURL).pathname;
  const isLoginRequest = requestUrl === "/api/admin/login";

  if (isLoginRequest) {
    delete config.headers.Authorization;
    return config;
  }

  const token = getAdminToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

http.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status;

    if (status === 401) {
      clearAdminToken();
      if (typeof window !== "undefined" && window.location.pathname.startsWith("/admin")) {
        window.location.href = "/admin/login";
      }
    }

    if (!error.response) {
      error.message = "网络连接失败，请检查后端服务是否启动。";
    } else if (status >= 500) {
      error.message = error.response?.data?.msg || "服务器错误，请稍后再试。";
    } else {
      error.message = error.response?.data?.msg || error.message;
    }

    return Promise.reject(error);
  },
);

export async function unwrap<T>(request: Promise<{ data: ApiResponse<T> | T }>) {
  const response = await request;
  const payload = response.data;

  if (Array.isArray(payload)) {
    return { code: 200, data: payload as T, msg: "" };
  }

  if (!payload || typeof payload !== "object") {
    return { code: 200, data: payload as T, msg: "" };
  }

  const wrapped = payload as ApiResponse<T>;
  const hasWrappedData = "data" in wrapped;
  const code = wrapped.code;
  const isSuccess = code === 200 || code === 0 || wrapped.success === true || (code == null && hasWrappedData && !wrapped.msg);

  if (!isSuccess) {
    throw new Error(wrapped.msg || "请求失败");
  }
  return wrapped;
}
