import { useState } from "react";
import { login } from "@/api/auth";
import { setAdminToken } from "@/lib/adminAuth";

export function LoginPanel() {
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    setLoading(true);
    setError("");
    try {
      const result = await login({
        username: String(form.get("username") || ""),
        password: String(form.get("password") || ""),
      });
      setAdminToken(result.data.token);
      window.location.href = "/admin";
    } catch (err) {
      setError(err instanceof Error ? err.message : "登录失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={submit} className="glass-card mx-auto w-full max-w-md rounded-[2rem] p-8">
      <h1 className="font-serif text-4xl italic text-slate-950">Vanillaxi</h1>
      <p className="mt-3 text-sm text-slate-500">进入你的写作后台。</p>
      <div className="mt-8 space-y-4">
        <input name="username" placeholder="用户名" className="w-full rounded-2xl border border-white/70 bg-white/75 px-4 py-3 outline-none" />
        <input name="password" type="password" placeholder="密码" className="w-full rounded-2xl border border-white/70 bg-white/75 px-4 py-3 outline-none" />
      </div>
      {error && <p className="mt-4 text-sm text-rose-500">{error}</p>}
      <button disabled={loading} className="mt-6 w-full rounded-full bg-slate-950 px-5 py-3 text-sm text-white disabled:opacity-50">{loading ? "登录中..." : "登录"}</button>
    </form>
  );
}
