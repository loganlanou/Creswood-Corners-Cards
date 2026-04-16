export function SiteFooter() {
  return (
    <footer className="border-t border-white/8 bg-[#050c16]">
      <div className="mx-auto grid w-full max-w-7xl gap-10 px-5 py-12 sm:px-8 md:grid-cols-[1.2fr_0.8fr_0.8fr]">
        <div>
          <p className="font-display text-2xl uppercase tracking-[0.16em] text-white">Creswood</p>
          <p className="mt-4 max-w-md text-sm leading-7 text-white/60">
            Built for football singles, live stream selling, and a cleaner collector experience
            than the typical cluttered card storefront.
          </p>
        </div>
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-white/75">Operations</p>
          <ul className="mt-4 space-y-3 text-sm text-white/60">
            <li>Stripe checkout and demo local testing</li>
            <li>Admin-managed inventory and live links</li>
            <li>Order, shipment, and customer visibility</li>
          </ul>
        </div>
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-white/75">Launch Path</p>
          <ul className="mt-4 space-y-3 text-sm text-white/60">
            <li>Localhost first</li>
            <li>Vercel for preliminary rollout</li>
            <li>Custom POS-ready backend structure</li>
          </ul>
        </div>
      </div>
    </footer>
  );
}
