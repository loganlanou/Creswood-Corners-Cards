import { NextResponse } from "next/server";
import { OrderSource, PaymentStatus, ProductStatus } from "@prisma/client";
import { z } from "zod";

import { getCurrentViewer } from "@/lib/auth";
import { appUrl } from "@/lib/env";
import { prisma } from "@/lib/prisma";
import { getStripeServer } from "@/lib/stripe";
import { createOrderNumber, normalizeEmail } from "@/lib/utils";

const payloadSchema = z.object({
  customerName: z.string().optional(),
  email: z.string().email(),
  items: z
    .array(
      z.object({
        productId: z.string(),
        quantity: z.number().int().min(1).max(10),
      }),
    )
    .min(1),
});

export async function POST(request: Request) {
  const body = await request.json();
  const parsed = payloadSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json({ error: "A valid email and at least one cart item are required." }, { status: 400 });
  }

  const viewer = await getCurrentViewer();
  const ids = parsed.data.items.map((item) => item.productId);
  const products = await prisma.product.findMany({
    where: {
      id: { in: ids },
      status: ProductStatus.ACTIVE,
    },
  });

  const productMap = new Map(products.map((product) => [product.id, product]));
  const lineItems = parsed.data.items
    .map((item) => {
      const product = productMap.get(item.productId);
      if (!product) {
        return null;
      }

      return {
        product,
        quantity: item.quantity,
      };
    })
    .filter(
      (
        entry,
      ): entry is {
        product: (typeof products)[number];
        quantity: number;
      } => Boolean(entry),
    );

  if (lineItems.length === 0) {
    return NextResponse.json({ error: "No valid products were found in the cart." }, { status: 400 });
  }

  const subtotalCents = lineItems.reduce(
    (sum, entry) => sum + entry.product.priceCents * entry.quantity,
    0,
  );

  const order = await prisma.order.create({
    data: {
      orderNumber: createOrderNumber(),
      email: normalizeEmail(parsed.data.email),
      customerName: parsed.data.customerName || null,
      clerkUserId: viewer.clerkUserId,
      source: process.env.STRIPE_SECRET_KEY ? OrderSource.WEB : OrderSource.WEB_DEMO,
      paymentStatus: process.env.STRIPE_SECRET_KEY ? PaymentStatus.PENDING : PaymentStatus.PAID,
      subtotalCents,
      totalCents: subtotalCents,
      items: {
        create: lineItems.map((entry) => ({
          productId: entry.product.id,
          productTitle: entry.product.title,
          productSlug: entry.product.slug,
          quantity: entry.quantity,
          unitPriceCents: entry.product.priceCents,
        })),
      },
    },
  });

  const stripe = getStripeServer();
  if (!stripe) {
    return NextResponse.json({
      url: `${appUrl}/checkout/success?order=${order.id}&demo=1`,
    });
  }

  const session = await stripe.checkout.sessions.create({
    mode: "payment",
    customer_email: parsed.data.email,
    success_url: `${appUrl}/checkout/success?order=${order.id}`,
    cancel_url: `${appUrl}/cart`,
    metadata: {
      orderId: order.id,
    },
    line_items: lineItems.map((entry) => ({
      quantity: entry.quantity,
      price_data: {
        currency: "usd",
        unit_amount: entry.product.priceCents,
        product_data: {
          name: entry.product.title,
          description: entry.product.summary,
        },
      },
    })),
  });

  await prisma.order.update({
    where: { id: order.id },
    data: { stripeSessionId: session.id },
  });

  return NextResponse.json({ url: session.url });
}
