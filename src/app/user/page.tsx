import Link from "next/link";

import { getCurrentViewer } from "@/lib/auth";
import { isClerkConfigured } from "@/lib/env";

export default async function UserPage() {
  const viewer = await getCurrentViewer();

  return (
    <div className="mx-auto max-w-4xl px-5 py-16 sm:px-8">
      <div className="rounded-[2rem] border border-white/10 bg-white/5 p-8">
        <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
          Account
        </p>
        <h1 className="mt-4 text-4xl font-semibold tracking-tight text-white">
          {viewer.isAuthenticated ? "Signed in account" : "Guest mode"}
        </h1>
        <p className="mt-4 text-base leading-8 text-white/62">
          {viewer.isAuthenticated
            ? `Signed in as ${viewer.email}. Admin access is ${viewer.isAdmin ? "enabled" : "not enabled"}.`
            : isClerkConfigured
              ? "Use Clerk sign-in to manage your account and unlock admin access when granted."
              : "Clerk keys have not been configured yet, so the local build is operating in demo mode."}
        </p>
        {!viewer.isAuthenticated && isClerkConfigured ? (
          <Link
            href="/sign-in"
            className="mt-6 inline-flex rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950"
          >
            Sign in
          </Link>
        ) : null}
      </div>
    </div>
  );
}
