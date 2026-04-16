import Link from "next/link";

import { prisma } from "@/lib/prisma";
import { formatCurrency } from "@/lib/utils";

export default async function CheckoutSuccessPage({
  searchParams,
}: {
  searchParams: Promise<{ order?: string; demo?: string }>;
}) {
  const params = await searchParams;
  const order = params.order
    ? await prisma.order.findUnique({
        where: { id: params.order },
        include: { items: true },
      })
    : null;

  return (
    <div className="mx-auto w-full max-w-4xl px-5 py-16 text-center sm:px-8">
      <div className="rounded-[2.4rem] border border-white/10 bg-white/5 p-8">
        <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
          Order received
        </p>
        <h1 className="mt-4 text-5xl font-semibold tracking-tight text-white">
          {params.demo ? "Demo checkout completed." : "Thanks for your order."}
        </h1>
        <p className="mt-4 text-base leading-8 text-white/62">
          Payment and shipment progress can now be tracked from the admin dashboard. This flow is
          wired for Stripe Checkout and also supports local demo testing without live keys.
        </p>
        {order ? (
          <div className="mt-8 rounded-[2rem] border border-white/10 bg-[#0a1525] p-6 text-left">
            <p className="text-sm uppercase tracking-[0.22em] text-white/50">{order.orderNumber}</p>
            <p className="mt-2 text-xl font-semibold text-white">{order.email}</p>
            <p className="mt-4 text-3xl font-semibold text-[#f7b733]">
              {formatCurrency(order.totalCents)}
            </p>
            <div className="mt-6 space-y-2 text-sm text-white/62">
              {order.items.map((item) => (
                <p key={item.id}>
                  {item.quantity} x {item.productTitle}
                </p>
              ))}
            </div>
          </div>
        ) : null}
        <div className="mt-8 flex flex-wrap justify-center gap-3">
          <Link
            href="/shop"
            className="inline-flex rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950"
          >
            Keep shopping
          </Link>
          <Link
            href="/admin"
            className="inline-flex rounded-full border border-white/12 bg-white/5 px-5 py-3 text-sm font-semibold text-white"
          >
            View admin dashboard
          </Link>
        </div>
      </div>
    </div>
  );
}
