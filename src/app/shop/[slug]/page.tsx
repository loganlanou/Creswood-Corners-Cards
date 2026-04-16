import Link from "next/link";
import { notFound } from "next/navigation";
import { CheckCircle2, ShieldCheck, Truck } from "lucide-react";

import { AddToCartButton } from "@/components/add-to-cart-button";
import { ProductVisual } from "@/components/product-visual";
import { formatCurrency } from "@/lib/utils";
import { getProductBySlug } from "@/lib/store";

export default async function ProductDetailPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const product = await getProductBySlug(slug);

  if (!product) {
    notFound();
  }

  return (
    <div className="mx-auto grid w-full max-w-7xl gap-10 px-5 py-12 sm:px-8 lg:grid-cols-[0.95fr_1.05fr]">
      <div className="card-sheen">
        <ProductVisual
          title={product.title}
          player={product.player}
          team={product.team}
          accent={product.accent}
        />
      </div>
      <div>
        <div className="flex flex-wrap gap-2">
          <span className="rounded-full bg-[#f7b733]/15 px-3 py-1 text-[0.72rem] font-semibold uppercase tracking-[0.2em] text-[#ffd16d]">
            {product.sport}
          </span>
          {product.grade ? (
            <span className="rounded-full border border-white/12 px-3 py-1 text-[0.72rem] uppercase tracking-[0.2em] text-white/60">
              {product.grade}
            </span>
          ) : null}
          {product.liveExclusive ? (
            <span className="rounded-full border border-emerald-300/20 bg-emerald-300/10 px-3 py-1 text-[0.72rem] uppercase tracking-[0.2em] text-emerald-200">
              Live stream priority
            </span>
          ) : null}
        </div>
        <h1 className="mt-5 text-5xl font-semibold tracking-tight text-white">{product.title}</h1>
        <p className="mt-4 max-w-2xl text-lg leading-8 text-white/64">{product.description}</p>
        <p className="mt-6 text-4xl font-semibold text-[#f7b733]">
          {formatCurrency(product.priceCents)}
        </p>
        <div className="mt-8 flex flex-wrap gap-3">
          <AddToCartButton
            className="px-6"
            item={{
              productId: product.id,
              slug: product.slug,
              title: product.title,
              summary: product.summary,
              priceCents: product.priceCents,
              accent: product.accent,
            }}
          />
          <Link
            href="/cart"
            className="inline-flex items-center justify-center rounded-full border border-white/12 bg-white/5 px-6 py-3 text-sm font-semibold text-white"
          >
            Go to cart
          </Link>
        </div>
        <div className="mt-10 grid gap-3 rounded-[2rem] border border-white/10 bg-white/5 p-6 text-sm text-white/68">
          {[
            {
              icon: CheckCircle2,
              label: `Condition: ${product.condition ?? "Collector-reviewed"}`,
            },
            {
              icon: ShieldCheck,
              label: "Secure checkout with Stripe when configured",
            },
            {
              icon: Truck,
              label: "Admin-managed shipment status and fulfillment tracking",
            },
          ].map((item) => (
            <div key={item.label} className="flex items-center gap-3">
              <item.icon className="size-4 text-[#f7b733]" />
              <span>{item.label}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
