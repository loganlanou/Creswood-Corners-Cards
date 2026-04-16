import { CartView } from "@/components/cart-view";

export default function CartPage() {
  return (
    <div className="mx-auto w-full max-w-7xl px-5 py-12 sm:px-8">
      <div className="max-w-3xl">
        <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[#ffd16d]">
          Cart
        </p>
        <h1 className="mt-4 text-5xl font-semibold tracking-tight text-white">
          Convert stream energy into a clean checkout flow.
        </h1>
      </div>
      <div className="mt-10">
        <CartView />
      </div>
    </div>
  );
}
