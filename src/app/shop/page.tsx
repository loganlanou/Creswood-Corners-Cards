import { CatalogBrowser } from "@/components/catalog-browser";
import { getCatalogProducts } from "@/lib/store";

export default async function ShopPage() {
  const products = await getCatalogProducts();

  return (
    <div className="mx-auto w-full max-w-7xl px-5 py-12 sm:px-8">
      <div className="max-w-3xl">
        <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
          Shop the catalog
        </p>
        <h1 className="mt-4 text-5xl font-semibold tracking-tight text-white">
          Fast discovery for football singles, slabs, and live-priority cards.
        </h1>
        <p className="mt-4 text-base leading-8 text-white/62">
          The catalog keeps search friction low and surfaces the same inventory that can be
          spotlighted during live streams and direct offers.
        </p>
      </div>
      <div className="mt-10">
        <CatalogBrowser products={products} />
      </div>
    </div>
  );
}
