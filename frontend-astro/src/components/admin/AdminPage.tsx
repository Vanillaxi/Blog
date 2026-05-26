import { AdminGuard } from "./AdminGuard";
import { AdminShell } from "./AdminShell";
import { ArticleEditor } from "./ArticleEditor";
import { ArticleManager } from "./ArticleManager";
import { Dashboard } from "./Dashboard";
import { ArticleCommentManager, CategoryManager, CommentManager, FriendLinkManager, GuestbookManager, TagManager } from "./SimpleManagers";

type AdminPageKind =
  | "dashboard"
  | "articles"
  | "article-new"
  | "article-edit"
  | "categories"
  | "tags"
  | "comments"
  | "comment-detail"
  | "guestbook"
  | "friends";

const meta: Record<AdminPageKind, { title: string; subtitle: string }> = {
  dashboard: { title: "Dashboard", subtitle: "今天也慢慢写一点。" },
  articles: { title: "文章", subtitle: "写作、整理和发布。" },
  "article-new": { title: "写文章", subtitle: "先写下来，再慢慢修。" },
  "article-edit": { title: "编辑文章", subtitle: "调整内容、状态和发布设置。" },
  categories: { title: "分类", subtitle: "保持文章结构清楚。" },
  tags: { title: "Tags", subtitle: "用轻标签连接文章。" },
  comments: { title: "评论", subtitle: "按文章查看和管理读者评论。" },
  "comment-detail": { title: "文章评论", subtitle: "删除或恢复这篇文章下的评论。" },
  guestbook: { title: "留言", subtitle: "整理留言板里的访客消息。" },
  friends: { title: "友链", subtitle: "维护朋友和常看的站点。" },
};

export function AdminPage({ page, articleId }: { page: AdminPageKind; articleId?: number }) {
  const current = meta[page];

  return (
    <AdminGuard>
      <AdminShell title={current.title} subtitle={current.subtitle}>
        {page === "dashboard" && <Dashboard />}
        {page === "articles" && <ArticleManager />}
        {page === "article-new" && <ArticleEditor />}
        {page === "article-edit" && <ArticleEditor articleId={articleId} />}
        {page === "categories" && <CategoryManager />}
        {page === "tags" && <TagManager />}
        {page === "comments" && <CommentManager />}
        {page === "comment-detail" && articleId && <ArticleCommentManager articleId={articleId} />}
        {page === "guestbook" && <GuestbookManager />}
        {page === "friends" && <FriendLinkManager />}
      </AdminShell>
    </AdminGuard>
  );
}
