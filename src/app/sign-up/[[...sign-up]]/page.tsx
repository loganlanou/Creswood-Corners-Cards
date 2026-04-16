import { SignUp } from "@clerk/nextjs";

import { isClerkConfigured } from "@/lib/env";

export default function SignUpPage() {
  if (!isClerkConfigured) {
    return null;
  }

  return (
    <div className="mx-auto flex min-h-[70vh] max-w-7xl items-center justify-center px-5 py-12 sm:px-8">
      <SignUp />
    </div>
  );
}
