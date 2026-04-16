import Link from "next/link";
import { ArrowRight, PackageCheck, ShieldCheck, Zap } from "lucide-react";

import { LiveBanner } from "@/components/live-banner";
import { ProductCard } from "@/components/product-card";
import { getActiveLiveSession, getCatalogProducts, getFeaturedProducts } from "@/lib/store";

export default async function Home() {
  const [liveSession, featuredProducts, catalogProducts] = await Promise.all([
    getActiveLiveSession(),
    getFeaturedProducts(),
    getCatalogProducts(),
  ]);

  return (
    <div>
      <LiveBanner session={liveSession} />
      <section className="mx-auto w-full max-w-7xl px-5 pb-10 pt-12 sm:px-8 lg:pb-16 lg:pt-18">
        <div className="grid gap-10 lg:grid-cols-[1.1fr_0.9fr] lg:items-end">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.32em] text-[#ffd16d]">
              Football singles. Live energy. Cleaner checkout.
            </p>
            <h1 className="mt-5 max-w-3xl font-display text-6xl uppercase leading-[0.92] tracking-[0.06em] text-white sm:text-7xl">
              A modern card shop built for the stream and the storefront.
            </h1>
            <p className="mt-6 max-w-2xl text-lg leading-8 text-white/66">
              Creswood Corners Cards combines direct product sales, live-stream traffic, account
              capture, and operational visibility in one football-focused storefront.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Link
                href="/shop"
                className="inline-flex items-center gap-2 rounded-full bg-[#f7b733] px-6 py-3 text-sm font-semibold text-slate-950"
              >
                Shop the catalog
                <ArrowRight className="size-4" />
              </Link>
              <Link
                href="/live"
                className="inline-flex items-center gap-2 rounded-full border border-white/12 bg-white/5 px-6 py-3 text-sm font-semibold text-white"
              >
                View live selling
              </Link>
            </div>
          </div>
          <div className="relative rounded-[2.25rem] border border-white/10 bg-white/5 p-5 backdrop-blur-xl">
            <div className="grid-fade absolute inset-0 rounded-[2.25rem]" />
            <div className="relative grid gap-4">
              <div className="grid gap-4 md:grid-cols-2">
                {featuredProducts.slice(0, 2).map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
              <div className="rounded-[2rem] border border-white/10 bg-[#09111d]/90 p-5">
                <p className="text-sm font-semibold uppercase tracking-[0.24em] text-white/55">
                  Selling pattern
                </p>
                <p className="mt-3 text-2xl font-semibold tracking-tight text-white">
                  Highlight the live show without sacrificing traditional product discovery.
                </p>
                <p className="mt-3 text-sm leading-7 text-white/60">
                  The best card sites keep stream urgency visible, but still make it easy to search,
                  compare, and buy singles in a normal storefront flow.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="mx-auto grid w-full max-w-7xl gap-4 px-5 py-6 sm:px-8 md:grid-cols-3">
        {[
          {
            icon: Zap,
            title: "Live-first merchandising",
            copy: "Live banners, stream links, and featured inventory stay visible whenever the show is active.",
          },
          {
            icon: ShieldCheck,
            title: "Trust-centric product detail",
            copy: "Condition, grade, shipping, and checkout clarity are front-loaded to reduce hesitation.",
          },
          {
            icon: PackageCheck,
            title: "Operational visibility",
            copy: "Admins can manage inventory, accounts, payment states, and shipment progress in one place.",
          },
        ].map((item) => (
          <article key={item.title} className="rounded-[2rem] border border-white/10 bg-white/5 p-6">
            <item.icon className="size-5 text-[#f7b733]" />
            <p className="mt-4 text-xl font-semibold text-white">{item.title}</p>
            <p className="mt-2 text-sm leading-7 text-white/60">{item.copy}</p>
          </article>
        ))}
      </section>

      <section className="mx-auto w-full max-w-7xl px-5 py-12 sm:px-8">
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
              Featured inventory
            </p>
            <h2 className="mt-3 text-4xl font-semibold tracking-tight text-white">
              Built around the cards that actually drive traffic.
            </h2>
          </div>
          <Link href="/shop" className="text-sm font-medium text-white/70 hover:text-white">
            View full catalog
          </Link>
        </div>
        <div className="mt-8 grid gap-5 lg:grid-cols-4">
          {catalogProducts.slice(0, 4).map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      </section>
    </div>
  );
}
