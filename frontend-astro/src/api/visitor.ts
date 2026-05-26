import { http, unwrap } from "./http";

export interface VisitorIdentityResponse {
  exists: boolean;
  nickname?: string;
}

export function getVisitorIdentity(email: string) {
  return unwrap<VisitorIdentityResponse>(http.get("/api/visitor-identity", { params: { email } }));
}
