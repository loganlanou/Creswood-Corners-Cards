import { ProductStatus } from "@prisma/client";

import { prisma } from "@/lib/prisma";

export async function getActiveLiveSession() {
  return prisma.liveSession.findFirst({
    where: { isActive: true },
    orderBy: { updatedAt: "desc" },
  });
}

export async function getFeaturedProducts() {
  return prisma.product.findMany({
    where: { status: ProductStatus.ACTIVE, featured: true },
    orderBy: [{ sortOrder: "asc" }, { updatedAt: "desc" }],
    take: 4,
  });
}

export async function getCatalogProducts() {
  return prisma.product.findMany({
    where: { status: ProductStatus.ACTIVE },
    orderBy: [{ sortOrder: "asc" }, { updatedAt: "desc" }],
  });
}

export async function getProductBySlug(slug: string) {
  return prisma.product.findUnique({
    where: { slug },
  });
}

export async function getAdminDashboardData() {
  const [products, liveSession, orders, grants] = await Promise.all([
    prisma.product.findMany({
      orderBy: [{ sortOrder: "asc" }, { updatedAt: "desc" }],
    }),
    prisma.liveSession.findFirst({
      orderBy: { updatedAt: "desc" },
    }),
    prisma.order.findMany({
      include: { items: true },
      orderBy: { createdAt: "desc" },
      take: 25,
    }),
    prisma.adminGrant.findMany({
      orderBy: { createdAt: "desc" },
    }),
  ]);

  return { products, liveSession, orders, grants };
}
