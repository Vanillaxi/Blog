import { useEffect, useState } from "react";
import { marked } from "marked";
import { getArticleDetail } from "@/api/article";
import { getCategories } from "@/api/category";
import { getTags } from "@/api/tag";
import type { Article, Category, Tag } from "@/lib/types";
import { getListPayload, mapArticle, mapCategory, mapTag } from "@/lib/publicMappers";
import { CommentSection } from "./CommentSection";

type DetailState =
  | { status: "loading"; article?: never; message?: never }
  | { status: "missing"; article?: never; message: string }
  | { status: "error"; article?: never; message: string }
  | { status: "ready"; article: Article; message?: never };

function isNotFoundError(error: unknown) {
  const maybeError = error as { response?: { status?: number }; message?: string };
  const status = maybeError.response?.status;
  const message = maybeError.message ?? "";

  return status === 404 || message.includes("不存在") || message.includes("删除") || message.includes("not found");
}

function isNetworkError(error: unknown) {
  const maybeError = error as { response?: unknown; request?: unknown; message?: string };
  return Boolean(maybeError.request && !maybeError.response) || maybeError.message?.includes("网络");
}

export function ArticleDetail() {
  const [state, setState] = useState<DetailState>({ status: "loading" });

  useEffect(() => {
    let ignore = false;

    async function loadArticle() {
      const id = new URLSearchParams(window.location.search).get("id");

      if (!id) {
        setState({ status: "missing", message: "文章不存在" });
        document.title = "文章不存在 - Vanillaxi";
        return;
      }

      setState({ status: "loading" });

      try {
        const [articleResponse, categoryResponse, tagResponse] = await Promise.all([
          getArticleDetail(id),
          getCategories(),
          getTags(),
        ]);
        if (ignore) return;

        const categories = getListPayload<Category>(categoryResponse.data).map(mapCategory);
        const tags = getListPayload<Tag>(tagResponse.data).map(mapTag);
        const article = mapArticle(articleResponse.data, categories, tags);

        if (!article.id || article.deleted) {
          setState({ status: "missing", message: "文章不存在或已被删除" });
          document.title = "文章不存在 - Vanillaxi";
          return;
        }

        setState({ status: "ready", article });
        document.title = `${article.title} - Vanillaxi`;
      } catch (error) {
        if (ignore) return;

        const message = isNotFoundError(error) ? "文章不存在或已被删除" : isNetworkError(error) ? "网络错误" : "网络错误";
        setState({ status: isNotFoundError(error) ? "missing" : "error", message });
        document.title = `${message} - Vanillaxi`;
      }
    }

    void loadArticle();

    return () => {
      ignore = true;
    };
  }, []);

  if (state.status === "loading") {
    return <p className="mx-auto max-w-4xl rounded-2xl border border-white/50 bg-white/40 px-4 py-3 text-sm text-slate-500">文章加载中...</p>;
  }

  if (state.status === "missing" || state.status === "error") {
    return <p className="mx-auto max-w-4xl rounded-2xl border border-rose-100 bg-rose-50/80 px-4 py-3 text-sm text-rose-600">{state.message}</p>;
  }

  const { article } = state;
  const articleHtml = marked.parse(article.content, { async: false });

  return (
    <>
      <article className="glass-card overflow-hidden rounded-[2rem] p-4 sm:p-5">
        {article.cover && (
          <div className="relative aspect-[21/9] overflow-hidden rounded-[1.5rem] shadow-soft">
            <img src={article.cover} alt="" className="h-full w-full object-cover opacity-90" />
            <div className="absolute inset-0 bg-gradient-to-t from-white/30 via-transparent to-white/10" />
          </div>
        )}

        <div className="px-2 py-7 sm:px-5">
          <div className="mb-5 flex flex-wrap items-center gap-2 text-xs text-slate-500">
            <span className="rounded-full bg-white/50 px-2.5 py-1">{article.category.name}</span>
            <span>{article.publishedAt || article.updatedAt}</span>
          </div>

          <h1 className="font-serif text-4xl italic tracking-tight text-[#24314f] sm:text-5xl">{article.title}</h1>

          {article.summary && <p className="mt-4 max-w-3xl text-sm leading-7 text-slate-600">{article.summary}</p>}

          {article.tags.length > 0 && (
            <div className="mt-5 flex flex-wrap gap-2">
              {article.tags.map((tag) => (
                <span key={tag.id} className="rounded-full border border-white/50 bg-white/40 px-3 py-1 text-xs text-slate-500">
                  #{tag.name}
                </span>
              ))}
            </div>
          )}

          <div className="mt-9 border-t border-white/55 pt-8">
            <div
              className="article-markdown prose prose-neutral max-w-none"
              dangerouslySetInnerHTML={{ __html: articleHtml }}
            />
          </div>
        </div>
      </article>

      <CommentSection articleId={article.id} />
    </>
  );
}
