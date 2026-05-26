import type { Article } from "@/lib/types";

export function ArticleCard({ article }: { article: Article }) {
  return (
    <a
      href={`/articles/detail?id=${article.id}`}
      className="glass-card group block overflow-hidden rounded-[2rem] p-0 transition duration-200 hover:-translate-y-0.5 hover:bg-white/60"
    >
      <div className="relative h-56 overflow-hidden sm:h-72">
        {article.cover ? (
          <img
            src={article.cover}
            alt=""
            className="h-full w-full object-cover opacity-95 transition duration-700 group-hover:scale-[1.025]"
            loading="lazy"
          />
        ) : (
          <div className="h-full w-full bg-[radial-gradient(circle_at_20%_20%,rgba(251,207,232,0.55),transparent_34%),linear-gradient(135deg,rgba(244,231,235,0.95),rgba(220,234,244,0.88))]" />
        )}

        <div className="absolute inset-0 bg-gradient-to-t from-slate-950/50 via-slate-800/10 to-white/10" />
        <div className="absolute inset-x-0 bottom-0 p-5 sm:p-7">
          <div className="mb-3 flex flex-wrap items-center gap-2 text-xs text-white/82">
            <span className="rounded-full border border-white/35 bg-white/20 px-2.5 py-1 backdrop-blur-md">{article.category.name}</span>
            <span>{article.updatedAt}</span>
            <span>·</span>
            <span>{article.commentCount ?? 0} 条评论</span>
          </div>
          <h2 className="max-w-3xl font-serif text-3xl italic tracking-tight text-white drop-shadow-[0_8px_24px_rgba(15,23,42,0.34)] sm:text-4xl">
            {article.title}
          </h2>
        </div>
      </div>

      <div className="p-5 sm:p-7">
        <p className="line-clamp-2 max-w-3xl text-sm leading-7 text-slate-600">{article.summary}</p>
        <div className="mt-5 flex flex-wrap gap-2">
          {article.tags.map((tag) => (
            <span key={tag.id} className="rounded-full border border-white/50 bg-white/40 px-3 py-1 text-xs text-slate-500">
              #{tag.name}
            </span>
          ))}
        </div>
      </div>
    </a>
  );
}
