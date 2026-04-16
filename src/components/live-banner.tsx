import Link from "next/link";
import { RadioTower } from "lucide-react";

export function LiveBanner({
  session,
}: {
  session:
    | {
        title: string;
        pitch: string;
        callout: string | null;
        streamUrl: string;
        platform: string;
      }
    | null;
}) {
  if (!session) {
    return null;
  }

  return (
    <div className="border-b border-[#f7b733]/22 bg-[linear-gradient(90deg,rgba(247,183,51,0.16),rgba(247,183,51,0.05),rgba(13,22,37,0.1))]">
      <div className="mx-auto flex w-full max-w-7xl flex-col gap-4 px-5 py-4 sm:px-8 md:flex-row md:items-center md:justify-between">
        <div className="flex items-start gap-3">
          <div className="mt-1 rounded-full bg-[#f7b733] p-2 text-slate-950">
            <RadioTower className="size-4" />
          </div>
          <div>
            <p className="text-[0.72rem] font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
              Live now on {session.platform}
            </p>
            <p className="mt-1 text-base font-semibold text-white">{session.title}</p>
            <p className="text-sm text-white/65">{session.callout ?? session.pitch}</p>
          </div>
        </div>
        <Link
          href={session.streamUrl}
          target="_blank"
          rel="noreferrer"
          className="inline-flex items-center justify-center rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-[#ffc95a]"
        >
          Join live selling
        </Link>
      </div>
    </div>
  );
}
