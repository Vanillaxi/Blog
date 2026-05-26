interface ArticleFiltersProps {
  categories: string[];
  tags: string[];
  selectedCategory: string;
  selectedTag: string;
  onCategoryChange: (category: string) => void;
  onTagChange: (tag: string) => void;
}

export function ArticleFilters({
  categories,
  tags,
  selectedCategory,
  selectedTag,
  onCategoryChange,
  onTagChange,
}: ArticleFiltersProps) {
  return (
    <section className="glass-card rounded-[2rem] p-5 sm:p-6">
      <div className="flex flex-wrap gap-2">
        {categories.map((category) => {
          const active = selectedCategory === category;
          return (
            <button
              key={category}
              type="button"
              onClick={() => onCategoryChange(category)}
              className={`rounded-full px-4 py-2 text-sm transition ${
                active
                  ? "border border-white/55 bg-gradient-to-r from-rose-200/55 via-pink-100/50 to-sky-100/50 text-slate-800 shadow-[0_8px_22px_rgba(244,114,182,0.12)]"
                  : "border border-white/50 bg-white/35 text-slate-600 hover:bg-white/50"
              }`}
            >
              {category}
            </button>
          );
        })}
      </div>

      <div className="mt-4 flex flex-wrap gap-2">
        {tags.map((tag) => {
          const active = selectedTag === tag;
          return (
            <button
              key={tag}
              type="button"
              onClick={() => onTagChange(active ? "" : tag)}
              className={`rounded-full px-3 py-1.5 text-xs transition ${
                active
                  ? "border border-rose-200/60 bg-rose-100/60 text-rose-900"
                  : "border border-white/35 bg-white/25 text-slate-500 hover:bg-white/45 hover:text-slate-700"
              }`}
            >
              #{tag}
            </button>
          );
        })}
      </div>
    </section>
  );
}
