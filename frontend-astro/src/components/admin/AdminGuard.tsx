import { useEffect, useState } from "react";
import { isAdminLoggedIn } from "@/lib/adminAuth";

export function AdminGuard({ children }: { children: React.ReactNode }) {
  const [ready, setReady] = useState(false);

  useEffect(() => {
    if (!isAdminLoggedIn()) {
      window.location.href = "/admin/login";
      return;
    }
    setReady(true);
  }, []);

  if (!ready) return <div className="p-8 text-sm text-slate-500">正在确认登录状态...</div>;
  return <>{children}</>;
}
