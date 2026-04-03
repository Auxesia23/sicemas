<script lang="ts">
    import "../layout.css";
    import { goto } from "$app/navigation";
    import { auth } from "$lib/states/auth.svelte";
    import { page } from "$app/state";
    import { slide, fade } from "svelte/transition";

    let { children } = $props();

    // State
    let isMenuOpen = $state(false);
    let isUserMenuOpen = $state(false);

    // Derived
    let currentPath = $derived(page.url.pathname);

    const navLinks = [
        {
            href: "/",
            label: "Beranda",
            icon: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6",
        },
        {
            href: "/situs",
            label: "Direktori",
            icon: "M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4",
        },
    ];

    function isActive(href: string): boolean {
        if (href === "/") return currentPath === "/";
        return currentPath.startsWith(href);
    }

    function handleNavClick(href: string) {
        isMenuOpen = false;
        goto(href);
    }
</script>

<nav
    class="navbar sticky top-4 z-50 bg-base-100/70 backdrop-blur-xl border border-base-200/50 px-4 max-w-7xl mx-auto rounded-2xl shadow-xl shadow-base-200/40 transition-all duration-300 h-16"
    aria-label="Main navigation"
>
    <div class="container mx-auto flex items-center justify-between w-full">
        <div class="flex items-center">
            <a
                href="/"
                class="btn btn-ghost hover:bg-transparent px-2 gap-3 group transition-all"
                aria-label="SICEMAS Beranda"
            >
                <div
                    class="relative w-10 h-10 flex items-center justify-center bg-linear-to-br from-primary to-secondary rounded-xl shadow-lg shadow-primary/20 transition-transform group-hover:scale-105"
                >
                    <img src="/kemenag.webp" alt="SICEMAS" class="w-6 h-6" />
                </div>
                <span
                    class="text-xl font-black tracking-tighter bg-linear-to-r from-primary to-secondary bg-clip-text text-transparent"
                >
                    SICEMAS
                </span>
            </a>
        </div>

        <div
            class="hidden lg:flex items-center gap-1 bg-base-200/50 p-1 rounded-xl border border-base-200/50"
        >
            {#each navLinks as link}
                <a
                    href={link.href}
                    class="btn btn-ghost btn-sm rounded-lg px-4 font-bold transition-all {isActive(
                        link.href,
                    )
                        ? 'bg-base-100 text-primary shadow-sm'
                        : 'text-base-content/60 hover:text-primary hover:bg-base-100/50'}"
                    aria-current={isActive(link.href) ? "page" : undefined}
                >
                    {link.label}
                </a>
            {/each}
        </div>

        <div class="flex items-center gap-3">
            {#if auth.user}
                <div class="relative">
                    <button
                        class="flex items-center gap-2 p-1.5 pr-3 rounded-xl hover:bg-base-200/50 transition-colors border border-transparent hover:border-base-200"
                        onclick={() => (isUserMenuOpen = !isUserMenuOpen)}
                    >
                        <div
                            class="w-8 h-8 rounded-lg bg-linear-to-br from-primary to-secondary flex items-center justify-center text-white text-sm font-black shadow-md shadow-primary/20"
                        >
                            {auth.user.nama_lengkap?.charAt(0).toUpperCase()}
                        </div>
                        <span
                            class="text-sm font-bold text-base-content hidden sm:block"
                        >
                            {auth.user.nama_lengkap.split(" ")[0]}
                        </span>
                        <svg
                            class="w-4 h-4 text-base-content/40 transition-transform {isUserMenuOpen
                                ? 'rotate-180'
                                : ''}"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="3"
                                d="M19 9l-7 7-7-7"
                            />
                        </svg>
                    </button>

                    {#if isUserMenuOpen}
                        <div
                            class="absolute right-0 mt-3 w-52 bg-base-100 rounded-2xl shadow-2xl border border-base-200 py-2 z-50 overflow-hidden"
                            transition:fade={{ duration: 150 }}
                        >
                            <div
                                class="px-4 py-2 mb-1 border-b border-base-200 bg-base-200/30"
                            >
                                <p
                                    class="text-[10px] font-black uppercase tracking-widest text-base-content/40"
                                >
                                    Sesi Aktif
                                </p>
                                <p
                                    class="text-sm font-bold truncate text-base-content"
                                >
                                    {auth.user.nama_lengkap}
                                </p>
                            </div>
                            <a
                                href="/dashboard"
                                class="flex items-center gap-3 px-4 py-2.5 text-sm font-bold text-base-content/70 hover:bg-primary/5 hover:text-primary transition-colors"
                            >
                                <svg
                                    class="w-4 h-4"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                    ><path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
                                    /></svg
                                >
                                Dashboard
                            </a>
                            <button
                                onclick={() => goto("/logout")}
                                class="flex items-center gap-3 px-4 py-2.5 text-sm font-bold text-error hover:bg-error/5 w-full text-left transition-colors"
                            >
                                <svg
                                    class="w-4 h-4"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                    ><path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
                                    /></svg
                                >
                                Keluar
                            </button>
                        </div>
                    {/if}
                </div>
            {:else}
                <a
                    href="/login"
                    class="btn btn-primary btn-sm rounded-xl px-6 font-bold shadow-lg shadow-primary/20 transition-all hover:-translate-y-0.5"
                >
                    Masuk
                </a>
            {/if}

            <button
                class="lg:hidden btn btn-ghost btn-square rounded-xl bg-base-200/50"
                onclick={() => (isMenuOpen = !isMenuOpen)}
                aria-label="Toggle menu"
            >
                <svg
                    class="w-6 h-6"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2.5"
                        d={isMenuOpen
                            ? "M6 18L18 6M6 6l12 12"
                            : "M4 6h16M4 12h16M4 18h16"}
                    />
                </svg>
            </button>
        </div>
    </div>

    {#if isMenuOpen}
        <div
            class="lg:hidden absolute top-[calc(100%+1rem)] left-0 right-0 bg-base-100 rounded-2xl border border-base-200 shadow-2xl p-4 space-y-2 z-50"
            transition:slide={{ duration: 200 }}
        >
            {#each navLinks as link}
                <a
                    href={link.href}
                    class="flex items-center gap-4 px-4 py-3 rounded-xl font-bold transition-all {isActive(
                        link.href,
                    )
                        ? 'bg-primary/10 text-primary'
                        : 'text-base-content/60 hover:bg-base-200'}"
                    onclick={() => (isMenuOpen = false)}
                >
                    <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d={link.icon}
                        />
                    </svg>
                    {link.label}
                </a>
            {/each}

            {#if !auth.user}
                <a href="/login" class="btn btn-primary w-full rounded-xl mt-4"
                    >Masuk ke Sistem</a
                >
            {/if}
        </div>
    {/if}
</nav>

<main class="relative z-0">
    {@render children()}
</main>

<footer
    class="relative z-10 bg-base-100/60 backdrop-blur-lg py-12 border-t border-base-200"
>
    <div class="container mx-auto px-4">
        <div class="flex flex-col items-center mb-8 text-center">
            <div class="flex items-center gap-2 font-bold text-xl mb-3">
                <div
                    class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center text-white text-sm"
                >
                    <img src="/kemenag.webp" alt="SICEMAS" class="w-6 h-6" />
                </div>
                <span>SICEMAS</span>
            </div>

            <p class="text-base-content/70 max-w-sm mb-6 text-sm">
                Sistem Informasi Catatan Elektronik Masjid Area Ciemas -
                Inisiatif digitalisasi data keagamaan oleh KUA Kecamatan Ciemas.
            </p>

            <div class="flex gap-4">
                <a
                    href="https://www.tiktok.com/@kuaciemas?is_from_webapp=1&sender_device=pc"
                    aria-label="TikTok"
                    class="hover:text-primary transition-colors"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        class="fill-current"
                    >
                        <path
                            d="M12.53.02C13.84 0 15.14.01 16.44 0c.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z"
                        ></path>
                    </svg>
                </a>
                <a
                    href="https://www.instagram.com/kuaciemas_official/"
                    aria-label="Instagram"
                    class="hover:text-primary transition-colors"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        class="fill-current"
                    >
                        <path
                            d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zM12 0C8.741 0 8.333.014 7.053.072 2.695.272.273 2.69.073 7.052.014 8.333 0 8.741 0 12c0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98C8.333 23.986 8.741 24 12 24c3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98C15.668.014 15.259 0 12 0zm0 5.838a6.162 6.162 0 100 12.324 6.162 6.162 0 000-12.324zM12 16a4 4 0 110-8 4 4 0 010 8zm6.406-11.845a1.44 1.44 0 100 2.881 1.44 1.44 0 000-2.881z"
                        ></path>
                    </svg>
                </a>
            </div>
        </div>

        <div
            class="border-t border-base-200 pt-6 flex flex-col md:flex-row justify-between items-center gap-6 text-sm text-base-content/60"
        >
            <div class="text-center md:text-left">
                <p class="font-medium text-base-content/80">
                    &copy; {new Date().getFullYear()} SICEMAS - KUA Kecamatan Ciemas.
                </p>
                <p class="mt-1">
                    Dikembangkan sebagai Proyek Skripsi oleh <span
                        class="font-semibold">Fajar Gustiana</span
                    > (Teknik Informatika, Unikom).
                </p>
            </div>

            <div class="flex items-center gap-4">
                <a
                    href="https://github.com/Auxesia23"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="flex items-center gap-2 hover:text-primary transition-colors tooltip tooltip-top"
                    data-tip="Lihat Source Code di GitHub"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        class="fill-current"
                    >
                        <path
                            d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
                        />
                    </svg>
                    <span class="font-medium text-sm">@Auxesia23</span>
                </a>
            </div>
        </div>
    </div>
</footer>

<style>
    .btn-ghost[aria-current="page"] {
        position: relative;
    }

    .btn-ghost[aria-current="page"]::after {
        content: "";
        position: absolute;
        bottom: 0;
        left: 50%;
        transform: translateX(-50%);
        width: 20px;
        height: 2px;
        background: currentColor;
        border-radius: 2px;
    }

    .navbar {
        transition:
            background-color 0.3s ease,
            backdrop-filter 0.3s ease,
            box-shadow 0.3s ease;
    }

    :global(html) {
        scroll-padding-top: 5rem;
    }
</style>
