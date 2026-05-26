export function getOsIcon(os?: string) {
  const value = os?.toLowerCase() || "";
  if (value.includes("mac")) return "🍎";
  if (value.includes("win")) return "🪟";
  if (value.includes("linux")) return "🐧";
  if (value.includes("android")) return "🤖";
  if (value.includes("ios")) return "📱";
  return "💻";
}

export function getBrowserIcon(browser?: string) {
  const value = browser?.toLowerCase() || "";
  if (value.includes("safari")) return "🧭";
  if (value.includes("firefox")) return "🦊";
  if (value.includes("edge")) return "🌊";
  if (value.includes("chrome")) return "🌐";
  return "🌐";
}

export function parseBrowser(userAgent?: string) {
  const value = userAgent || "";
  if (/Edg\//.test(value)) return "Edge";
  if (/Firefox\//.test(value)) return "Firefox";
  if (/Chrome\//.test(value) || /CriOS\//.test(value)) return "Chrome";
  if (/Safari\//.test(value)) return "Safari";
  return "";
}

export function parseOS(userAgent?: string) {
  const value = userAgent || "";
  if (/Android/i.test(value)) return "Android";
  if (/iPhone|iPad|iPod/i.test(value)) return "iOS";
  if (/Mac OS X|Macintosh/i.test(value)) return "macOS";
  if (/Windows NT/i.test(value)) return "Windows";
  if (/Linux/i.test(value)) return "Linux";
  return "";
}

export function locationFromIP(ip?: string) {
  const value = ip?.trim();
  if (!value) return "未知地区";
  if (value === "127.0.0.1" || value === "::1" || value === "localhost") return "本地";
  if (/^10\./.test(value) || /^192\.168\./.test(value) || /^172\.(1[6-9]|2\d|3[0-1])\./.test(value)) return "内网";
  return "未知地区";
}

export function publicMetaParts(item: {
  location?: string;
  ip_location?: string;
  ipLocation?: string;
  ip?: string;
  browser?: string;
  os?: string;
  user_agent?: string;
  userAgent?: string;
}) {
  const userAgent = item.userAgent || item.user_agent;
  const location = item.location || item.ipLocation || item.ip_location || locationFromIP(item.ip);
  const browser = item.browser || parseBrowser(userAgent);
  const os = item.os || parseOS(userAgent);
  return [location || "未知地区", browser, os].filter(Boolean);
}
