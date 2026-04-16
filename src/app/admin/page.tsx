import {
  FulfillmentStatus,
  PaymentStatus,
  ProductStatus,
  ShippingStatus,
} from "@prisma/client";
import Link from "next/link";
import { redirect } from "next/navigation";

import {
  saveAdminGrantAction,
  saveLiveSessionAction,
  saveProductAction,
  updateOrderStatusAction,
} from "@/app/admin/actions";
import { listRegisteredUsers, getCurrentViewer } from "@/lib/auth";
import { isClerkConfigured, isDemoAdminMode, isStripeConfigured } from "@/lib/env";
import { getAdminDashboardData } from "@/lib/store";
import { formatCurrency } from "@/lib/utils";

export default async function AdminPage() {
  const viewer = await getCurrentViewer();

  if (!viewer.isAdmin) {
    redirect("/sign-in");
  }

  const [{ products, liveSession, orders, grants }, users] = await Promise.all([
    getAdminDashboardData(),
    listRegisteredUsers(),
  ]);

  const grossSales = orders.reduce((sum, order) => sum + order.totalCents, 0);

  return (
    <div className="mx-auto w-full max-w-7xl px-5 py-12 sm:px-8">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
            Admin
          </p>
          <h1 className="mt-4 text-5xl font-semibold tracking-tight text-white">
            Inventory, live selling, customers, and operations in one dashboard.
          </h1>
        </div>
        <div className="rounded-[1.5rem] border border-white/10 bg-white/5 px-5 py-4 text-sm text-white/68">
          <p>Clerk: {isClerkConfigured ? "configured" : `demo mode ${isDemoAdminMode ? "enabled" : "disabled"}`}</p>
          <p>Stripe: {isStripeConfigured ? "configured" : "demo checkout fallback"}</p>
        </div>
      </div>

      <section className="mt-8 grid gap-4 md:grid-cols-4">
        {[
          ["Products", `${products.length}`],
          ["Orders", `${orders.length}`],
          ["Registered accounts", `${users.length}`],
          ["Gross tracked sales", formatCurrency(grossSales)],
        ].map(([label, value]) => (
          <article key={label} className="rounded-[1.8rem] border border-white/10 bg-white/5 p-5">
            <p className="text-sm uppercase tracking-[0.22em] text-white/50">{label}</p>
            <p className="mt-3 text-3xl font-semibold text-white">{value}</p>
          </article>
        ))}
      </section>

      <div className="mt-10 grid gap-8 xl:grid-cols-[1.1fr_0.9fr]">
        <section className="space-y-8">
          <div className="rounded-[2rem] border border-white/10 bg-white/5 p-6">
            <h2 className="text-2xl font-semibold text-white">Live selling control</h2>
            <form action={saveLiveSessionAction} className="mt-5 grid gap-4">
              <input
                name="title"
                defaultValue={liveSession?.title ?? "Friday Night Football Heat Check"}
                placeholder="Stream title"
                className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm"
              />
              <input
                name="streamUrl"
                defaultValue={liveSession?.streamUrl ?? ""}
                placeholder="https://..."
                className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm"
              />
              <div className="grid gap-4 md:grid-cols-2">
                <input
                  name="platform"
                  defaultValue={liveSession?.platform ?? "Whatnot"}
                  placeholder="Platform"
                  className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm"
                />
                <label className="flex items-center gap-3 rounded-2xl border border-white/10 bg-[#0a1525] px-4">
                  <input
                    type="checkbox"
                    name="isActive"
                    defaultChecked={liveSession?.isActive ?? false}
                  />
                  <span className="text-sm text-white/72">Stream is live right now</span>
                </label>
              </div>
              <textarea
                name="pitch"
                defaultValue={liveSession?.pitch ?? ""}
                placeholder="What this live session is focused on"
                className="min-h-28 rounded-2xl border border-white/10 bg-[#0a1525] px-4 py-3 text-sm"
              />
              <textarea
                name="callout"
                defaultValue={liveSession?.callout ?? ""}
                placeholder="Short supporting callout for banners"
                className="min-h-24 rounded-2xl border border-white/10 bg-[#0a1525] px-4 py-3 text-sm"
              />
              <button className="inline-flex w-fit rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950">
                Save live settings
              </button>
            </form>
          </div>

          <div className="rounded-[2rem] border border-white/10 bg-white/5 p-6">
            <h2 className="text-2xl font-semibold text-white">Create or update product</h2>
            <form action={saveProductAction} className="mt-5 grid gap-4">
              <div className="grid gap-4 md:grid-cols-2">
                <input name="title" placeholder="Product title" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="slug" placeholder="slug-format-title" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
              </div>
              <div className="grid gap-4 md:grid-cols-2">
                <input name="player" placeholder="Player" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="team" placeholder="Team" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
              </div>
              <div className="grid gap-4 md:grid-cols-3">
                <input name="brand" placeholder="Brand" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="setName" placeholder="Set" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="cardNumber" placeholder="Card #" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
              </div>
              <div className="grid gap-4 md:grid-cols-4">
                <input name="year" type="number" placeholder="Year" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="grade" placeholder="Grade" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="condition" placeholder="Condition" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="accent" placeholder="amber / blue / emerald" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
              </div>
              <div className="grid gap-4 md:grid-cols-4">
                <input name="sport" defaultValue="Football" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="priceCents" type="number" placeholder="Price cents" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <input name="quantity" type="number" placeholder="Quantity" defaultValue={1} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
                <select name="status" defaultValue={ProductStatus.ACTIVE} className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm">
                  {Object.values(ProductStatus).map((status) => (
                    <option key={status} value={status}>
                      {status}
                    </option>
                  ))}
                </select>
              </div>
              <textarea
                name="summary"
                placeholder="Short card summary"
                className="min-h-24 rounded-2xl border border-white/10 bg-[#0a1525] px-4 py-3 text-sm"
              />
              <textarea
                name="description"
                placeholder="Longer product description"
                className="min-h-32 rounded-2xl border border-white/10 bg-[#0a1525] px-4 py-3 text-sm"
              />
              <div className="flex flex-wrap gap-6 text-sm text-white/72">
                <label className="flex items-center gap-2">
                  <input type="checkbox" name="featured" />
                  Featured
                </label>
                <label className="flex items-center gap-2">
                  <input type="checkbox" name="liveExclusive" />
                  Live priority
                </label>
              </div>
              <button className="inline-flex w-fit rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950">
                Save product
              </button>
            </form>
          </div>
        </section>

        <section className="space-y-8">
          <div className="rounded-[2rem] border border-white/10 bg-white/5 p-6">
            <h2 className="text-2xl font-semibold text-white">Admin access</h2>
            <p className="mt-2 text-sm leading-7 text-white/60">
              Grant admin access by email so additional Clerk accounts can manage inventory and operations.
            </p>
            <form action={saveAdminGrantAction} className="mt-5 grid gap-4">
              <input name="email" placeholder="admin@example.com" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
              <input name="label" placeholder="Optional label" className="h-12 rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm" />
              <button className="inline-flex w-fit rounded-full bg-white px-5 py-3 text-sm font-semibold text-slate-950">
                Grant admin access
              </button>
            </form>
            <div className="mt-6 space-y-3">
              {grants.map((grant) => (
                <div key={grant.id} className="rounded-2xl border border-white/10 bg-[#0a1525] px-4 py-3 text-sm text-white/68">
                  <p className="font-medium text-white">{grant.email}</p>
                  <p>{grant.label ?? "No label"}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="rounded-[2rem] border border-white/10 bg-white/5 p-6">
            <h2 className="text-2xl font-semibold text-white">Registered accounts</h2>
            <p className="mt-2 text-sm leading-7 text-white/60">
              Clerk users are surfaced here with order count and tracked spend so outreach and post-order communication can be managed centrally.
            </p>
            <div className="mt-5 space-y-3">
              {users.map((user) => (
                <div key={user.id} className="rounded-2xl border border-white/10 bg-[#0a1525] px-4 py-3 text-sm">
                  <div className="flex items-start justify-between gap-4">
                    <div>
                      <p className="font-medium text-white">{user.email}</p>
                      <p className="text-white/55">{user.fullName}</p>
                    </div>
                    <div className="text-right text-white/55">
                      <p>{user.orderCount} orders</p>
                      <p>{formatCurrency(user.totalSpentCents)}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>
      </div>

      <section className="mt-10 rounded-[2rem] border border-white/10 bg-white/5 p-6">
        <h2 className="text-2xl font-semibold text-white">Current products</h2>
        <div className="mt-5 overflow-x-auto">
          <table className="min-w-full text-left text-sm text-white/68">
            <thead className="text-white/45">
              <tr>
                <th className="pb-3 pr-6">Title</th>
                <th className="pb-3 pr-6">Price</th>
                <th className="pb-3 pr-6">Qty</th>
                <th className="pb-3 pr-6">Status</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <tr key={product.id} className="border-t border-white/8">
                  <td className="py-3 pr-6 text-white">{product.title}</td>
                  <td className="py-3 pr-6">{formatCurrency(product.priceCents)}</td>
                  <td className="py-3 pr-6">{product.quantity}</td>
                  <td className="py-3 pr-6">{product.status}</td>
                  <td className="py-3 pr-6">
                    <Link
                      href={`/admin/products/${product.id}`}
                      className="text-[#ffd16d] hover:text-[#f7b733]"
                    >
                      Edit
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section className="mt-10 rounded-[2rem] border border-white/10 bg-white/5 p-6">
        <h2 className="text-2xl font-semibold text-white">Orders and shipment progress</h2>
        <div className="mt-6 space-y-4">
          {orders.map((order) => (
            <form
              key={order.id}
              action={updateOrderStatusAction}
              className="rounded-[1.75rem] border border-white/10 bg-[#0a1525] p-5"
            >
              <input type="hidden" name="orderId" value={order.id} />
              <div className="flex flex-wrap items-start justify-between gap-4">
                <div>
                  <p className="text-lg font-semibold text-white">{order.orderNumber}</p>
                  <p className="text-sm text-white/55">{order.email}</p>
                </div>
                <p className="text-lg font-semibold text-[#f7b733]">
                  {formatCurrency(order.totalCents)}
                </p>
              </div>
              <div className="mt-4 grid gap-4 md:grid-cols-3">
                <select name="paymentStatus" defaultValue={order.paymentStatus} className="h-11 rounded-2xl border border-white/10 bg-white/5 px-4 text-sm">
                  {Object.values(PaymentStatus).map((status) => (
                    <option key={status} value={status}>
                      {status}
                    </option>
                  ))}
                </select>
                <select name="fulfillmentStatus" defaultValue={order.fulfillmentStatus} className="h-11 rounded-2xl border border-white/10 bg-white/5 px-4 text-sm">
                  {Object.values(FulfillmentStatus).map((status) => (
                    <option key={status} value={status}>
                      {status}
                    </option>
                  ))}
                </select>
                <select name="shippingStatus" defaultValue={order.shippingStatus} className="h-11 rounded-2xl border border-white/10 bg-white/5 px-4 text-sm">
                  {Object.values(ShippingStatus).map((status) => (
                    <option key={status} value={status}>
                      {status}
                    </option>
                  ))}
                </select>
              </div>
              <textarea
                name="notes"
                defaultValue={order.notes ?? ""}
                placeholder="Packing or shipment notes"
                className="mt-4 min-h-24 w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm"
              />
              <button className="mt-4 inline-flex rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950">
                Update order
              </button>
            </form>
          ))}
        </div>
      </section>
    </div>
  );
}
