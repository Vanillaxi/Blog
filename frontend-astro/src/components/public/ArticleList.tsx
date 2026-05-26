import { useEffect, useMemo, useState } from "react";
import { getArticlesByCategory, getArticlesByTag, searchArticles } from "@/api/article";
import { getCategories } from "@/api/category";
import { getTags } from "@/api/tag";
import type { Article, Category, Tag } from "@/lib/types";
import { getListPayload, mapArticle, mapCategory, mapTag } from "@/lib/publicMappers";
import { ArticleCard } from "./ArticleCard";
import { ArticleFilters } from "./ArticleFilters";
import { EmptyState } from "./EmptyState";

function mapArticleSafely(item: unknown, categories: Category[], tags: Tag[]) {
  try {
    return mapArticle(item, categories, tags);
  } catch (error) {
    console.error("[articles] article mapping failed", error, item);
    return null;
  }
}

export function ArticleList() {
  const params = typeof window === "undefined" ? new URLSearchParams() : new URLSearchParams(window.location.search);
  const keyword = params.get("keyword") || "";
  const initialCategory = params.get("category") || "全部";
  const initialTag = params.get("tag") || "";

  const [articles, setArticles] = useState<Article[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [selectedCategory, setSelectedCategory] = useState(initialCategory);
  const [selectedTag, setSelectedTag] = useState(initialTag);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const categoryOptions = useMemo(() => ["全部", ...categories.map((category) => category.name)], [categories]);
  const tagOptions = useMemo(() => tags.map((tag) => tag.name), [tags]);

  function updateCategory(value: string) {
    setError("");
    setSelectedCategory(value);
  }

  function updateTag(value: string) {
    setError("");
    setSelectedTag(value);
  }

  useEffect(() => {
    let ignore = false;

    async function loadFilters() {
      try {
        const [categoryResponse, tagResponse] = await Promise.all([getCategories(), getTags()]);
        if (ignore) return;
        const nextCategories = getListPayload(categoryResponse.data).map(mapCategory);
        const nextTags = getListPayload(tagResponse.data).map(mapTag);
        setCategories(nextCategories);
        setTags(nextTags);
      } catch (err) {
        if (!ignore) setError(err instanceof Error ? err.message : "筛选项加载失败");
      }
    }

    void loadFilters();
    return () => {
      ignore = true;
    };
  }, []);

  useEffect(() => {
    let ignore = false;

    async function loadArticles() {
      setLoading(true);
      setError("");

      try {
        const hasCategoryFilter = selectedCategory !== "全部";
        const category = categories.find((item) => item.name === selectedCategory);
        const tag = tags.find((item) => item.name === selectedTag);
        const params = { page: 1, pageSize: 50 };
        const response = tag && hasCategoryFilter
          ? await searchArticles({ ...params, keyword })
          : tag
          ? await getArticlesByTag(tag.id, params)
          : category && hasCategoryFilter
            ? await getArticlesByCategory(category.id, params)
            : await searchArticles({ ...params, keyword });

        if (ignore) return;
        const nextArticles = getListPayload(response.data)
          .map((item) => mapArticleSafely(item, categories, tags))
          .filter((article): article is Article => Boolean(article));
        setArticles(nextArticles);
      } catch (err) {
        if (!ignore) setError(err instanceof Error ? err.message : "文章列表加载失败");
      } finally {
        if (!ignore) setLoading(false);
      }
    }

    void loadArticles();
    return () => {
      ignore = true;
    };
  }, [categories, keyword, selectedCategory, selectedTag, tags]);

  const filtered = useMemo(() => {
    const value = keyword.trim().toLowerCase();
    return articles.filter((article) => {
      if (value && !article.title.toLowerCase().includes(value) && !article.summary.toLowerCase().includes(value)) {
        return false;
      }
      if (selectedCategory !== "全部" && article.category.name !== selectedCategory) {
        return false;
      }
      if (selectedTag && !article.tags.some((tag) => tag.name === selectedTag)) {
        return false;
      }
      return true;
    });
  }, [articles, keyword, selectedCategory, selectedTag]);

  return (
    <div className="space-y-6">
      <ArticleFilters
        categories={categoryOptions}
        tags={tagOptions}
        selectedCategory={selectedCategory}
        selectedTag={selectedTag}
        onCategoryChange={updateCategory}
        onTagChange={updateTag}
      />

      {error && <p className="rounded-2xl border border-rose-100 bg-rose-50/80 px-4 py-3 text-sm text-rose-600">{error}</p>}
      {loading && <p className="rounded-2xl border border-white/50 bg-white/40 px-4 py-3 text-sm text-slate-500">文章加载中...</p>}

      {(keyword || selectedCategory !== "全部" || selectedTag) && (
        <p className="rounded-full border border-white/50 bg-white/40 px-4 py-2 text-sm text-slate-500 backdrop-blur">
          当前筛选
          {keyword && <span> · 关键词「{keyword}」</span>}
          {selectedCategory !== "全部" && <span> · 分类「{selectedCategory}」</span>}
          {selectedTag && <span> · Tag「{selectedTag}」</span>}
          <span> · 共 {filtered.length} 篇</span>
        </p>
      )}

      <div className="space-y-6">
        {filtered.map((article) => (
          <ArticleCard key={article.id} article={article} />
        ))}
        {!loading && filtered.length === 0 && <EmptyState />}
      </div>
    </div>
  );
}
