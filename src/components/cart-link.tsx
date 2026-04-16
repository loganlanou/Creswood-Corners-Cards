"use client";

import Link from "next/link";
import { ShoppingBag } from "lucide-react";

import { useCart } from "@/components/cart-provider";

export function CartLink() {
  const { count } = useCart();

  return (
    <Link
      href="/cart"
      className="inline-flex items-center gap-2 rounded-full border border-white/12 bg-white/7 px-4 py-2 text-sm font-medium text-white transition hover:bg-white/12"
    >
      <ShoppingBag className="size-4" />
      Cart
      <span className="rounded-full bg-white/12 px-2 py-0.5 text-xs">{count}</span>
    </Link>
  );
}
