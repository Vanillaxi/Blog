import { Search } from "lucide-react";
import { useEffect, useState } from "react";

const ADMIN_PHRASE = "vanillaxi-admin-2007";

const navItems = [
  { label: "主页", href: "/" },
  { label: "文章", href: "/articles" },
  { label: "留言板", href: "/guestbook" },
  { label: "时间轴", href: "/timeline" },
  { label: "友链", href: "/friends" },
  { label: "关于我", href: "/about" },
];

export function PublicNav() {
  const [keyword, setKeyword] = useState("");
  const [pathname, setPathname] = useState("/");
  const isHome = pathname === "/";

  useEffect(() => {
    setPathname(window.location.pathname);
  }, []);

  function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const value = keyword.trim();

    if (value === ADMIN_PHRASE) {
      window.location.href = "/admin/login";
      return;
    }

    if (value.length > 0) {
      window.location.href = `/articles?keyword=${encodeURIComponent(value)}`;
    }
  }

  return (
    <header className="fixed left-0 right-0 top-0 z-40 px-4 pt-5 sm:px-8 sm:pt-8 lg:px-14">
      <div className="mx-auto flex max-w-7xl flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <form
          onSubmit={onSubmit}
          className={`group flex h-11 w-full items-center gap-2 rounded-full border px-4 text-slate-700 shadow-sm backdrop-blur-xl transition-all focus-within:border-rose-200/80 focus-within:ring-2 focus-within:ring-rose-200/40 sm:h-12 sm:w-60 ${
            isHome ? "border-white/60 bg-white/40" : "border-white/50 bg-white/30"
          }`}
        >
          <Search className="h-4 w-4 shrink-0 text-slate-500" strokeWidth={1.7} />
          <input
            value={keyword}
            onChange={(event) => setKeyword(event.target.value)}
            placeholder="Search..."
            className="min-w-0 flex-1 bg-transparent text-sm text-slate-800 caret-rose-400 placeholder:text-slate-400 outline-none"
            aria-label="搜索文章"
          />
        </form>

        <nav
          className={`no-scrollbar flex max-w-full items-center gap-1 overflow-x-auto rounded-full border p-1 text-sm text-slate-700 shadow-[0_14px_45px_rgba(148,163,184,0.12)] backdrop-blur-2xl ${
            isHome ? "border-white/60 bg-white/42" : "border-white/50 bg-white/30"
          }`}
        >
          {navItems.map((item) => {
            const active = item.href === "/" ? pathname === "/" : pathname.startsWith(item.href);
            return (
              <a
                key={item.href}
                href={item.href}
                className={`shrink-0 rounded-full border px-3.5 py-2 transition duration-200 ${
                  active
                    ? isHome
                      ? "border-rose-200/70 bg-gradient-to-r from-rose-200/70 to-sky-100/70 text-slate-800 shadow-[0_8px_24px_rgba(255,182,193,0.28)]"
                      : "border-rose-200/55 bg-gradient-to-r from-rose-200/55 to-sky-100/55 text-slate-800 shadow-[0_8px_22px_rgba(255,182,193,0.14)]"
                    : isHome
                      ? "border-transparent text-slate-700 hover:bg-white/62 hover:text-slate-900"
                      : "border-transparent text-slate-700 hover:bg-white/50 hover:text-slate-900"
                }`}
              >
                {item.label}
              </a>
            );
          })}
        </nav>
      </div>
    </header>
  );
}
