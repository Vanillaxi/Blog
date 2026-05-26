export function PublicAvatar({ name, avatar, size = "md" }: { name: string; avatar?: string; size?: "sm" | "md" }) {
  const sizeClass = size === "sm" ? "h-8 w-8 text-xs" : "h-11 w-11 text-sm";
  const initial = name.trim().slice(0, 1).toUpperCase() || "?";
  const safeAvatar = avatar && !avatar.includes("gravatar.com/avatar/") ? avatar : "";

  return (
    <div className={`flex shrink-0 items-center justify-center overflow-hidden rounded-full border border-white/70 bg-rose-50/70 font-semibold text-slate-700 shadow-[inset_0_1px_0_rgba(255,255,255,0.75)] ${sizeClass}`}>
      {safeAvatar ? <img src={safeAvatar} alt={name} className="h-full w-full object-cover" loading="lazy" /> : initial}
    </div>
  );
}
