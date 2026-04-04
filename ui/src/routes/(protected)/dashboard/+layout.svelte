<script lang="ts">
    import { onMount } from "svelte";
    import { auth } from "$lib/states/auth.svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/state";

    let { children } = $props();

    let isMobileMenuOpen = $state(false);
    let user = $derived(auth.user!);
    let currentPath = $derived(page.url.pathname);

    const navigate = (path: string) => {
        goto(path);
        isMobileMenuOpen = false;
    };

    function getInitials(name: string = ""): string {
        if (!name) return "??";
        return name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2);
    }

    function isActive(path: string): boolean {
        if (path === "/dashboard") return currentPath === "/dashboard";
        return currentPath.startsWith(path);
    }

    const menuGroups = [
        {
            title: "Menu Utama",
            items: [
                {
                    path: "/dashboard",
                    label: "Dashboard",
                    icon: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6",
                    permissions: ["user:create", "situs:verify"],
                },
            ],
        },
        {
            title: "Manajemen Data",
            items: [
                {
                    path: "/dashboard/situs",
                    label: "Daftar Situs",
                    icon: "M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4",
                    permissions: ["situs:read_all", "situs:read_own"],
                },
                {
                    path: "/dashboard/jenis-situs",
                    label: "Jenis Situs",
                    icon: "M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z",
                    permissions: ["jenis-situs:read"],
                },
                {
                    path: "/dashboard/petugas",
                    label: "Petugas",
                    icon: "M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z",
                    permissions: ["user:read"],
                },
            ],
        },
        {
            title: "Sistem",
            items: [
                {
                    path: "/dashboard/policy",
                    label: "Hak Akses",
                    icon: "M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z",
                    permissions: ["policy:read"],
                },
                {
                    path: "/dashboard/profile",
                    label: "Profil",
                    icon: "M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z",
                    permissions: null,
                },
            ],
        },
    ];
</script>

<div class="min-h-screen bg-base-200 flex">
    <button
        type="button"
        class="lg:hidden fixed top-4 left-4 z-50 btn btn-circle btn-primary shadow-lg"
        onclick={() => (isMobileMenuOpen = !isMobileMenuOpen)}
        aria-label={isMobileMenuOpen ? "Tutup menu" : "Buka menu"}
        aria-expanded={isMobileMenuOpen}
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
        >
            {#if isMobileMenuOpen}
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                />
            {:else}
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 6h16M4 12h16M4 18h16"
                />
            {/if}
        </svg>
    </button>

    <aside
        class="fixed lg:sticky top-0 h-dvh z-40 w-72 bg-base-100 border-r border-base-200
        transition-transform duration-300 ease-in-out flex flex-col
        {isMobileMenuOpen
            ? 'translate-x-0'
            : '-translate-x-full lg:translate-x-0'}"
    >
        <div
            class="p-6 border-b border-base-200 bg-linear-to-br from-base-100 to-base-200/50"
        >
            <div class="flex items-center gap-4">
                <div
                    class="w-12 h-12 rounded-xl bg-linear-to-br from-primary to-secondary flex items-center justify-center shadow-lg"
                >
                    <img
                        src="/kemenag.webp"
                        alt="KUA"
                        class="h-8 w-8 object-contain"
                    />
                </div>
                <div>
                    <h3 class="text-lg font-bold text-primary leading-tight">
                        KUA Ciemas
                    </h3>
                    <p class="text-xs text-base-content/60 font-medium">
                        SICEMAS Admin
                    </p>
                </div>
            </div>
        </div>

        <div class="p-4 mx-4 mt-4 mb-2">
            <div class="bg-base-200/50 rounded-xl p-4 border border-base-200">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-linear-to-br from-primary to-secondary flex items-center justify-center text-white text-sm font-bold shadow-md"
                    >
                        {getInitials(user.nama_lengkap)}
                    </div>
                    <div class="flex-1 min-w-0">
                        <p
                            class="text-sm font-semibold text-base-content truncate"
                        >
                            {user.nama_lengkap}
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <div class="flex-1 overflow-y-auto py-4 px-3 space-y-6">
            {#each menuGroups as group}
                {@const visibleItems = group.items.filter(
                    (item) =>
                        !item.permissions ||
                        auth.hasAnyPermission(item.permissions),
                )}

                {#if visibleItems.length > 0}
                    <div>
                        <h4
                            class="px-3 mb-2 text-xs font-bold text-base-content/40 uppercase tracking-wider"
                        >
                            {group.title}
                        </h4>
                        <ul class="space-y-1">
                            {#each visibleItems as item}
                                <li>
                                    <button
                                        type="button"
                                        onclick={() => navigate(item.path)}
                                        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200 group text-left {isActive(
                                            item.path,
                                        )
                                            ? 'bg-primary/10 text-primary font-medium'
                                            : 'text-base-content/70 hover:bg-base-200 hover:text-base-content'}"
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-5 w-5 transition-transform group-hover:scale-110 {isActive(
                                                item.path,
                                            )
                                                ? 'text-primary'
                                                : 'text-base-content/50 group-hover:text-base-content'}"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d={item.icon}
                                            />
                                        </svg>
                                        <span class="flex-1">{item.label}</span>
                                        {#if isActive(item.path)}
                                            <div
                                                class="w-1.5 h-1.5 rounded-full bg-primary animate-pulse"
                                            ></div>
                                        {/if}
                                    </button>
                                </li>
                            {/each}
                        </ul>
                    </div>
                {/if}
            {/each}
        </div>

        <div class="p-4 border-t border-base-200 bg-base-200/30">
            <div class="space-y-2">
                <button
                    type="button"
                    onclick={() => goto("/logout")}
                    class="w-full flex items-center gap-3 px-3 py-2 rounded-xl text-error hover:bg-error/5 transition-colors text-left"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
                        />
                    </svg>
                    Keluar
                </button>
            </div>

            <div class="mt-4 bg-base-100 rounded-xl p-3 border border-base-200">
                <div class="flex items-center gap-3">
                    <div
                        class="w-8 h-8 rounded-lg bg-info/10 flex items-center justify-center text-info"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-4 w-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                        </svg>
                    </div>
                    <div>
                        <p class="text-xs font-semibold text-base-content">
                            v0.2.0
                        </p>
                        <p class="text-[10px] text-base-content/60">
                            Sistem Terbaru
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </aside>

    {#if isMobileMenuOpen}
        <button
            type="button"
            class="lg:hidden fixed inset-0 bg-black/50 z-30 backdrop-blur-sm"
            onclick={() => (isMobileMenuOpen = false)}
            aria-label="Tutup menu"
        ></button>
    {/if}

    <main class="flex-1 min-w-0 min-h-screen overflow-x-hidden">
        <div class="p-4 sm:p-6 lg:p-8 max-w-7xl mx-auto pt-20 lg:pt-8">
            {@render children()}
        </div>
    </main>
</div>

<style>
    aside::-webkit-scrollbar {
        width: 5px;
    }

    aside::-webkit-scrollbar-track {
        background: transparent;
    }

    aside::-webkit-scrollbar-thumb {
        background: hsl(var(--b3) / 0.5);
        border-radius: 3px;
    }

    aside::-webkit-scrollbar-thumb:hover {
        background: hsl(var(--bc) / 0.3);
    }

    aside {
        scrollbar-width: thin;
        scrollbar-color: hsl(var(--b3) / 0.5) transparent;
    }
</style>
