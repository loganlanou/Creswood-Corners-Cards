import { cn } from "@/lib/utils";

const accentMap: Record<string, string> = {
  amber:
    "from-amber-300/90 via-amber-500/65 to-orange-700/80 shadow-[0_25px_80px_-30px_rgba(245,158,11,0.65)]",
  crimson:
    "from-rose-300/85 via-red-500/60 to-red-900/85 shadow-[0_25px_80px_-30px_rgba(239,68,68,0.7)]",
  emerald:
    "from-emerald-300/90 via-emerald-500/60 to-emerald-900/80 shadow-[0_25px_80px_-30px_rgba(16,185,129,0.65)]",
  blue: "from-sky-300/90 via-blue-500/60 to-indigo-900/80 shadow-[0_25px_80px_-30px_rgba(59,130,246,0.7)]",
  violet:
    "from-fuchsia-300/85 via-violet-500/60 to-violet-900/80 shadow-[0_25px_80px_-30px_rgba(139,92,246,0.7)]",
  slate:
    "from-slate-300/85 via-slate-500/60 to-slate-900/80 shadow-[0_25px_80px_-30px_rgba(100,116,139,0.7)]",
};

export function ProductVisual({
  title,
  player,
  team,
  accent,
  compact = false,
}: {
  title: string;
  player?: string | null;
  team?: string | null;
  accent?: string | null;
  compact?: boolean;
}) {
  const accentClasses = accent ? accentMap[accent] : accentMap.amber;

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-[1.75rem] border border-white/12 bg-slate-950 text-white",
        compact ? "aspect-[4/5]" : "aspect-[5/6]",
      )}
    >
      <div className={cn("absolute inset-0 bg-gradient-to-br", accentClasses)} />
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(255,255,255,0.38),_transparent_28%),linear-gradient(160deg,rgba(255,255,255,0.16),transparent_45%)]" />
      <div className="absolute inset-x-4 top-4 flex items-start justify-between text-[0.65rem] uppercase tracking-[0.28em] text-white/70">
        <span>Creswood</span>
        <span>Verified</span>
      </div>
      <div className="absolute inset-x-5 bottom-5 rounded-[1.4rem] border border-white/15 bg-slate-950/40 p-4 backdrop-blur">
        <p className="text-[0.7rem] uppercase tracking-[0.28em] text-white/60">
          {team ?? "Football"}
        </p>
        <h3 className={cn("mt-2 font-semibold tracking-tight", compact ? "text-lg" : "text-2xl")}>
          {player ?? title}
        </h3>
        <p className="mt-2 line-clamp-2 text-sm text-white/72">{title}</p>
      </div>
    </div>
  );
}
