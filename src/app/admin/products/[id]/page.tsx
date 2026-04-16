import { ProductStatus } from "@prisma/client";
import Link from "next/link";
import { notFound, redirect } from "next/navigation";

import { saveProductAction } from "@/app/admin/actions";
import { getCurrentViewer } from "@/lib/auth";
import { prisma } from "@/lib/prisma";

export default async function AdminProductEditPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const viewer = await getCurrentViewer();
  if (!viewer.isAdmin) {
    redirect("/sign-in");
  }

  const { id } = await params;
  const product = await prisma.product.findUnique({ where: { id } });

  if (!product) {
    notFound();
  }

  return (
    <div className="mx-auto max-w-5xl px-5 py-12 sm:px-8">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
            Edit product
          </p>
          <h1 className="mt-3 text-4xl font-semibold tracking-tight text-white">{product.title}</h1>
        </div>
        <Link
          href="/admin"
          className="rounded-full border border-white/12 bg-white/5 px-5 py-3 text-sm font-semibold text-white"
        >
          Back to admin
        </Link>
      </div>

      <form action={saveProductAction} className="mt-8 grid gap-4 rounded-[2rem] border border-white/10 bg-white/5 p-6">
        <input type="hidden" name="id" value={product.id} />
        <div className="grid gap-4 md:grid-cols-2">
          <input name="title" defaultValue={product.title} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="slug" defaultValue={product.slug} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <input name="player" defaultValue={product.player ?? ""} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="team" defaultValue={product.team ?? ""} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
        </div>
        <div className="grid gap-4 md:grid-cols-3">
          <input name="brand" defaultValue={product.brand ?? ""} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="setName" defaultValue={product.setName ?? ""} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="cardNumber" defaultValue={product.cardNumber ?? ""} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
        </div>
        <div className="grid gap-4 md:grid-cols-4">
          <input name="year" type="number" defaultValue={product.year ?? undefined} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="grade" defaultValue={product.grade ?? ""} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="condition" defaultValue={product.condition ?? ""} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="accent" defaultValue={product.accent ?? ""} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
        </div>
        <div className="grid gap-4 md:grid-cols-4">
          <input name="sport" defaultValue={product.sport} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="priceCents" type="number" defaultValue={product.priceCents} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <input name="quantity" type="number" defaultValue={product.quantity} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
          <select name="status" defaultValue={product.status} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm">
            {Object.values(ProductStatus).map((status) => (
              <option key={status} value={status}>
                {status}
              </option>
            ))}
          </select>
        </div>
        <textarea name="summary" defaultValue={product.summary} className="min-h-24 rounded-2xl border border-white/10 bg-[#0a1525] px-4 py-3 text-sm" />
        <textarea name="description" defaultValue={product.description} className="min-h-36 rounded-2xl border border-white/10 bg-[#0a1525] px-4 py-3 text-sm" />
        <div className="flex flex-wrap gap-6 text-sm text-white/72">
          <label className="flex items-center gap-2">
            <input type="checkbox" name="featured" defaultChecked={product.featured} />
            Featured
          </label>
          <label className="flex items-center gap-2">
            <input type="checkbox" name="liveExclusive" defaultChecked={product.liveExclusive} />
            Live priority
          </label>
        </div>
        <button className="inline-flex w-fit rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950">
          Save product changes
        </button>
      </form>
    </div>
  );
}
