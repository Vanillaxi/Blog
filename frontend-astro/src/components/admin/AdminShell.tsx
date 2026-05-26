import { BookOpen, Folder, Hash, Home, Link2, LogOut, MessageCircle, MessagesSquare } from "lucide-react";
import { clearAdminToken } from "@/lib/adminAuth";

const nav = [
  { href: "/admin", label: "Dashboard", icon: Home },
  { href: "/admin/articles", label: "文章", icon: BookOpen },
  { href: "/admin/categories", label: "分类", icon: Folder },
  { href: "/admin/tags", label: "Tags", icon: Hash },
  { href: "/admin/comments", label: "评论", icon: MessageCircle },
  { href: "/admin/guestbook", label: "留言", icon: MessagesSquare },
  { href: "/admin/friends", label: "友链", icon: Link2 },
];

export function AdminShell({ title, subtitle, children }: { title: string; subtitle?: string; children: React.ReactNode }) {
  function logout() {
    clearAdminToken();
    window.location.href = "/admin/login";
  }

  return (
    <div className="min-h-screen px-4 py-4 text-slate-900 lg:px-6">
      <div className="mx-auto grid max-w-7xl gap-5 lg:grid-cols-[220px_1fr]">
        <aside className="rounded-[2rem] border border-white/70 bg-white/65 p-4 shadow-soft backdrop-blur-xl lg:sticky lg:top-4 lg:h-[calc(100vh-2rem)]">
          <a href="/" className="block px-3 py-3 font-serif text-2xl italic transition hover:opacity-75">Vanillaxi</a>
          <nav className="mt-6 space-y-1">
            {nav.map((item) => {
              const Icon = item.icon;
              return (
                <a key={item.href} href={item.href} className="flex items-center gap-3 rounded-2xl px-3 py-2.5 text-sm text-slate-600 transition hover:bg-white hover:text-slate-950">
                  <Icon className="h-4 w-4" />
                  {item.label}
                </a>
              );
            })}
          </nav>
          <button onClick={logout} className="mt-8 flex w-full items-center gap-3 rounded-2xl px-3 py-2.5 text-sm text-slate-500 hover:bg-white hover:text-slate-950">
            <LogOut className="h-4 w-4" />
            退出登录
          </button>
        </aside>
        <main className="min-w-0">
          <header className="mb-6 flex flex-wrap items-end justify-between gap-4 rounded-[2rem] border border-white/70 bg-white/55 p-6 shadow-soft backdrop-blur-xl">
            <div>
              <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
              {subtitle && <p className="mt-2 text-sm text-slate-500">{subtitle}</p>}
            </div>
          </header>
          {children}
        </main>
      </div>
    </div>
  );
}
