import type { Metadata } from "next";

import { AuthProvider } from "@/components/auth-provider";
import { CartProvider } from "@/components/cart-provider";
import { SiteFooter } from "@/components/site-footer";
import { SiteHeader } from "@/components/site-header";
import { appUrl } from "@/lib/env";
import { getActiveLiveSession } from "@/lib/store";
import "./globals.css";

export const metadata: Metadata = {
  metadataBase: new URL(appUrl),
  title: "Creswood Corners Cards",
  description:
    "Modern football trading-card storefront with live selling, admin inventory control, and Stripe-ready checkout.",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const liveSession = await getActiveLiveSession();

  return (
    <html lang="en" className="h-full antialiased">
      <body className="min-h-full bg-[#08111f] text-white">
        <AuthProvider>
          <CartProvider>
            <div className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top,rgba(247,183,51,0.12),transparent_25%),linear-gradient(180deg,#08111f_0%,#0d1726_45%,#07101b_100%)]">
              <div className="absolute inset-x-0 top-0 h-[32rem] bg-[radial-gradient(circle_at_top_left,rgba(96,165,250,0.16),transparent_32%),radial-gradient(circle_at_top_right,rgba(247,183,51,0.18),transparent_30%)]" />
              <div className="relative flex min-h-screen flex-col">
                <SiteHeader />
                {liveSession ? (
                  <div className="border-b border-white/8 bg-transparent">
                    <div className="mx-auto w-full max-w-7xl px-5 sm:px-8">
                      <div className="flex flex-wrap items-center gap-3 py-3 text-sm text-white/68">
                        <span className="rounded-full bg-emerald-400/18 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-emerald-200">
                          Live now
                        </span>
                        <span>{liveSession.title}</span>
                        <span className="text-white/38">/</span>
                        <span>{liveSession.platform}</span>
                      </div>
                    </div>
                  </div>
                ) : null}
                <main className="flex-1">{children}</main>
                <SiteFooter />
              </div>
            </div>
          </CartProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
