"use client";

import { useDeferredValue, useState } from "react";

import { ProductCard } from "@/components/product-card";

type CatalogProduct = Parameters<typeof ProductCard>[0]["product"];

export function CatalogBrowser({ products }: { products: CatalogProduct[] }) {
  const [query, setQuery] = useState("");
  const [filter, setFilter] = useState<"all" | "live" | "graded">("all");
  const deferredQuery = useDeferredValue(query);

  const normalizedQuery = deferredQuery.trim().toLowerCase();
  const filtered = products.filter((product) => {
    const matchesQuery =
      normalizedQuery.length === 0 ||
      [product.title, product.summary, product.player ?? "", product.team ?? ""]
        .join(" ")
        .toLowerCase()
        .includes(normalizedQuery);
    const matchesFilter =
      filter === "all" ||
      (filter === "live" && product.liveExclusive) ||
      (filter === "graded" && Boolean(product.grade));

    return matchesQuery && matchesFilter;
  });

  return (
    <div>
      <div className="flex flex-col gap-4 rounded-[2rem] border border-white/10 bg-white/5 p-5 backdrop-blur md:flex-row md:items-center md:justify-between">
        <input
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Search player, team, set, or card"
          className="h-12 w-full rounded-full border border-white/10 bg-[#0a1525] px-5 text-sm text-white outline-none placeholder:text-white/35 md:max-w-md"
        />
        <div className="flex flex-wrap gap-2">
          {[
            ["all", "All inventory"],
            ["live", "Live priority"],
            ["graded", "Graded cards"],
          ].map(([value, label]) => {
            const active = filter === value;
            return (
              <button
                key={value}
                type="button"
                onClick={() => setFilter(value as "all" | "live" | "graded")}
                className={`rounded-full px-4 py-2 text-sm font-medium transition ${
                  active
                    ? "bg-[#f7b733] text-slate-950"
                    : "border border-white/10 bg-white/5 text-white/70 hover:bg-white/10"
                }`}
              >
                {label}
              </button>
            );
          })}
        </div>
      </div>
      <div className="mt-8 grid gap-5 lg:grid-cols-3">
        {filtered.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
      {filtered.length === 0 ? (
        <div className="mt-8 rounded-[2rem] border border-dashed border-white/12 bg-white/3 p-10 text-center text-white/60">
          No cards matched the current filter.
        </div>
      ) : null}
    </div>
  );
}
