import Link from "next/link";
import { ShieldCheck, UserRound } from "lucide-react";

import { CartLink } from "@/components/cart-link";
import { getCurrentViewer } from "@/lib/auth";

export async function SiteHeader() {
  const viewer = await getCurrentViewer();

  return (
    <header className="sticky top-0 z-40 border-b border-white/8 bg-[#08111f]/80 backdrop-blur-xl">
      <div className="mx-auto flex w-full max-w-7xl items-center justify-between gap-6 px-5 py-4 sm:px-8">
        <Link href="/" className="flex items-center gap-3">
          <div className="grid h-10 w-10 place-items-center rounded-2xl bg-[#f7b733] font-black text-slate-950">
            CC
          </div>
          <div>
            <p className="font-display text-2xl uppercase leading-none tracking-[0.18em] text-white">
              Creswood
            </p>
            <p className="text-xs uppercase tracking-[0.28em] text-white/45">Corners Cards</p>
          </div>
        </Link>
        <nav className="hidden items-center gap-6 text-sm text-white/72 md:flex">
          <Link href="/shop" className="transition hover:text-white">
            Shop
          </Link>
          <Link href="/live" className="transition hover:text-white">
            Live selling
          </Link>
          <Link href="/shop?view=featured" className="transition hover:text-white">
            Featured
          </Link>
        </nav>
        <div className="flex items-center gap-3">
          <CartLink />
          {viewer.isAdmin ? (
            <Link
              href="/admin"
              className="hidden items-center gap-2 rounded-full border border-[#f7b733]/30 bg-[#f7b733]/12 px-4 py-2 text-sm font-medium text-[#ffd16d] transition hover:bg-[#f7b733]/20 sm:inline-flex"
            >
              <ShieldCheck className="size-4" />
              Admin
            </Link>
          ) : null}
          <Link
            href={viewer.isAuthenticated ? "/user" : "/sign-in"}
            className="inline-flex items-center gap-2 rounded-full border border-white/12 bg-white/7 px-4 py-2 text-sm font-medium text-white transition hover:bg-white/12"
          >
            <UserRound className="size-4" />
            {viewer.isAuthenticated ? "Account" : "Sign in"}
          </Link>
        </div>
      </div>
    </header>
  );
}
