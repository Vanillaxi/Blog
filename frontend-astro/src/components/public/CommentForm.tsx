import { useRef, useState } from "react";
import { ChevronDown, MessageCircle, RefreshCw, Smile } from "lucide-react";
import { addComment } from "@/api/comment";
import { getVisitorIdentity } from "@/api/visitor";

interface Props {
  mode: "comment" | "guestbook";
  articleId?: number;
  targetType?: 1 | 2;
  targetId?: number;
  replyTo?: {
    id: number | string;
    nickname: string;
  } | null;
  onCancelReply?: () => void;
  onSubmitted?: () => void | Promise<void>;
}

type Errors = Partial<Record<"nickname" | "email" | "content" | "captchaAnswer", string>>;

type Captcha = {
  id: string;
  question: string;
  answer: string;
};

const EMOJI_OPTIONS = ["😊", "😂", "😭", "😳", "😍", "😎", "🤔", "👍", "🎉", "🌸", "🍵", "✨", "💻", "🐱", "🫶"];
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function generateMockCaptcha(): Captcha {
  const operators = ["+", "-", "*"] as const;
  const operator = operators[Math.floor(Math.random() * operators.length)];
  let left = Math.floor(Math.random() * 9) + 1;
  let right = Math.floor(Math.random() * 9) + 1;

  if (operator === "-" && right > left) {
    [left, right] = [right, left];
  }

  const answer = operator === "+" ? left + right : operator === "-" ? left - right : left * right;

  return {
    id: `mock-captcha-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    question: `${left} ${operator} ${right} =`,
    answer: String(answer),
  };
}

export function CommentForm({ mode, articleId, targetType, targetId, replyTo, onCancelReply, onSubmitted }: Props) {
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState("");
  const [errors, setErrors] = useState<Errors>({});
  const [submitError, setSubmitError] = useState("");
  const [identityMessage, setIdentityMessage] = useState("");
  const [identityLoading, setIdentityLoading] = useState(false);
  const [markdownEnabled, setMarkdownEnabled] = useState(false);
  const [emojiOpen, setEmojiOpen] = useState(false);
  const [showWebsite, setShowWebsite] = useState(false);
  const [captcha, setCaptcha] = useState<Captcha>(() => generateMockCaptcha());
  const [nickname, setNickname] = useState("");
  const [boundNickname, setBoundNickname] = useState("");
  const [email, setEmail] = useState("");
  const [content, setContent] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const isComment = mode === "comment";
  const resolvedTargetType = targetType ?? (isComment ? 1 : 2);
  const resolvedTargetId = targetId ?? articleId ?? 0;
  const replyNickname = replyTo?.nickname;

  function validate(form: HTMLFormElement) {
    const data = new FormData(form);
    const next: Errors = {};
    const emailValue = String(data.get("email") || "").trim();
    const nicknameValue = boundNickname || String(data.get("nickname") || "").trim();
    const captchaAnswer = String(data.get("captchaAnswer") || "").trim();
    if (!content.trim()) next.content = isComment ? "请写下你的想法" : "请写一点内容";
    if (!nicknameValue) next.nickname = "请填写昵称";
    if (!emailValue) next.email = "请填写邮箱";
    else if (!EMAIL_PATTERN.test(emailValue)) next.email = "邮箱格式不正确";
    if (!captchaAnswer) next.captchaAnswer = "请填写验证码";
    else if (captchaAnswer !== captcha.answer) next.captchaAnswer = "验证码不正确，请再试一次。";
    return next;
  }

  async function lookupVisitorIdentity(nextEmail = email) {
    const value = nextEmail.trim();
    setIdentityMessage("");

    if (!value || !EMAIL_PATTERN.test(value)) {
      setBoundNickname("");
      return;
    }

    setIdentityLoading(true);
    try {
      const response = await getVisitorIdentity(value);
      if (response.data.exists && response.data.nickname) {
        setNickname(response.data.nickname);
        setBoundNickname(response.data.nickname);
        setErrors((current) => ({ ...current, nickname: undefined, email: undefined }));
        setIdentityMessage(`该邮箱已绑定昵称：${response.data.nickname}，已自动填入。`);
      } else {
        setBoundNickname("");
      }
    } catch (err) {
      console.error("[visitor-identity] lookup failed", err);
    } finally {
      setIdentityLoading(false);
    }
  }

  function refreshCaptcha() {
    setCaptcha(generateMockCaptcha());
    setErrors((current) => ({ ...current, captchaAnswer: undefined }));
  }

  function insertEmoji(emoji: string) {
    const textarea = textareaRef.current;
    const start = textarea?.selectionStart ?? content.length;
    const end = textarea?.selectionEnd ?? content.length;
    const nextContent = `${content.slice(0, start)}${emoji}${content.slice(end)}`;

    setContent(nextContent);
    setErrors((current) => ({ ...current, content: undefined }));

    requestAnimationFrame(() => {
      textarea?.focus();
      const cursor = start + emoji.length;
      textarea?.setSelectionRange(cursor, cursor);
    });
  }

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const nextErrors = validate(form);
    setErrors(nextErrors);
    setSuccess("");
    setSubmitError("");
    if (Object.keys(nextErrors).length > 0) return;
    setLoading(true);
    const data = new FormData(form);
    try {
      await addComment({
        target_type: resolvedTargetType,
        target_id: resolvedTargetId,
        parent_id: Number(replyTo?.id ?? 0),
        nickname: boundNickname || String(data.get("nickname") || "").trim(),
        email: String(data.get("email") || "").trim(),
        website: String(data.get("website") || "") || undefined,
        content,
      });
      setSuccess(isComment ? "评论已提交。" : "留言已提交。");
      form.reset();
      setNickname("");
      setBoundNickname("");
      setEmail("");
      setContent("");
      setIdentityMessage("");
      setShowWebsite(false);
      setMarkdownEnabled(false);
      setEmojiOpen(false);
      setCaptcha(generateMockCaptcha());
      onCancelReply?.();
      await onSubmitted?.();
    } catch (err) {
      console.error("[comment] submit failed", err);
      setSubmitError(err instanceof Error ? err.message : "提交失败，请稍后再试。");
      setCaptcha(generateMockCaptcha());
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="glass-card relative z-10 rounded-[1.75rem] p-6">
      <div className="mb-5 flex items-center gap-2">
        <MessageCircle className="h-5 w-5 text-rose-400" />
        <h2 className="text-lg font-semibold text-[#24314f]">{isComment ? "写评论" : "写留言"}</h2>
      </div>
      {replyNickname && (
        <div className="mb-4 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-white/50 bg-white/45 px-4 py-3 text-sm text-slate-600">
          <span>
            正在回复 <span className="font-medium text-[#24314f]">@{replyNickname}</span>
          </span>
          <button type="button" onClick={onCancelReply} className="rounded-full bg-white/60 px-3 py-1 text-xs text-slate-500 transition hover:bg-white hover:text-slate-700">
            取消回复
          </button>
        </div>
      )}
      <div>
        <label className="text-sm font-medium text-slate-700">内容 *</label>
        <textarea
          ref={textareaRef}
          name="content"
          value={content}
          onChange={(event) => setContent(event.target.value)}
          placeholder={replyNickname ? `回复 @${replyNickname}...` : isComment ? "写下你的想法..." : "写点什么吧..."}
          className="mt-2 min-h-32 w-full rounded-2xl border border-white/55 bg-white/55 px-4 py-3 text-sm outline-none focus:border-rose-200"
        />
        {errors.content && <p className="mt-1 text-xs text-rose-500">{errors.content}</p>}
      </div>
      <div className="mt-4 grid gap-4 sm:grid-cols-2">
        <Field
          name="nickname"
          label="昵称"
          placeholder="你的名字"
          error={errors.nickname}
          required
          value={nickname}
          onChange={(value) => {
            if (boundNickname) return;
            setNickname(value);
            setErrors((current) => ({ ...current, nickname: undefined }));
          }}
          readOnly={!!boundNickname}
        />
        <Field
          name="email"
          label="邮箱"
          placeholder="your@email.com"
          error={errors.email}
          required
          type="email"
          value={email}
          onChange={(value) => {
            setEmail(value);
            setIdentityMessage("");
            setBoundNickname("");
            setErrors((current) => ({ ...current, email: undefined }));
          }}
          onBlur={() => void lookupVisitorIdentity()}
        />
      </div>
      {identityLoading && <p className="mt-2 text-xs text-slate-400">正在查询邮箱绑定昵称...</p>}
      {identityMessage && <p className="mt-2 rounded-2xl bg-sky-50/80 px-4 py-2 text-xs text-sky-700">{identityMessage}</p>}
      <div className="mt-4">
        <button
          type="button"
          onClick={() => setShowWebsite((current) => !current)}
          className="inline-flex items-center gap-1 rounded-full border border-white/45 bg-white/35 px-3 py-1.5 text-xs text-slate-500 transition hover:bg-white/50 hover:text-slate-700"
          aria-expanded={showWebsite}
        >
          {showWebsite ? "收起个人网站" : "+ 添加个人网站（可选）"}
          <ChevronDown className={`h-3.5 w-3.5 transition ${showWebsite ? "rotate-180" : ""}`} />
        </button>
        {showWebsite && (
          <div className="mt-3">
            <Field name="website" label="个人网站" placeholder="https://your-site.com" />
          </div>
        )}
      </div>
      <div className="mt-4 grid gap-4 sm:grid-cols-[minmax(0,1fr)_auto_auto_auto] sm:items-start">
        <label className="block text-sm font-medium text-slate-700">
          <span className="sr-only">验证码</span>
          <div className="flex flex-wrap items-center gap-2">
            <input name="captchaAnswer" placeholder="输入答案" className="min-w-0 flex-1 rounded-2xl border border-white/55 bg-white/55 px-4 py-3 text-sm outline-none focus:border-rose-200" />
            <span className="rounded-full border border-white/55 bg-white/50 px-4 py-2.5 text-sm text-slate-600">{captcha.question}</span>
            <button
              type="button"
              onClick={refreshCaptcha}
              className="inline-flex items-center gap-1 rounded-full border border-white/55 bg-white/40 px-3 py-2 text-xs text-slate-500 transition hover:bg-white/60 hover:text-slate-700"
            >
              <RefreshCw className="h-3.5 w-3.5" />
              换一题
            </button>
          </div>
          {errors.captchaAnswer && <p className="mt-1 text-xs text-rose-500">{errors.captchaAnswer}</p>}
        </label>
        <button
          type="button"
          onClick={() => setMarkdownEnabled((current) => !current)}
          className={`rounded-2xl border border-white/55 px-4 py-3 text-sm ${
            markdownEnabled ? "bg-rose-100/60 text-rose-900" : "bg-white/50 text-slate-600"
          }`}
        >
          Markdown
        </button>
        <div className="relative">
          <button
            type="button"
            onClick={() => setEmojiOpen((current) => !current)}
            className={`inline-flex w-full items-center justify-center gap-2 rounded-2xl border border-white/55 px-4 py-3 text-sm transition ${
              emojiOpen ? "bg-rose-100/60 text-rose-900" : "bg-white/50 text-slate-600 hover:bg-white/65"
            }`}
            aria-expanded={emojiOpen}
          >
            <Smile className="h-4 w-4" /> 表情
          </button>
          {emojiOpen && (
            <div className="absolute bottom-full right-0 z-[999] mb-2 grid w-72 max-w-[calc(100vw-32px)] grid-cols-5 gap-2 rounded-2xl border border-white/65 bg-white/85 p-3 shadow-[0_18px_45px_rgba(15,23,42,0.14)] backdrop-blur-xl">
              {EMOJI_OPTIONS.map((emoji) => (
                <button
                  key={emoji}
                  type="button"
                  onClick={() => insertEmoji(emoji)}
                  className="flex h-9 w-9 items-center justify-center rounded-xl bg-white/55 text-lg transition hover:bg-rose-50"
                  aria-label={`插入 ${emoji}`}
                >
                  {emoji}
                </button>
              ))}
            </div>
          )}
        </div>
        <button disabled={loading} className="rounded-full bg-slate-900 px-5 py-2.5 text-sm text-white disabled:cursor-not-allowed disabled:opacity-50">
          {loading ? "发送中..." : isComment ? "发表评论" : "发送留言"}
        </button>
      </div>
      <p className="mt-3 text-xs text-slate-400">邮箱不会完整展示，仅用于固定你的昵称。</p>
      {success && <p className="mt-4 rounded-2xl bg-emerald-50 px-4 py-3 text-sm text-emerald-700">{success}</p>}
      {submitError && <p className="mt-4 rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">{submitError}</p>}
    </form>
  );
}

function Field({
  name,
  label,
  placeholder,
  error,
  required,
  type = "text",
  value,
  onChange,
  onBlur,
  readOnly,
}: {
  name: string;
  label: string;
  placeholder: string;
  error?: string;
  required?: boolean;
  type?: string;
  value?: string;
  onChange?: (value: string) => void;
  onBlur?: () => void;
  readOnly?: boolean;
}) {
  return (
    <label className="block text-sm font-medium text-slate-700">
      {label}{required ? " *" : ""}
      <input
        name={name}
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={(event) => onChange?.(event.target.value)}
        onBlur={onBlur}
        readOnly={readOnly}
        aria-readonly={readOnly}
        className={`mt-2 w-full rounded-2xl border border-white/55 px-4 py-3 text-sm outline-none focus:border-rose-200 ${
          readOnly ? "cursor-not-allowed bg-white/35 text-slate-500" : "bg-white/55"
        }`}
      />
      {error && <p className="mt-1 text-xs text-rose-500">{error}</p>}
    </label>
  );
}
