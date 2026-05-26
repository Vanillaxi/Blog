import { useEffect, useState } from "react";
import { getArticlesTimeline, searchArticles } from "@/api/article";
import { getCategories } from "@/api/category";
import { getTags } from "@/api/tag";
import type { Article, Category, Tag } from "@/lib/types";
import { getListPayload, mapArticle, mapCategory, mapTag } from "@/lib/publicMappers";

type RawArticle = Record<string, unknown>;
const DEFAULT_COVER = "/images/sakura-mountain.png";

function isDev() {
  return import.meta.env.DEV;
}

function isPublishedArticle(item: RawArticle) {
  const status = Number(item.status ?? 1);
  const deleted = Number(item.is_deleted ?? item.isDeleted ?? 0);
  return status === 1 && deleted !== 1;
}

function getArticleTime(item: RawArticle) {
  const value = item.published_time ?? item.publishedTime ?? item.date ?? item.create_time ?? item.createTime ?? "";
  const timestamp = new Date(String(value)).getTime();
  return Number.isFinite(timestamp) ? timestamp : 0;
}

function byPublishedTimeDesc(left: RawArticle, right: RawArticle) {
  return getArticleTime(right) - getArticleTime(left) || Number(right.id ?? 0) - Number(left.id ?? 0);
}

function mapArticleSafely(item: RawArticle, categories: Category[], tags: Tag[]) {
  try {
    return mapArticle(item, categories, tags);
  } catch (error) {
    console.error("[timeline] article mapping failed", error, item);
    return null;
  }
}

export function Timeline() {
  const [items, setItems] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let ignore = false;

    async function loadTimeline() {
      setLoading(true);
      setError("");
      try {
        const [timelineResponse, categoryResponse, tagResponse] = await Promise.all([
          getArticlesTimeline({ page: 1, pageSize: 50 }),
          getCategories(),
          getTags(),
        ]);
        const categories = getListPayload<Category>(categoryResponse.data).map(mapCategory);
        const tags = getListPayload<Tag>(tagResponse.data).map(mapTag);

        let rawArticles = getListPayload<RawArticle>(timelineResponse.data);

        if (isDev()) {
          console.log("[timeline] timeline response", timelineResponse);
          console.log("[timeline] parsed timeline list", rawArticles);
        }

        if (rawArticles.length === 0) {
          const fallbackResponse = await searchArticles({ keyword: "", page: 1, pageSize: 50 });
          rawArticles = getListPayload<RawArticle>(fallbackResponse.data);

          if (isDev()) {
            console.log("[timeline] fallback search response", fallbackResponse);
            console.log("[timeline] parsed fallback list", rawArticles);
          }
        }

        const articles = rawArticles
          .filter(isPublishedArticle)
          .sort(byPublishedTimeDesc)
          .map((item) => mapArticleSafely(item, categories, tags))
          .filter((article): article is Article => Boolean(article));

        if (isDev()) {
          console.log("[timeline] mapped articles", articles);
        }

        if (!ignore) setItems(articles);
      } catch (err) {
        console.error("[timeline] parse or request failed", err);
        if (!ignore) setError("加载时间轴失败");
      } finally {
        if (!ignore) setLoading(false);
      }
    }

    void loadTimeline();
    return () => {
      ignore = true;
    };
  }, []);

  if (loading) {
    return <p className="mx-auto max-w-4xl rounded-2xl border border-white/50 bg-white/40 px-4 py-3 text-sm text-slate-500">时间轴加载中...</p>;
  }

  if (error) {
    return <p className="mx-auto max-w-4xl rounded-2xl border border-rose-100 bg-rose-50/80 px-4 py-3 text-sm text-rose-600">{error}</p>;
  }

  if (items.length === 0) {
    return <p className="mx-auto max-w-4xl rounded-2xl border border-white/50 bg-white/40 px-4 py-3 text-sm text-slate-500">暂时还没有文章。</p>;
  }

  return (
    <div className="relative mx-auto max-w-5xl py-2 sm:py-6">
      <div className="absolute bottom-0 left-5 top-0 w-[3px] bg-gradient-to-b from-transparent via-[rgba(201,116,140,0.82)] to-transparent shadow-[0_0_14px_rgba(201,116,140,0.24)] md:left-1/2 md:-translate-x-1/2" />

      <div className="space-y-8 md:space-y-10">
        {items.map((article, index) => {
          const isLeft = index % 2 === 0;
          return (
            <div className="relative grid pl-12 md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)] md:pl-0" key={article.id}>
              <span className="absolute left-5 top-7 z-10 h-[18px] w-[18px] -translate-x-1/2 rounded-full border-[3px] border-white/95 bg-[#D97F97] shadow-[0_0_12px_rgba(217,127,151,0.32)] md:left-1/2" />

              <a
                href={`/articles/detail?id=${article.id}`}
                className={`glass-card group block w-full max-w-[460px] overflow-hidden rounded-3xl p-0 transition duration-200 hover:-translate-y-0.5 hover:bg-white/60 ${
                  isLeft ? "md:mr-[72px] md:justify-self-end" : "md:col-start-2 md:ml-[72px] md:justify-self-start"
                }`}
              >
                <div className="relative h-44 overflow-hidden rounded-t-3xl sm:h-52">
                  <img
                    src={article.cover || DEFAULT_COVER}
                    alt=""
                    className="h-full w-full object-cover opacity-95 transition duration-700 group-hover:scale-[1.025]"
                    loading="lazy"
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-white/25 via-transparent to-white/10" />
                </div>

                <div className="p-5 sm:p-6">
                  <p className="text-xs text-slate-500">{article.publishedAt}</p>
                  <h2 className="mt-2 font-serif text-2xl italic leading-tight text-[#25324f]">{article.title}</h2>
                  <p className="mt-3 line-clamp-3 text-sm leading-7 text-slate-600">{article.summary}</p>
                  <div className="mt-4 flex flex-wrap gap-2">
                    <span className="rounded-full border border-white/45 bg-white/40 px-3 py-1 text-xs text-slate-500">{article.category.name}</span>
                    {article.tags.slice(0, 2).map((tag) => (
                      <span className="rounded-full border border-white/40 bg-white/30 px-3 py-1 text-xs text-slate-500" key={tag.id}>
                        #{tag.name}
                      </span>
                    ))}
                  </div>
                </div>
              </a>
            </div>
          );
        })}
      </div>
    </div>
  );
}
