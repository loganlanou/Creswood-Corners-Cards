import Link from "next/link";

import { ProductCard } from "@/components/product-card";
import { getActiveLiveSession, getCatalogProducts } from "@/lib/store";

export default async function LivePage() {
  const [liveSession, products] = await Promise.all([getActiveLiveSession(), getCatalogProducts()]);
  const livePriority = products.filter((product) => product.liveExclusive).slice(0, 3);

  return (
    <div className="mx-auto w-full max-w-7xl px-5 py-12 sm:px-8">
      <div className="grid gap-6 rounded-[2.4rem] border border-white/10 bg-white/5 p-6 md:grid-cols-[1fr_auto] md:items-center">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
            Live selling
          </p>
          <h1 className="mt-4 text-5xl font-semibold tracking-tight text-white">
            Stream traffic gets a dedicated home instead of a buried link.
          </h1>
          <p className="mt-4 max-w-3xl text-base leading-8 text-white/62">
            When a show is active, the homepage banner, this page, and live-priority inventory all
            reinforce the same destination so customers can move between the stream and the store
            without confusion.
          </p>
        </div>
        {liveSession ? (
          <Link
            href={liveSession.streamUrl}
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center justify-center rounded-full bg-[#f7b733] px-6 py-3 text-sm font-semibold text-slate-950"
          >
            Join {liveSession.platform}
          </Link>
        ) : null}
      </div>

      <div className="mt-8 rounded-[2rem] border border-white/10 bg-[#0b1422] p-6">
        {liveSession ? (
          <>
            <p className="text-sm font-semibold uppercase tracking-[0.24em] text-emerald-200">
              Live now
            </p>
            <p className="mt-3 text-3xl font-semibold text-white">{liveSession.title}</p>
            <p className="mt-3 max-w-3xl text-base leading-8 text-white/60">{liveSession.pitch}</p>
            {liveSession.callout ? (
              <p className="mt-2 text-sm text-white/50">{liveSession.callout}</p>
            ) : null}
          </>
        ) : (
          <>
            <p className="text-sm font-semibold uppercase tracking-[0.24em] text-white/55">
              No active stream
            </p>
            <p className="mt-3 text-3xl font-semibold text-white">
              Admin can turn live selling on as soon as the next show is scheduled.
            </p>
          </>
        )}
      </div>

      <section className="mt-10">
        <div className="flex items-end justify-between gap-4">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
              Stream-ready cards
            </p>
            <h2 className="mt-3 text-4xl font-semibold tracking-tight text-white">
              Inventory tagged for live mentions and quick claims.
            </h2>
          </div>
          <Link href="/shop" className="text-sm font-medium text-white/70 hover:text-white">
            Browse all products
          </Link>
        </div>
        <div className="mt-8 grid gap-5 lg:grid-cols-3">
          {livePriority.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      </section>
    </div>
  );
}
