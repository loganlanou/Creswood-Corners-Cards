import { clerkClient, currentUser } from "@clerk/nextjs/server";

import { adminEmails, isClerkConfigured, isDemoAdminMode } from "@/lib/env";
import { prisma } from "@/lib/prisma";
import { normalizeEmail } from "@/lib/utils";

export async function getCurrentViewer() {
  if (!isClerkConfigured) {
    return {
      isAuthenticated: false,
      isAdmin: isDemoAdminMode,
      email: null as string | null,
      clerkUserId: null as string | null,
      mode: "demo" as const,
    };
  }

  const user = await currentUser();
  const email = user?.emailAddresses.find((entry) => entry.id === user.primaryEmailAddressId)?.emailAddress ??
    user?.emailAddresses[0]?.emailAddress ??
    null;

  const normalizedEmail = email ? normalizeEmail(email) : null;
  const grant = normalizedEmail
    ? await prisma.adminGrant.findUnique({ where: { email: normalizedEmail } })
    : null;

  return {
    isAuthenticated: Boolean(user),
    isAdmin: Boolean(normalizedEmail && (adminEmails.has(normalizedEmail) || grant)),
    email,
    clerkUserId: user?.id ?? null,
    mode: "clerk" as const,
  };
}

export async function listRegisteredUsers() {
  const orderAggregates = await prisma.order.groupBy({
    by: ["email"],
    _count: { email: true },
    _sum: { totalCents: true },
  });

  const orderMap = new Map(
    orderAggregates.map((entry) => [
      normalizeEmail(entry.email),
      {
        orderCount: entry._count.email,
        totalSpentCents: entry._sum.totalCents ?? 0,
      },
    ]),
  );

  if (!isClerkConfigured) {
    return Array.from(orderMap.entries()).map(([email, stats]) => ({
      id: email,
      email,
      fullName: "Local demo customer",
      createdAt: null as Date | null,
      ...stats,
    }));
  }

  const client = await clerkClient();
  const users = await client.users.getUserList({ limit: 100 });

  return users.data.map((user) => {
    const email =
      user.emailAddresses.find((entry) => entry.id === user.primaryEmailAddressId)?.emailAddress ??
      user.emailAddresses[0]?.emailAddress ??
      "unknown@example.com";
    const normalizedEmail = normalizeEmail(email);
    const stats = orderMap.get(normalizedEmail) ?? {
      orderCount: 0,
      totalSpentCents: 0,
    };

    return {
      id: user.id,
      email,
      fullName: [user.firstName, user.lastName].filter(Boolean).join(" ") || "Registered account",
      createdAt: user.createdAt ? new Date(user.createdAt) : null,
      ...stats,
    };
  });
}
