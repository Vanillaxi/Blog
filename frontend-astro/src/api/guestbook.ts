import { addComment, getComments, type CommentPayload } from "./comment";

export type GuestbookPayload = Omit<CommentPayload, "target_type" | "target_id">;

export function getGuestbookMessages(page?: number, pageSize?: number) {
  return getComments({
    target_type: 2,
    target_id: 0,
    page,
    pageSize,
  });
}

export function submitGuestbook(payload: GuestbookPayload) {
  return addComment({
    ...payload,
    target_type: 2,
    target_id: 0,
  });
}
