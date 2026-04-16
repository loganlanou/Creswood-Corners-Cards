"use server";

import {
  FulfillmentStatus,
  PaymentStatus,
  ProductStatus,
  ShippingStatus,
} from "@prisma/client";
import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import { z } from "zod";

import { getCurrentViewer } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import { normalizeEmail } from "@/lib/utils";

const productSchema = z.object({
  id: z.string().optional(),
  title: z.string().min(3),
  slug: z
    .string()
    .min(3)
    .regex(/^[a-z0-9-]+$/),
  summary: z.string().min(10),
  description: z.string().min(20),
  sport: z.string().min(3),
  player: z.string().optional(),
  team: z.string().optional(),
  brand: z.string().optional(),
  setName: z.string().optional(),
  year: z.coerce.number().int().optional(),
  cardNumber: z.string().optional(),
  grade: z.string().optional(),
  condition: z.string().optional(),
  priceCents: z.coerce.number().int().min(100),
  quantity: z.coerce.number().int().min(0),
  accent: z.string().optional(),
  status: z.nativeEnum(ProductStatus),
  featured: z.boolean(),
  liveExclusive: z.boolean(),
});

async function requireAdminAction() {
  const viewer = await getCurrentViewer();
  if (!viewer.isAdmin) {
    redirect("/sign-in");
  }
}

export async function saveProductAction(formData: FormData) {
  await requireAdminAction();

  const parsed = productSchema.parse({
    id: formData.get("id") || undefined,
    title: formData.get("title"),
    slug: formData.get("slug"),
    summary: formData.get("summary"),
    description: formData.get("description"),
    sport: formData.get("sport"),
    player: formData.get("player") || undefined,
    team: formData.get("team") || undefined,
    brand: formData.get("brand") || undefined,
    setName: formData.get("setName") || undefined,
    year: formData.get("year") || undefined,
    cardNumber: formData.get("cardNumber") || undefined,
    grade: formData.get("grade") || undefined,
    condition: formData.get("condition") || undefined,
    priceCents: Number(formData.get("priceCents")),
    quantity: Number(formData.get("quantity")),
    accent: formData.get("accent") || undefined,
    status: formData.get("status"),
    featured: formData.get("featured") === "on",
    liveExclusive: formData.get("liveExclusive") === "on",
  });

  const payload = {
    ...parsed,
    player: parsed.player || null,
    team: parsed.team || null,
    brand: parsed.brand || null,
    setName: parsed.setName || null,
    year: parsed.year || null,
    cardNumber: parsed.cardNumber || null,
    grade: parsed.grade || null,
    condition: parsed.condition || null,
    accent: parsed.accent || null,
  };

  if (parsed.id) {
    await prisma.product.update({
      where: { id: parsed.id },
      data: payload,
    });
  } else {
    await prisma.product.create({
      data: payload,
    });
  }

  revalidatePath("/");
  revalidatePath("/shop");
  revalidatePath("/live");
  revalidatePath("/admin");
}

export async function saveLiveSessionAction(formData: FormData) {
  await requireAdminAction();

  const existing = await prisma.liveSession.findFirst({ orderBy: { updatedAt: "desc" } });
  const payload = {
    title: String(formData.get("title") || "Live Selling Session"),
    pitch: String(formData.get("pitch") || ""),
    callout: String(formData.get("callout") || "") || null,
    streamUrl: String(formData.get("streamUrl") || ""),
    platform: String(formData.get("platform") || "Whatnot"),
    isActive: formData.get("isActive") === "on",
  };

  if (existing) {
    await prisma.liveSession.update({
      where: { id: existing.id },
      data: payload,
    });
  } else {
    await prisma.liveSession.create({ data: payload });
  }

  revalidatePath("/");
  revalidatePath("/live");
  revalidatePath("/admin");
}

export async function saveAdminGrantAction(formData: FormData) {
  await requireAdminAction();

  const email = normalizeEmail(String(formData.get("email") || ""));
  if (!email) {
    return;
  }

  await prisma.adminGrant.upsert({
    where: { email },
    update: {
      label: String(formData.get("label") || "") || null,
    },
    create: {
      email,
      label: String(formData.get("label") || "") || null,
    },
  });

  revalidatePath("/admin");
}

export async function updateOrderStatusAction(formData: FormData) {
  await requireAdminAction();

  const orderId = String(formData.get("orderId"));
  await prisma.order.update({
    where: { id: orderId },
    data: {
      paymentStatus: formData.get("paymentStatus") as PaymentStatus,
      fulfillmentStatus: formData.get("fulfillmentStatus") as FulfillmentStatus,
      shippingStatus: formData.get("shippingStatus") as ShippingStatus,
      notes: String(formData.get("notes") || "") || null,
    },
  });

  revalidatePath("/admin");
}
