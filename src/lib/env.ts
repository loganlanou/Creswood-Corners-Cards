const hasValue = (value: string | undefined) => Boolean(value && value.trim().length > 0);

export const appUrl =
  process.env.NEXT_PUBLIC_APP_URL?.replace(/\/$/, "") ?? "http://localhost:3000";

export const isClerkConfigured =
  hasValue(process.env.NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY) &&
  hasValue(process.env.CLERK_SECRET_KEY);

export const isStripeConfigured =
  hasValue(process.env.STRIPE_SECRET_KEY) &&
  hasValue(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY);

export const isDemoAdminMode = process.env.DEMO_ADMIN_MODE !== "false";

export const adminEmails = new Set(
  (process.env.ADMIN_EMAILS ?? "")
    .split(",")
    .map((email) => email.trim().toLowerCase())
    .filter(Boolean),
);
