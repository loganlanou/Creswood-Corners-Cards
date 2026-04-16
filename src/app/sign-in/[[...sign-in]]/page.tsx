import Link from "next/link";
import { SignIn } from "@clerk/nextjs";

import { isClerkConfigured } from "@/lib/env";

export default function SignInPage() {
  if (!isClerkConfigured) {
    return (
      <div className="mx-auto max-w-3xl px-5 py-16 sm:px-8">
        <div className="rounded-[2rem] border border-white/10 bg-white/5 p-8">
          <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
            Clerk setup required
          </p>
          <h1 className="mt-4 text-4xl font-semibold tracking-tight text-white">
            Add Clerk keys to enable production auth.
          </h1>
          <p className="mt-4 text-base leading-8 text-white/62">
            The app is currently running in local demo mode so the storefront and admin can be
            built and tested before live credentials are attached.
          </p>
          <Link
            href="/admin"
            className="mt-6 inline-flex rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950"
          >
            Open demo admin
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto flex min-h-[70vh] max-w-7xl items-center justify-center px-5 py-12 sm:px-8">
      <SignIn />
    </div>
  );
}
