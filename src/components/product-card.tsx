import Link from "next/link";

import { AddToCartButton } from "@/components/add-to-cart-button";
import { ProductVisual } from "@/components/product-visual";
import { formatCurrency } from "@/lib/utils";

type ProductCardProps = {
  product: {
    id: string;
    slug: string;
    title: string;
    summary: string;
    priceCents: number;
    player: string | null;
    team: string | null;
    grade: string | null;
    accent: string | null;
    liveExclusive: boolean;
  };
};

export function ProductCard({ product }: ProductCardProps) {
  return (
    <article className="group rounded-[2rem] border border-white/10 bg-white/5 p-4 backdrop-blur transition hover:-translate-y-1 hover:bg-white/7">
      <Link href={`/shop/${product.slug}`} className="block">
        <ProductVisual
          title={product.title}
          player={product.player}
          team={product.team}
          accent={product.accent}
          compact
        />
      </Link>
      <div className="mt-4 flex items-start justify-between gap-4">
        <div>
          <div className="flex flex-wrap gap-2">
            {product.liveExclusive ? (
              <span className="rounded-full bg-[#f7b733]/15 px-3 py-1 text-[0.7rem] font-semibold uppercase tracking-[0.2em] text-[#ffd16d]">
                Live priority
              </span>
            ) : null}
            {product.grade ? (
              <span className="rounded-full border border-white/12 px-3 py-1 text-[0.7rem] uppercase tracking-[0.2em] text-white/65">
                {product.grade}
              </span>
            ) : null}
          </div>
          <Link href={`/shop/${product.slug}`} className="mt-3 block text-xl font-semibold tracking-tight">
            {product.title}
          </Link>
          <p className="mt-2 max-w-xs text-sm leading-6 text-white/65">{product.summary}</p>
        </div>
        <p className="text-lg font-semibold text-[#f7b733]">{formatCurrency(product.priceCents)}</p>
      </div>
      <div className="mt-5 flex items-center justify-between gap-3">
        <Link
          href={`/shop/${product.slug}`}
          className="text-sm font-medium text-white/75 transition hover:text-white"
        >
          View details
        </Link>
        <AddToCartButton
          item={{
            productId: product.id,
            slug: product.slug,
            title: product.title,
            summary: product.summary,
            priceCents: product.priceCents,
            accent: product.accent,
          }}
        />
      </div>
    </article>
  );
}
