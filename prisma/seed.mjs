import pkg from "@prisma/client";
import { PrismaBetterSqlite3 } from "@prisma/adapter-better-sqlite3";

const { PrismaClient, ProductStatus } = pkg;

const adapter = new PrismaBetterSqlite3({
  url: process.env.DATABASE_URL || "file:./dev.db",
});

const prisma = new PrismaClient({ adapter });

const products = [
  {
    slug: "jayden-daniels-prizm-rc-gold-wave",
    title: "Jayden Daniels Prizm RC Gold Wave",
    summary: "Low-number rookie color built for live hype and premium web buyers.",
    description:
      "A featured rookie card for your flagship football launch. This listing is structured the way modern card shops sell: clean title, quick condition notes, and immediate buying confidence.",
    player: "Jayden Daniels",
    team: "Washington Commanders",
    brand: "Panini Prizm",
    setName: "Prizm Football",
    year: 2024,
    cardNumber: "312",
    grade: "Raw",
    condition: "Near Mint",
    priceCents: 32900,
    quantity: 1,
    featured: true,
    liveExclusive: true,
    status: ProductStatus.ACTIVE,
    accent: "amber",
    sortOrder: 10,
  },
  {
    slug: "cj-stroud-downtown-sgc-10",
    title: "C.J. Stroud Downtown SGC 10",
    summary: "High-visibility centerpiece slab for homepage feature modules.",
    description:
      "Premier slab inventory is critical for trust and average order value. This card gives the storefront a premium anchor product for both home and live traffic.",
    player: "C.J. Stroud",
    team: "Houston Texans",
    brand: "Donruss",
    setName: "Downtown",
    year: 2023,
    cardNumber: "DT-7",
    grade: "SGC 10",
    condition: "Gem Mint",
    priceCents: 79900,
    quantity: 1,
    featured: true,
    status: ProductStatus.ACTIVE,
    accent: "crimson",
    sortOrder: 20,
  },
  {
    slug: "bo-nix-mosaic-genesis-rc",
    title: "Bo Nix Mosaic Genesis RC",
    summary: "A clean chase rookie to represent mid-tier impulse buys.",
    description:
      "Football inventory needs accessible chase cards, not just grails. This product fills that role and supports fast add-to-cart behavior.",
    player: "Bo Nix",
    team: "Denver Broncos",
    brand: "Panini Mosaic",
    setName: "Mosaic Football",
    year: 2024,
    cardNumber: "287",
    grade: "Raw",
    condition: "Near Mint-Mint",
    priceCents: 18500,
    quantity: 2,
    featured: true,
    status: ProductStatus.ACTIVE,
    accent: "emerald",
    sortOrder: 30,
  },
  {
    slug: "drake-maye-select-field-level-auto",
    title: "Drake Maye Select Field Level Auto",
    summary: "Autograph inventory ready for both direct checkout and stream mentions.",
    description:
      "A balanced autograph listing that supports both storefront conversion and live mention cross-sell when the stream is active.",
    player: "Drake Maye",
    team: "New England Patriots",
    brand: "Panini Select",
    setName: "Select Football",
    year: 2024,
    cardNumber: "FLA-DM",
    grade: "Raw",
    condition: "Near Mint",
    priceCents: 26900,
    quantity: 1,
    status: ProductStatus.ACTIVE,
    accent: "blue",
    sortOrder: 40,
  },
  {
    slug: "malik-nabers-optic-holo-rc",
    title: "Malik Nabers Optic Holo RC",
    summary: "Popular rookie profile that makes category rails feel current.",
    description:
      "Designed as a flexible inventory piece for homepage, category, and live merchandising modules.",
    player: "Malik Nabers",
    team: "New York Giants",
    brand: "Donruss Optic",
    setName: "Optic Football",
    year: 2024,
    cardNumber: "214",
    grade: "Raw",
    condition: "Near Mint",
    priceCents: 12900,
    quantity: 3,
    status: ProductStatus.ACTIVE,
    accent: "violet",
    sortOrder: 50,
  },
  {
    slug: "rome-odunze-prizm-silver-rc",
    title: "Rome Odunze Prizm Silver RC",
    summary: "An affordable parallel that keeps the grid approachable.",
    description:
      "Not every listing should be a high-end piece. This card gives the catalog a healthier price spread and supports basket building.",
    player: "Rome Odunze",
    team: "Chicago Bears",
    brand: "Panini Prizm",
    setName: "Prizm Football",
    year: 2024,
    cardNumber: "341",
    grade: "Raw",
    condition: "Near Mint",
    priceCents: 8900,
    quantity: 4,
    status: ProductStatus.ACTIVE,
    accent: "slate",
    sortOrder: 60,
  },
];

async function main() {
  await prisma.liveSession.upsert({
    where: { id: "launch-live-session" },
    update: {
      title: "Friday Night Football Heat Check",
      pitch: "Live singles, rookie color, and quick-hit auctions built to move inventory fast.",
      callout: "Running live now on Whatnot and mirrored from the store banner.",
      streamUrl: "https://www.whatnot.com",
      platform: "Whatnot",
      isActive: true,
    },
    create: {
      id: "launch-live-session",
      title: "Friday Night Football Heat Check",
      pitch: "Live singles, rookie color, and quick-hit auctions built to move inventory fast.",
      callout: "Running live now on Whatnot and mirrored from the store banner.",
      streamUrl: "https://www.whatnot.com",
      platform: "Whatnot",
      isActive: true,
    },
  });

  for (const product of products) {
    await prisma.product.upsert({
      where: { slug: product.slug },
      update: product,
      create: product,
    });
  }
}

main()
  .then(async () => {
    await prisma.$disconnect();
  })
  .catch(async (error) => {
    console.error(error);
    await prisma.$disconnect();
    process.exit(1);
  });
